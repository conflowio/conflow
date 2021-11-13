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

	"github.com/conflowio/conflow/conflow/generator/parser"
	"github.com/conflowio/conflow/conflow/schema"
	schemadirectives "github.com/conflowio/conflow/conflow/schema/directives"
)

type Function struct {
	Name            string
	InterpreterPath string
	ReturnsError    bool
	Schema          schema.Schema
}

// ParseFunction parses all results of a given go function
func ParseFunction(
	parseCtx *parser.Context,
	fun *ast.FuncType,
	pkg string,
	name string,
	metadata *parser.Metadata,
) (*Function, error) {
	var fd *schemadirectives.Function
	for _, d := range metadata.Directives {
		if v, ok := d.(*schemadirectives.Function); ok {
			fd = v
			break
		}
	}

	s := &schema.Function{
		Metadata: schema.Metadata{
			Description: metadata.Description,
		},
	}

	if strings.HasPrefix(s.Metadata.Description, name+" ") {
		s.Metadata.Description = strings.Replace(s.Metadata.Description, name+" ", "It ", 1)
	}

	for _, directive := range metadata.Directives {
		if err := directive.ApplyToSchema(s); err != nil {
			return nil, err
		}
	}

	parseCtx = parseCtx.WithParent(fun)

	var resultTypeField string

	pos := fun.Pos()
	for _, param := range fun.Params.List {
		var comments []*ast.Comment
		if parseCtx.File.End() < fun.Pos() || parseCtx.File.Pos() > param.Pos() {
			continue
		}

		for _, cg := range parseCtx.File.Comments {
			if cg.Pos() >= pos && cg.Pos() < param.Pos() {
				comments = append(comments, cg.List...)
			}
		}

		for _, name := range param.Names {
			field, err := parser.ParseField(parseCtx, name.String(), param, pkg, comments...)
			if err != nil {
				return nil, fmt.Errorf("parameter %s is invalid: %w", name.String(), err)
			}

			if field.Variadic {
				s.AdditionalParameters = &schema.NamedSchema{
					Name:   name.String(),
					Schema: field.Schema,
				}
			} else {
				s.Parameters = append(s.Parameters, schema.NamedSchema{
					Name:   name.String(),
					Schema: field.Schema,
				})
			}

			if field.ResultTypeFrom {
				if resultTypeField != "" {
					return nil, errors.New("only one parameter can be marked as return type")
				}
				resultTypeField = name.String()
			}
		}
		pos = param.End()
	}

	s.ResultTypeFrom = resultTypeField

	if fun.Results == nil || len(fun.Results.List) == 0 || len(fun.Results.List) > 2 {
		return nil, fmt.Errorf("the function must return with a single value, or a single value and an error")
	}

	field, err := parser.ParseField(parseCtx, "", fun.Results.List[0], pkg)
	if err != nil {
		return nil, fmt.Errorf("result value is invalid: %w", err)
	}

	s.Result = field.Schema

	if len(fun.Results.List) == 2 {
		if id, ok := fun.Results.List[1].Type.(*ast.Ident); !ok || id.String() != "error" {
			return nil, fmt.Errorf("the function must return an error as the second return value")
		}
	}

	var interpreterPath string
	if fd != nil {
		interpreterPath = fd.Path
	}

	return &Function{
		Name:            name,
		InterpreterPath: interpreterPath,
		ReturnsError:    len(fun.Results.List) == 2,
		Schema:          s,
	}, nil
}
