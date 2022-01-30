// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generator

import (
	"errors"
	"fmt"
	"go/ast"
	"strings"

	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/conflow/generator/parser"
	"github.com/conflowio/conflow/src/conflow/schema"
	schemadirectives "github.com/conflowio/conflow/src/conflow/schema/directives"
)

type Struct struct {
	Name            string
	InterpreterPath string
	Schema          schema.Schema
	Dependencies    []parser.Dependency
}

// ParseStruct parses all fields of a given go struct
func ParseStruct(
	parseCtx *parser.Context,
	str *ast.StructType,
	pkg string,
	name string,
	metadata *parser.Metadata,
) (*Struct, error) {
	var bd *schemadirectives.Block
	for _, d := range metadata.Directives {
		if v, ok := d.(*schemadirectives.Block); ok {
			bd = v
			break
		}
	}

	s := &schema.Object{
		Metadata: schema.Metadata{
			Description: metadata.Description,
		},
		Name:       name,
		Parameters: map[string]schema.Schema{},
	}

	if strings.HasPrefix(s.Metadata.Description, name+" ") {
		s.Metadata.Description = strings.Replace(s.Metadata.Description, name+" ", "It ", 1)
	}

	for _, directive := range metadata.Directives {
		if err := directive.ApplyToSchema(s); err != nil {
			return nil, err
		}
	}

	var idField string
	var valueField string
	var dependencies []parser.Dependency

	parseCtx = parseCtx.WithParent(str)

	for _, field := range str.Fields.List {
		if len(field.Names) > 0 {
			fieldName := field.Names[0].String()

			field, err := parser.ParseField(parseCtx, fieldName, field, pkg)
			if err != nil {
				return nil, fmt.Errorf("failed to parse field %q: %w", fieldName, err)
			}

			if field == nil {
				continue
			}

			if field.Dependency != "" {
				dependencies = append(dependencies, parser.Dependency{
					Name:      fieldName,
					FieldName: field.Dependency,
				})
				continue
			}

			if err := addField(s, &idField, &valueField, *field); err != nil {
				return nil, err
			}
		} else {
			fieldStr, err := ParseEmbeddedField(parseCtx, pkg, field)
			if err != nil {
				return nil, fmt.Errorf("failed to parse embedded struct %q: %w", field.Type, err)
			}

			fieldStrSchema := fieldStr.Schema.(*schema.Object)

			for parameterName, parameter := range fieldStrSchema.Parameters {
				if parameter.GetAnnotation(conflow.AnnotationID) == "true" {
					continue
				}

				if parameter.GetAnnotation(conflow.AnnotationValue) == "true" {
					parameter.(schema.MetadataAccessor).SetAnnotation(conflow.AnnotationValue, "")
				}

				f := parser.Field{
					Name:             fieldStrSchema.GetFieldName(parameterName),
					ParameterName:    parameterName,
					Required:         fieldStrSchema.IsParameterRequired(parameterName),
					Schema:           parameter,
					JSONPropertyName: fieldStrSchema.GetJSONPropertyName(parameterName),
				}

				if err := addField(s, &idField, &valueField, f); err != nil {
					return nil, err
				}
			}
		}
	}

	if s.GetAnnotation(conflow.AnnotationType) == conflow.BlockTypeGenerator {
		hasGeneratedField := false
		for _, p := range s.Parameters {
			if p.GetAnnotation(conflow.AnnotationGenerated) == "true" {
				hasGeneratedField = true
				break
			}
		}

		if !hasGeneratedField {
			return nil, errors.New("a generator block must have at least one field marked as @generated")
		}
	}

	var interpreterPath string
	if bd != nil {
		interpreterPath = bd.Path
	}

	return &Struct{
		Name:            name,
		InterpreterPath: interpreterPath,
		Schema:          s,
		Dependencies:    dependencies,
	}, nil
}

func ParseEmbeddedField(
	parseCtx *parser.Context,
	pkg string,
	astField *ast.Field,
) (*Struct, error) {
	var comments []*ast.Comment
	if astField.Comment != nil {
		comments = astField.Comment.List
	}

	metadata, err := parser.ParseMetadataFromComments(comments)
	if err != nil {
		return nil, err
	}
	if len(metadata.Directives) > 0 {
		return nil, fmt.Errorf("directives are not allowed on an embedded struct field")
	}

	switch t := astField.Type.(type) {
	case *ast.Ident:
		astFile, astStruct, metadata, err := parser.FindType(parseCtx, pkg, t.Name)
		if err != nil {
			return nil, err
		}

		str, err := ParseStruct(parseCtx.WithFile(astFile), astStruct.(*ast.StructType), pkg, t.Name, metadata)
		if err != nil {
			return nil, err
		}

		return str, nil
	case *ast.SelectorExpr:
		pkg := parser.GetImportPath(parseCtx.File, t.X.(*ast.Ident).Name)
		if pkg == "" {
			return nil, fmt.Errorf("failed to find package import for %s", t.X.(*ast.Ident).Name)
		}

		astFile, astStruct, metadata, err := parser.FindType(parseCtx, pkg, t.Sel.Name)
		if err != nil {
			return nil, err
		}

		str, err := ParseStruct(parseCtx.WithFile(astFile), astStruct.(*ast.StructType), pkg, t.Sel.Name, metadata)
		if err != nil {
			return nil, err
		}

		return str, nil
	default:
		panic(fmt.Errorf("unexpected ast node type: %T", t))
	}
}

func addField(s *schema.Object, idField, valueField *string, field parser.Field) error {
	if _, exists := s.Parameters[field.ParameterName]; exists {
		return fmt.Errorf("multiple fields has the same property name: %s", field.ParameterName)
	}

	if field.Schema.GetAnnotation(conflow.AnnotationID) == "true" {
		if *idField != "" {
			return fmt.Errorf("multiple id fields were found: %s, %s", *idField, field.Name)
		}
		*idField = field.Name
	}

	if field.Schema.GetAnnotation(conflow.AnnotationValue) == "true" {
		if *valueField != "" {
			return fmt.Errorf("multiple value fields were found: %s, %s", *valueField, field.Name)
		}
		*valueField = field.Name
	}

	if field.Required {
		if *valueField != "" && *valueField != field.Name {
			return errors.New("when setting a value field then no other fields can be required")
		}
		s.Required = append(s.Required, field.ParameterName)
	}

	s.Parameters[field.ParameterName] = field.Schema

	if field.ParameterName != field.JSONPropertyName {
		if s.JSONPropertyNames == nil {
			s.JSONPropertyNames = map[string]string{}
		}
		s.JSONPropertyNames[field.ParameterName] = field.JSONPropertyName
	}

	if field.Name != field.JSONPropertyName {
		if s.FieldNames == nil {
			s.FieldNames = map[string]string{}
		}
		s.FieldNames[field.JSONPropertyName] = field.Name
	}

	return nil
}
