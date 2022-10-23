// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generator

import (
	"bytes"
	"go/ast"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/conflowio/conflow/pkg/conflow/annotations"
	"github.com/conflowio/conflow/pkg/conflow/block"
	"github.com/conflowio/conflow/pkg/conflow/generator/parser"
	generatortemplate "github.com/conflowio/conflow/pkg/conflow/generator/template"
	"github.com/conflowio/conflow/pkg/internal/utils"
	"github.com/conflowio/conflow/pkg/schema"
)

// GenerateInterpreter generates an interpreter for the given block
func GenerateInterpreter(
	parseCtx *parser.Context,
	str *ast.StructType,
	pkg string,
	name string,
	comments []*ast.Comment,
) ([]byte, *Struct, error) {
	metadata, err := parser.ParseMetadataFromComments(comments)
	if err != nil {
		return nil, nil, err
	}

	s, err := ParseStruct(parseCtx, str, pkg, name, metadata)
	if err != nil {
		return nil, nil, err
	}

	params := generateTemplateParams(parseCtx, s, pkg)

	bodyTmpl := template.New("body")
	bodyTmpl.Funcs(map[string]interface{}{
		"assignValue": func(s schema.Schema, valueName, resultName string) string {
			return s.AssignValue(params.Imports, valueName, resultName)
		},
		"filterInputs": func(props map[string]schema.Schema) map[string]schema.Schema {
			return filterSchemaProperties(props, func(s schema.Schema) bool {
				return !s.GetReadOnly()
			})
		},
		"filterParams": func(props map[string]schema.Schema) map[string]schema.Schema {
			return filterSchemaProperties(props, func(s schema.Schema) bool {
				return !block.IsBlockSchema(s)
			})
		},
		"filterBlocks": func(props map[string]schema.Schema) map[string]schema.Schema {
			return filterSchemaProperties(props, block.IsBlockSchema)
		},
		"filterNonID": func(props map[string]schema.Schema) map[string]schema.Schema {
			return filterSchemaProperties(props, func(s schema.Schema) bool {
				return s.GetAnnotation(annotations.ID) != "true"
			})
		},
		"filterDefaults": func(props map[string]schema.Schema) map[string]schema.Schema {
			return filterSchemaProperties(props, func(s schema.Schema) bool {
				return s.DefaultValue() != nil
			})
		},
		"getParameterName": func(name string) string {
			return params.Schema.(*schema.Object).ParameterName(name)
		},
		"getFieldName": func(name string) string {
			return params.Schema.(*schema.Object).FieldName(name)
		},
		"getType": func(s schema.Schema) string {
			return s.GoType(params.Imports)
		},
		"isArray": func(s schema.Schema) bool {
			_, ok := s.(*schema.Array)
			return ok
		},
		"isMap": func(s schema.Schema) bool {
			_, ok := s.(*schema.Map)
			return ok
		},
		"title": func(s string) string {
			return cases.Title(language.English, cases.NoLower).String(s)
		},
	})
	if _, parseErr := bodyTmpl.Parse(interpreterTemplate); parseErr != nil {
		return nil, nil, parseErr
	}

	body := &bytes.Buffer{}
	err = bodyTmpl.Execute(body, params)
	if err != nil {
		return nil, nil, err
	}

	header, err := generatortemplate.GenerateHeader(generatortemplate.HeaderParams{
		Package:       params.Package,
		Imports:       params.Imports,
		LocalPrefixes: parseCtx.LocalPrefixes,
	})
	if err != nil {
		return nil, nil, err
	}

	return append(header, body.Bytes()...), s, nil
}

func filterSchemaProperties(props map[string]schema.Schema, f func(p schema.Schema) bool) map[string]schema.Schema {
	res := map[string]schema.Schema{}
	for pn, p := range props {
		if f(p) {
			res[pn] = p
		}
	}
	return res
}

func generateTemplateParams(
	parseCtx *parser.Context,
	s *Struct,
	pkg string,
) *InterpreterTemplateParams {
	imports := map[string]string{
		s.InterpreterPkg: "",
		"fmt":            "fmt",
		"github.com/conflowio/conflow/pkg/conflow": "conflow",
	}

	nameSelector := utils.EnsureUniqueGoPackageSelector(imports, pkg)

	var idPropertyName, valuePropertyName string
	o := s.Schema.(*schema.Object)
	for jsonPropertyName, property := range o.Properties {
		parameterName := o.ParameterName(jsonPropertyName)
		switch {
		case property.GetAnnotation(annotations.ID) == "true":
			idPropertyName = parameterName
		case property.GetAnnotation(annotations.Value) == "true":
			valuePropertyName = parameterName
		}
	}

	pkgName := parseCtx.File.Name.Name
	if s.InterpreterPath != "" {
		parts := strings.Split(strings.Trim(s.InterpreterPath, "/"), "/")
		pkgName = parts[len(parts)-1]
	}

	return &InterpreterTemplateParams{
		Package:               pkgName,
		NameSelector:          nameSelector,
		Name:                  s.Name,
		Schema:                s.Schema,
		IDPropertyName:        idPropertyName,
		ValuePropertyName:     valuePropertyName,
		Imports:               imports,
		Dependencies:          s.Dependencies,
		SchemaPackageSelector: utils.EnsureUniqueGoPackageSelector(imports, "github.com/conflowio/conflow/pkg/schema"),
	}
}
