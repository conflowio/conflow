// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"path"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/schema"
	schemadirectives "github.com/conflowio/conflow/conflow/schema/directives"
	"github.com/conflowio/conflow/util"
)

type Field struct {
	Dependency     string
	Name           string
	PropertyName   string
	Required       bool
	ResultTypeFrom bool
	Schema         schema.Schema
}

func ParseField(
	parseCtx *Context,
	fieldName string,
	astField *ast.Field,
	pkg string,
	fileComments ...*ast.Comment,
) (*Field, error) {
	var comments []*ast.Comment
	if astField.Doc != nil {
		comments = append(comments, astField.Doc.List...)
	}
	comments = append(comments, fileComments...)

	metadata, err := ParseMetadataFromComments(fieldName, comments)
	if err != nil {
		return nil, err
	}

	required := false
	resultType := false

	propertyName := fieldName
	if propertyName != "" && !conflow.IDRegExp.MatchString(propertyName) {
		propertyName = ToSnakeCase(propertyName)
	}

	var jsonPropertyName string
	if astField.Tag != nil {
		tag, err := strconv.Unquote(astField.Tag.Value)
		if err != nil {
			return nil, errors.New("tag is invalid")
		}
		jsonTags := reflect.StructTag(tag).Get("json")
		jsonTagParts := strings.Split(jsonTags, ",")
		jsonName := strings.TrimSpace(jsonTagParts[0])

		if jsonName == "-" {
			return nil, nil
		}

		if jsonName != "" && conflow.IDRegExp.MatchString(jsonName) {
			propertyName = jsonName
			jsonPropertyName = jsonName
		}
	}

	var dependencyName string

	for _, directive := range metadata.Directives {
		if _, ok := directive.(*schemadirectives.Ignore); ok {
			if _, ok := parseCtx.Parent.(*ast.StructType); !ok {
				return nil, errors.New("the @ignore annotation can only be used on struct fields")
			}
			return nil, nil
		}
	}

	fieldSchema, _, err := getSchemaForField(parseCtx, astField.Type, pkg)
	if err != nil {
		return nil, err
	}

	for _, directive := range metadata.Directives {
		switch d := directive.(type) {
		case *schemadirectives.Dependency:
			if _, ok := parseCtx.Parent.(*ast.StructType); !ok {
				return nil, errors.New("the @dependency annotation can only be used on struct fields")
			}

			if d.Name != "" {
				if !util.StringSliceContains(validDependencies, d.Name) {
					return nil, fmt.Errorf("%s dependency is invalid, valid values are: %s", d.Name, strings.Join(validDependencies, ", "))
				}
				dependencyName = d.Name
			} else {
				if util.StringSliceContains(validDependencies, fieldName) {
					dependencyName = fieldName
				} else {
					return nil, errors.New("dependency can not be inferred from the field name, please set the name explicitly (@dependency \"name\"")
				}
			}

			var actualType string
			switch s := fieldSchema.(type) {
			case *schema.Reference:
				actualType = strings.TrimPrefix(s.Ref, "http://conflow.schema/")
			case *schema.Untyped:
				actualType = "interface{}"
			default:
				actualType = s.TypeString()
			}

			if dependencyTypes[dependencyName] != actualType {
				return nil, fmt.Errorf("%s dependency type can only be defined on a %s field", dependencyName, dependencyTypes[dependencyName])
			}

		case *schemadirectives.Required:
			if _, ok := parseCtx.Parent.(*ast.StructType); !ok {
				return nil, errors.New("the @required annotation can only be used on struct fields")
			}
			required = true
		case *schemadirectives.ResultType:
			if _, ok := parseCtx.Parent.(*ast.FuncType); !ok {
				return nil, errors.New("the @result_type annotation can only be used on function parameters")
			}
			resultType = true
		case *schemadirectives.Name:
			if jsonPropertyName != "" && d.Value != jsonPropertyName {
				return nil, errors.New("name directive's value must match the name in the json struct tag")
			}
			propertyName = d.Value
		}
	}

	fieldSchema.(schema.MetadataAccessor).SetDescription(metadata.Description)

	meta, ok := fieldSchema.(schema.MetadataAccessor)
	if !ok {
		panic(fmt.Errorf("metadata is not writable on schema %T", fieldSchema))
	}

	for _, directive := range metadata.Directives {
		if err := directive.ApplyToSchema(fieldSchema); err != nil {
			return nil, err
		}
	}

	if schema.HasAnnotationValue(fieldSchema, conflow.AnnotationID, "true") &&
		schema.HasAnnotationValue(fieldSchema, conflow.AnnotationValue, "true") {
		return nil, errors.New("the id field can not be marked as the value field")
	}

	if fieldSchema.GetReadOnly() && !schema.HasAnnotationValue(fieldSchema, conflow.AnnotationID, "true") {
		meta.SetAnnotation(conflow.AnnotationEvalStage, "close")
	}

	if schema.HasAnnotationValue(fieldSchema, conflow.AnnotationGenerated, "true") {
		meta.SetAnnotation(conflow.AnnotationEvalStage, "init")
		required = true
	}

	return &Field{
		Dependency:     dependencyName,
		Name:           fieldName,
		PropertyName:   propertyName,
		Required:       required,
		ResultTypeFrom: resultType,
		Schema:         fieldSchema,
	}, nil
}

func getSchemaForField(parseCtx *Context, typeNode ast.Expr, pkg string) (schema.Schema, bool, error) {
	switch tn := typeNode.(type) {
	case *ast.Ident:
		var s schema.Schema
		switch tn.String() {
		case "int64":
			s = &schema.Integer{}
		case "float64":
			s = &schema.Number{}
		case "bool":
			s = &schema.Boolean{}
		case "string":
			s = &schema.String{}
		default:
			r, _ := utf8.DecodeRuneInString(tn.String())
			if !unicode.IsUpper(r) {
				return nil, false, fmt.Errorf("type %s is not allowed", tn.String())
			}

			filePath := parseCtx.FileSet.File(parseCtx.File.Pos()).Name()

			_, _, err := FindStruct(parseCtx.WithWorkdir(path.Dir(filePath)), pkg, tn.String())
			if err != nil {
				if _, notFound := err.(errStructNotFound); notFound {
					return nil, false, fmt.Errorf("type %s is not allowed", tn.String())
				}
				return nil, false, fmt.Errorf("failed to parse %s: %w", tn.String(), err)
			}

			s = &schema.Reference{
				Ref: fmt.Sprintf("http://conflow.schema/%s.%s", pkg, tn.String()),
			}

			return s, true, nil
		}

		return s, false, nil
	case *ast.ArrayType:
		itemsSchema, isRef, err := getSchemaForField(parseCtx, tn.Elt, pkg)
		if err != nil {
			return nil, false, err
		}

		return &schema.Array{
			Items: itemsSchema,
		}, isRef, nil
	case *ast.MapType:
		keyIdent, ok := tn.Key.(*ast.Ident)
		if !ok || keyIdent.String() != "string" {
			return nil, false, fmt.Errorf("only string map keys are supported")
		}

		propertiesSchema, isRef, err := getSchemaForField(parseCtx, tn.Value, pkg)
		if err != nil {
			return nil, false, err
		}
		if isRef {
			return nil, false, fmt.Errorf("only single blocks or slice of blocks are supported")
		}

		return &schema.Map{
			AdditionalProperties: propertiesSchema,
		}, false, nil
	case *ast.StarExpr:
		res, isRef, err := getSchemaForField(parseCtx, tn.X, pkg)
		if err != nil {
			return nil, false, err
		}

		res.(schema.MetadataAccessor).SetPointer(true)

		return res, isRef, nil
	case *ast.SelectorExpr:
		if xIdent, ok := tn.X.(*ast.Ident); ok {
			path := GetImportPath(parseCtx.File, xIdent.Name)
			if path == "" {
				return nil, false, fmt.Errorf("could not find import path for %s", xIdent.Name)
			}

			var s schema.Schema
			switch path + "." + tn.Sel.Name {
			case "github.com/conflowio/conflow/conflow.ID":
				s = &schema.String{
					Format: schema.FormatConflowID,
					Metadata: schema.Metadata{
						ReadOnly: true,
					},
				}
			case "io.ReadCloser":
				s = &schema.ByteStream{}
			case "time.Time":
				s = &schema.Time{}
			case "time.Duration":
				s = &schema.TimeDuration{}
			default:
				_, _, err := FindStruct(parseCtx, path, tn.Sel.Name)
				if err != nil {
					if _, notFound := err.(errStructNotFound); notFound {
						return nil, false, fmt.Errorf("type is not allowed: %s.%s", xIdent.Name, tn.Sel.Name)
					}
					return nil, false, fmt.Errorf("failed to parse %s.%s: %w", xIdent.Name, tn.Sel.Name, err)
				}

				s = &schema.Reference{
					Ref: fmt.Sprintf("http://conflow.schema/%s.%s", path, tn.Sel.Name),
				}
				return s, true, nil
			}

			return s, false, nil
		}
		return nil, false, fmt.Errorf("unexpected ast node: %T: %v", typeNode, typeNode)
	case *ast.InterfaceType:
		return &schema.Untyped{}, false, nil
	default:
		return nil, false, fmt.Errorf("unexpected ast node: %T: %v", typeNode, typeNode)
	}
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(name string) string {
	name = matchFirstCap.ReplaceAllString(name, "${1}_${2}")
	name = matchAllCap.ReplaceAllString(name, "${1}_${2}")
	return strings.ToLower(name)
}

func GetImportPath(file *ast.File, name string) string {
	if name == "time" || name == "io" {
		return name
	}

	for _, i := range file.Imports {
		path, _ := strconv.Unquote(i.Path.Value)
		if i.Name != nil {
			if i.Name.Name == name {
				return path
			}
		} else {
			parts := strings.Split(path, "/")
			if parts[len(parts)-1] == name {
				return path
			}
		}
	}
	return ""
}
