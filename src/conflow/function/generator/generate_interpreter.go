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

	"github.com/conflowio/conflow/src/conflow/generator/parser"
	generatortemplate "github.com/conflowio/conflow/src/conflow/generator/template"
	"github.com/conflowio/conflow/src/internal/utils"
	"github.com/conflowio/conflow/src/schema"
)

// GenerateInterpreter generates an interpreter for the given function
func GenerateInterpreter(
	parseCtx *parser.Context,
	fun *ast.FuncType,
	pkg string,
	name string,
	comments []*ast.Comment,
) ([]byte, *Function, error) {
	metadata, err := parser.ParseMetadataFromComments(comments)
	if err != nil {
		return nil, nil, err
	}

	if strings.HasPrefix(metadata.Description, name+" ") {
		metadata.Description = strings.Replace(metadata.Description, name+" ", "It ", 1)
	}

	f, err := ParseFunction(parseCtx, fun, pkg, name, metadata)
	if err != nil {
		return nil, nil, err
	}

	params := generateTemplateParams(parseCtx, f, pkg)

	bodyTmpl := template.New("body")
	bodyTmpl.Funcs(map[string]interface{}{
		"assignValue": func(s schema.Schema, valueName, resultName string) string {
			return s.AssignValue(params.Imports, valueName, resultName)
		},
		"getType": func(s schema.Schema) string {
			return s.GoType(params.Imports)
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

	res, err = generatortemplate.GenerateHeader(generatortemplate.HeaderParams{
		Package: params.Package,
		Imports: params.Imports,
	})
	if err != nil {
		return nil, nil, err
	}

	res.Write(body)

	return res.Bytes(), f, nil
}

func generateTemplateParams(
	parseCtx *parser.Context,
	f *Function,
	pkg string,
) *InterpreterTemplateParams {
	imports := map[string]string{
		pkg: "",
	}

	var nameSelector string
	if f.InterpreterPath != "" {
		nameSelector = utils.EnsureUniqueGoPackageSelector(imports, pkg)
	}

	pkgName := parseCtx.File.Name.Name
	if f.InterpreterPath != "" {
		parts := strings.Split(strings.Trim(f.InterpreterPath, "/"), "/")
		pkgName = parts[len(parts)-1]
	}

	return &InterpreterTemplateParams{
		Package:               pkgName,
		Name:                  strings.ToUpper(string(f.Name[0])) + f.Name[1:],
		FuncNameSelector:      nameSelector,
		FuncName:              f.Name,
		Schema:                f.Schema,
		Imports:               imports,
		ReturnsError:          f.ReturnsError,
		SchemaPackageSelector: utils.EnsureUniqueGoPackageSelector(imports, "github.com/conflowio/conflow/src/schema"),
	}
}
