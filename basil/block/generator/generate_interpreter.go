// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generator

import (
	"bytes"
	"go/ast"
	"text/template"

	"github.com/opsidian/basil/basil/variable"
)

// GenerateInterpreter generates an interpreter for the given block
func GenerateInterpreter(str *ast.StructType, file *ast.File, pkgName string, name string) ([]byte, error) {
	params, err := generateTemplateParams(str, file, pkgName, name)
	if err != nil {
		return nil, err
	}

	tmpl := template.New("block_interpreter")
	tmpl.Funcs(map[string]interface{}{
		"filterInputs":   func(fs Fields) Fields { return fs.Filter(func(f *Field) bool { return !f.IsOutput }) },
		"filterParams":   func(fs Fields) Fields { return fs.Filter(func(f *Field) bool { return !f.IsBlock }) },
		"filterBlocks":   func(fs Fields) Fields { return fs.Filter(func(f *Field) bool { return f.IsBlock }) },
		"filterNonID":    func(fs Fields) Fields { return fs.Filter(func(f *Field) bool { return !f.IsID }) },
		"filterDefaults": func(fs Fields) Fields { return fs.Filter(func(f *Field) bool { return f.Default != nil }) },
	})
	if _, parseErr := tmpl.Parse(interpreterTemplate); parseErr != nil {
		return nil, parseErr
	}

	res := &bytes.Buffer{}
	err = tmpl.Execute(res, params)
	if err != nil {
		return nil, err
	}

	return res.Bytes(), nil
}

func generateTemplateParams(str *ast.StructType, file *ast.File, pkgName string, name string) (*InterpreterTemplateParams, error) {
	fields, err := ParseFields(str, file)
	if err != nil {
		return nil, err
	}

	var idField, valueField *Field
	for _, field := range fields {
		switch {
		case field.IsID:
			idField = field
		case field.IsValue:
			valueField = field
		}
	}

	return &InterpreterTemplateParams{
		Package:            pkgName,
		Name:               name,
		Fields:             fields,
		IDField:            idField,
		ValueField:         valueField,
		ValueFunctionNames: variable.ValueFunctionNames,
	}, nil
}
