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

	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/block"
	"github.com/conflowio/conflow/conflow/generator/parser"
	"github.com/conflowio/conflow/conflow/schema"
	"github.com/conflowio/conflow/internal/utils"
)

// GenerateInterpreter generates an interpreter for the given block
func GenerateInterpreter(
	parseCtx *parser.Context,
	str *ast.StructType,
	pkg string,
	name string,
	comments []*ast.Comment,
) ([]byte, *Struct, error) {
	metadata, err := parser.ParseMetadataFromComments(name, comments)
	if err != nil {
		return nil, nil, err
	}

	s, err := ParseStruct(parseCtx, str, pkg, name, metadata)
	if err != nil {
		return nil, nil, err
	}

	params := generateTemplateParams(parseCtx, s, pkg)

	bodyTmpl := template.New("block_interpreter_body")
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
				return !schema.HasAnnotationValue(s, conflow.AnnotationID, "true")
			})
		},
		"filterDefaults": func(props map[string]schema.Schema) map[string]schema.Schema {
			return filterSchemaProperties(props, func(s schema.Schema) bool {
				return s.DefaultValue() != nil
			})
		},
		"getPropertyName": func(name string) string {
			if p, ok := params.Schema.(schema.ObjectKind).GetPropertyNames()[name]; ok {
				return p
			}
			return name
		},
		"getType": func(s schema.Schema) string {
			return s.GoType(params.Imports)
		},
		"isArray": func(s schema.Schema) bool {
			_, ok := s.(schema.ArrayKind)
			return ok
		},
		"title": func(s string) string {
			return strings.Title(s)
		},
	})
	if _, parseErr := bodyTmpl.Parse(interpreterTemplate); parseErr != nil {
		return nil, nil, parseErr
	}

	res := &bytes.Buffer{}
	err = bodyTmpl.Execute(res, params)
	if err != nil {
		return nil, nil, err
	}

	body := res.Bytes()

	headerTmpl := template.New("block_interpreter_header")
	headerTmpl.Funcs(map[string]interface{}{
		"last": func(path string) string {
			parts := strings.Split(path, "/")
			return parts[len(parts)-1]
		},
		"sortedImportKeys": parser.SortedImportKeys,
	})
	if _, parseErr := headerTmpl.Parse(interpreterHeaderTemplate); parseErr != nil {
		return nil, nil, parseErr
	}

	res = &bytes.Buffer{}
	err = headerTmpl.Execute(res, params)
	if err != nil {
		return nil, nil, err
	}

	res.Write(body)

	return res.Bytes(), s, nil
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
		".":       pkg,
		"fmt":     "fmt",
		"conflow": "github.com/conflowio/conflow/conflow",
		"schema":  "github.com/conflowio/conflow/conflow/schema",
	}

	var nameSelector string
	if s.InterpreterPath != "" {
		nameSelector = utils.EnsureUniqueGoPackageName(imports, pkg) + "."
	}

	var idPropertyName, valuePropertyName string
	for name, property := range s.Schema.(schema.ObjectKind).GetProperties() {
		switch {
		case schema.HasAnnotationValue(property, conflow.AnnotationID, "true"):
			idPropertyName = name
		case schema.HasAnnotationValue(property, conflow.AnnotationValue, "true"):
			valuePropertyName = name
		}
	}

	pkgName := parseCtx.File.Name.Name
	if s.InterpreterPath != "" {
		parts := strings.Split(strings.Trim(s.InterpreterPath, "/"), "/")
		pkgName = parts[len(parts)-1]
	}

	return &InterpreterTemplateParams{
		Package:           pkgName,
		NameSelector:      nameSelector,
		Name:              s.Name,
		Schema:            s.Schema,
		IDPropertyName:    idPropertyName,
		ValuePropertyName: valuePropertyName,
		Imports:           imports,
		Dependencies:      s.Dependencies,
	}
}
