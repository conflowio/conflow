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
	"path"
	"strings"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/conflow/annotations"
	"github.com/conflowio/conflow/pkg/conflow/generator/parser"
	"github.com/conflowio/conflow/pkg/schema"
	schemadirectives "github.com/conflowio/conflow/pkg/schema/directives"
)

type Struct struct {
	Name            string
	InterpreterPath string
	InterpreterPkg  string
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

	var interpreterPath string
	interpreterPkg := pkg
	if bd != nil && bd.Path != "" {
		interpreterPath = bd.Path
		interpreterPkg = path.Clean(path.Join(interpreterPkg, bd.Path))
	}

	s := &schema.Object{
		Metadata: schema.Metadata{
			ID:          fmt.Sprintf("%s.%s", pkg, name),
			Description: metadata.Description,
		},
		Properties: map[string]schema.Schema{},
	}

	if strings.HasPrefix(s.Metadata.Description, name+" ") {
		s.Metadata.Description = strings.Replace(s.Metadata.Description, name+" ", "It ", 1)
	}

	for _, directive := range metadata.Directives {
		if err := directive.ApplyToSchema(s); err != nil {
			return nil, err
		}
	}

	var idField, valueField, keyField string
	var dependencies []parser.Dependency

	parseCtx = parseCtx.WithParent(str)

	for _, strField := range str.Fields.List {
		if len(strField.Names) > 0 {
			fieldName := strField.Names[0].String()

			field, err := parser.ParseField(parseCtx, fieldName, strField, pkg)
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

			if err := addField(s, &idField, &valueField, &keyField, *field); err != nil {
				return nil, err
			}
		} else {
			fieldStr, err := ParseEmbeddedField(parseCtx, pkg, strField)
			if err != nil {
				return nil, fmt.Errorf("failed to parse embedded struct %q: %w", strField.Type, err)
			}

			fieldStrSchema := fieldStr.Schema.(*schema.Object)

			for jsonPropertyName, property := range fieldStrSchema.Properties {
				parameterName := fieldStrSchema.ParameterName(jsonPropertyName)
				if property.GetAnnotation(annotations.ID) == "true" {
					continue
				}

				f := parser.Field{
					Name:             fieldStrSchema.FieldName(parameterName),
					ParameterName:    parameterName,
					Required:         fieldStrSchema.IsPropertyRequired(jsonPropertyName),
					Schema:           property,
					JSONPropertyName: jsonPropertyName,
				}

				if err := addField(s, &idField, &valueField, &keyField, f); err != nil {
					return nil, err
				}
			}
		}
	}

	if valueField != "" && len(s.Required) > 0 && (len(s.Required) > 1 || s.Required[0] != valueField) {
		return nil, errors.New("when setting a value field then no other fields can be required")
	}

	if s.GetAnnotation(annotations.Type) == conflow.BlockTypeGenerator {
		hasGeneratedField := false
		for _, p := range s.Properties {
			if p.GetAnnotation(annotations.Generated) == "true" {
				hasGeneratedField = true
				break
			}
		}

		if !hasGeneratedField {
			return nil, errors.New("a generator block must have at least one field marked as @generated")
		}
	}

	return &Struct{
		Name:            name,
		InterpreterPath: interpreterPath,
		InterpreterPkg:  interpreterPkg,
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

func addField(s *schema.Object, idField, valueField, keyField *string, field parser.Field) error {
	if _, exists := s.Properties[field.JSONPropertyName]; exists {
		return fmt.Errorf("multiple fields has the same JSON property name: %s", field.JSONPropertyName)
	}

	if field.Schema.GetAnnotation(annotations.ID) == "true" {
		if *idField != "" {
			return fmt.Errorf("multiple id fields were found: %s, %s", *idField, field.ParameterName)
		}
		*idField = field.ParameterName
	}

	if field.Schema.GetAnnotation(annotations.Value) == "true" {
		if *valueField != "" {
			return fmt.Errorf("multiple value fields were found: %s, %s", *valueField, field.ParameterName)
		}
		*valueField = field.ParameterName
	}

	if field.Schema.GetAnnotation(annotations.Key) == "true" {
		if *keyField != "" {
			return fmt.Errorf("multiple key fields were found: %s, %s", *keyField, field.ParameterName)
		}
		*keyField = field.ParameterName
	}

	if field.Required {
		s.Required = append(s.Required, field.ParameterName)
	}

	s.Properties[field.JSONPropertyName] = field.Schema

	if field.ParameterName != field.JSONPropertyName {
		if s.ParameterNames == nil {
			s.ParameterNames = map[string]string{}
		}
		s.ParameterNames[field.JSONPropertyName] = field.ParameterName
	}

	if field.Name != field.JSONPropertyName {
		if s.FieldNames == nil {
			s.FieldNames = map[string]string{}
		}
		s.FieldNames[field.JSONPropertyName] = field.Name
	}

	return nil
}
