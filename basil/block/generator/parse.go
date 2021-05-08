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

	"github.com/opsidian/basil/basil/generator/parser"
	"github.com/opsidian/basil/basil/schema"
	schemadirectives "github.com/opsidian/basil/basil/schema/directives"
)

type Struct struct {
	Name            string
	InterpreterPath string
	Schema          schema.Schema
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
		Properties: map[string]schema.Schema{},
	}

	for _, directive := range metadata.Directives {
		if err := directive.ApplyToSchema(s); err != nil {
			return nil, err
		}
	}

	var idField string
	var valueField string

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

			if err := addField(s, &idField, &valueField, *field); err != nil {
				return nil, err
			}
		} else {
			fieldStr, err := ParseEmbeddedField(parseCtx, pkg, field)
			if err != nil {
				return nil, fmt.Errorf("failed to parse embedded struct %q: %w", field.Type, err)
			}

			fieldStrSchema := fieldStr.Schema.(*schema.Object)

			for propertyName, property := range fieldStrSchema.Properties {
				if property.GetAnnotation("id") == "true" {
					continue
				}

				if property.GetAnnotation("value") == "true" {
					property.(schema.MetadataAccessor).SetAnnotation("value", "")
				}

				fieldName := propertyName
				if v, ok := fieldStrSchema.PropertyNames[propertyName]; ok {
					fieldName = v
				}

				f := parser.Field{
					Name:         fieldName,
					PropertyName: propertyName,
					Required:     fieldStrSchema.IsPropertyRequired(propertyName),
					Schema:       property,
				}

				if err := addField(s, &idField, &valueField, f); err != nil {
					return nil, err
				}
			}
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

	metadata, err := parser.ParseMetadataFromComments("", comments)
	if err != nil {
		return nil, err
	}
	if len(metadata.Directives) > 0 {
		return nil, fmt.Errorf("directives are not allowed on an embedded struct field")
	}

	switch t := astField.Type.(type) {
	case *ast.Ident:
		astStruct, metadata, err := parser.FindStruct(parseCtx, pkg, t.Name)
		if err != nil {
			return nil, err
		}

		str, err := ParseStruct(parseCtx, astStruct, pkg, t.Name, metadata)
		if err != nil {
			return nil, err
		}

		return str, nil
	case *ast.SelectorExpr:
		pkg := parser.GetImportPath(parseCtx.File, t.X.(*ast.Ident).Name)
		if pkg == "" {
			return nil, fmt.Errorf("failed to find package import for %s", t.X.(*ast.Ident).Name)
		}

		astStruct, metadata, err := parser.FindStruct(parseCtx, pkg, t.Sel.Name)
		if err != nil {
			return nil, err
		}

		str, err := ParseStruct(parseCtx, astStruct, pkg, t.Sel.Name, metadata)
		if err != nil {
			return nil, err
		}

		return str, nil
	default:
		panic(fmt.Errorf("unexpected ast node type: %T", t))
	}
}

func addField(s *schema.Object, idField, valueField *string, field parser.Field) error {
	if _, exists := s.Properties[field.PropertyName]; exists {
		return fmt.Errorf("multiple fields has the same property name: %s", field.PropertyName)
	}

	if schema.HasAnnotationValue(field.Schema, "id", "true") {
		if *idField != "" {
			return fmt.Errorf("multiple id fields were found: %s, %s", *idField, field.Name)
		}
		*idField = field.Name
	}

	if schema.HasAnnotationValue(field.Schema, "value", "true") {
		if *valueField != "" {
			return fmt.Errorf("multiple value fields were found: %s, %s", *valueField, field.Name)
		}
		*valueField = field.Name
	}

	if field.Required {
		if *valueField != "" && *valueField != field.Name {
			return errors.New("when setting a value field then no other fields can be required")
		}
		s.Required = append(s.Required, field.PropertyName)
	}

	s.Properties[field.PropertyName] = field.Schema
	if field.PropertyName != field.Name {
		if s.PropertyNames == nil {
			s.PropertyNames = map[string]string{}
		}
		s.PropertyNames[field.PropertyName] = field.Name
	}

	return nil
}
