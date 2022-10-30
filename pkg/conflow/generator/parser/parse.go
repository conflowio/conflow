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
	"strconv"
	"strings"

	"github.com/conflowio/conflow/pkg/conflow/annotations"
	"github.com/conflowio/conflow/pkg/internal/utils"
	"github.com/conflowio/conflow/pkg/schema"
	schemadirectives "github.com/conflowio/conflow/pkg/schema/directives"
	"github.com/conflowio/conflow/pkg/util"
)

type Field struct {
	Dependency       string
	JSONPropertyName string
	Name             string
	ParameterName    string
	Required         bool
	ResultTypeFrom   bool
	Schema           schema.Schema
	Variadic         bool
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

	metadata, err := ParseMetadataFromComments(comments)
	if err != nil {
		return nil, err
	}

	jsonPropertyName := fieldName
	if astField.Tag != nil {
		tag, err := strconv.Unquote(astField.Tag.Value)
		if err != nil {
			return nil, errors.New("tag is invalid")
		}
		jsonTags := reflect.StructTag(tag).Get("json")
		jsonTagParts := strings.Split(jsonTags, ",")
		jsonName := strings.TrimSpace(jsonTagParts[0])

		// TODO: ignore "-"
		if jsonName != "" && jsonName != "-" {
			jsonPropertyName = jsonName
		}
	}

	var parameterName string
	if jsonPropertyName != "" && schema.NameRegExp.MatchString(jsonPropertyName) {
		parameterName = jsonPropertyName
	}
	if parameterName == "" && fieldName != "" {
		if schema.NameRegExp.MatchString(fieldName) {
			parameterName = fieldName
		} else {
			parameterName = utils.ToSnakeCase(fieldName)
			if !schema.NameRegExp.MatchString(parameterName) {
				return nil, fmt.Errorf("unable to generate a valid parameter name from field %q, please use @name to define one", fieldName)
			}
		}
	}

	for _, directive := range metadata.Directives {
		if _, ok := directive.(*schemadirectives.Ignore); ok {
			if _, ok := parseCtx.Parent.(*ast.StructType); !ok {
				return nil, errors.New("the @ignore annotation can only be used on struct fields")
			}
			return nil, nil
		}
	}

	fieldSchema, _, err := GetSchemaForType(parseCtx, astField.Type, pkg, metadata)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(fieldSchema.GetDescription(), fieldName+" ") {
		fieldSchema.(schema.MetadataAccessor).SetDescription(
			strings.Replace(fieldSchema.GetDescription(), fieldName+" ", "It ", 1),
		)
	}

	var dependencyName string
	var required bool
	var resultType bool

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
				name := strings.ToLower(fieldName[0:1]) + fieldName[1:]
				if util.StringSliceContains(validDependencies, name) {
					dependencyName = name
				} else {
					return nil, errors.New("dependency can not be inferred from the field name, please set the name explicitly (@dependency \"name\"")
				}
			}

			var actualType string
			switch s := fieldSchema.(type) {
			case *schema.Reference:
				actualType = s.Ref
			case *schema.Any:
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
			parameterName = d.Value
		}
	}

	if fieldSchema.GetAnnotation(annotations.ID) == "true" &&
		fieldSchema.GetAnnotation(annotations.Value) == "true" {
		return nil, errors.New("the id field can not be marked as the value field")
	}

	if fieldSchema.GetReadOnly() && fieldSchema.GetAnnotation(annotations.ID) != "true" {
		fieldSchema.(schema.MetadataAccessor).SetAnnotation(annotations.EvalStage, "close")
	}

	if fieldSchema.GetAnnotation(annotations.Generated) == "true" {
		fieldSchema.(schema.MetadataAccessor).SetAnnotation(annotations.EvalStage, "init")
		required = true
	}

	var variadic bool
	if _, ok := astField.Type.(*ast.Ellipsis); ok {
		variadic = true
	}

	return &Field{
		Dependency:       dependencyName,
		JSONPropertyName: jsonPropertyName,
		Name:             fieldName,
		ParameterName:    parameterName,
		Required:         required,
		ResultTypeFrom:   resultType,
		Schema:           fieldSchema,
		Variadic:         variadic,
	}, nil
}

func GetSchemaForType(parseCtx *Context, typeNode ast.Expr, pkg string, metadata *Metadata) (schema.Schema, bool, error) {
	s, isRef, err := getBaseSchemaForType(parseCtx, typeNode, pkg)
	if err != nil {
		return nil, false, err
	}

	if metadata.Description != "" {
		s.(schema.MetadataAccessor).SetDescription(metadata.Description)
	}

	for _, directive := range metadata.Directives {
		if err := directive.ApplyToSchema(s); err != nil {
			return nil, false, err
		}
	}

	return s, isRef, nil
}

func getBaseSchemaForType(parseCtx *Context, typeNode ast.Expr, pkg string) (schema.Schema, bool, error) {
	switch tn := typeNode.(type) {
	case *ast.Ident:
		switch tn.String() {
		case "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "float32",
			"complex64", "complex128", "int", "uint", "uintptr", "byte", "rune":
			return nil, false, fmt.Errorf("type %s is not allowed", tn.String())
		case "int64":
			return &schema.Integer{}, false, nil
		case "float64":
			return &schema.Number{}, false, nil
		case "bool":
			return &schema.Boolean{}, false, nil
		case "string":
			return &schema.String{}, false, nil
		default:
			filePath := parseCtx.FileSet.File(parseCtx.File.Pos()).Name()

			astFile, expr, metadata, err := FindType(parseCtx.WithWorkdir(path.Dir(filePath)), pkg, tn.String())
			if err != nil {
				return nil, false, fmt.Errorf("failed to parse %s: %w", tn.String(), err)
			}

			switch e := expr.(type) {
			case *ast.StructType:
				return &schema.Reference{
					Ref: fmt.Sprintf("%s.%s", pkg, tn.String()),
				}, true, nil
			default:
				return GetSchemaForType(parseCtx.WithFile(astFile), e, pkg, metadata)
			}
		}
	case *ast.ArrayType:
		// []byte (or []uint8) is a special case
		if ident, ok := tn.Elt.(*ast.Ident); ok {
			if ident.String() == "byte" || ident.String() == "uint8" {
				if formatName, _, ok := schema.GetFormatForType("[]uint8"); ok {
					return &schema.String{Format: formatName}, false, nil
				}
			}
		}
		itemsSchema, isRef, err := getBaseSchemaForType(parseCtx, tn.Elt, pkg)
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

		propertiesSchema, isRef, err := getBaseSchemaForType(parseCtx, tn.Value, pkg)
		if err != nil {
			return nil, false, err
		}

		return &schema.Map{
			AdditionalProperties: propertiesSchema,
		}, isRef, nil
	case *ast.StarExpr:
		res, isRef, err := getBaseSchemaForType(parseCtx, tn.X, pkg)
		if err != nil {
			return nil, false, err
		}

		nullable, ok := res.(schema.Nullable)
		if !ok {
			return nil, false, fmt.Errorf("%s schema type doesn't allow pointer types", res.Type())
		}
		nullable.SetNullable(true)

		return res, isRef, nil
	case *ast.SelectorExpr:
		if xIdent, ok := tn.X.(*ast.Ident); ok {
			path := GetImportPath(parseCtx.File, xIdent.Name)
			if path == "" {
				return nil, false, fmt.Errorf("could not find import path for %s", xIdent.Name)
			}

			switch path + "." + tn.Sel.Name {
			case "github.com/conflowio/conflow/pkg/conflow.ID":
				return &schema.String{
					Format: schema.FormatConflowID,
					Metadata: schema.Metadata{
						ReadOnly: true,
					},
				}, false, nil
			case "io.ReadCloser":
				return &schema.ByteStream{}, false, nil
			}

			if formatName, _, ok := schema.GetFormatForType(path + "." + tn.Sel.Name); ok {
				return &schema.String{Format: formatName}, false, nil
			}

			astFile, expr, metadata, err := FindType(parseCtx, path, tn.Sel.Name)
			if err != nil {
				return nil, false, fmt.Errorf("failed to parse %s.%s: %w", xIdent.Name, tn.Sel.Name, err)
			}

			switch e := expr.(type) {
			case *ast.StructType:
				return &schema.Reference{
					Ref: fmt.Sprintf("%s.%s", path, tn.Sel.Name),
				}, true, nil
			default:
				return GetSchemaForType(parseCtx.WithFile(astFile), e, path, metadata)
			}
		}
		return nil, false, fmt.Errorf("unexpected ast node: %T: %v", typeNode, typeNode)
	case *ast.InterfaceType:
		return &schema.Any{}, false, nil
	case *ast.Ellipsis:
		return getBaseSchemaForType(parseCtx, tn.Elt, pkg)
	default:
		return nil, false, fmt.Errorf("unexpected ast node: %T: %v", typeNode, typeNode)
	}
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
