// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generator

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/conflowio/conflow/pkg/internal/utils"

	"github.com/conflowio/conflow/pkg/schema"
)

// GenerateInterpreter generates an interpreter for the given function
func GenerateInterpreter(
	f *Function,
	pkg string,
	imports map[string]string,
) ([]byte, error) {
	params := &InterpreterTemplateParams{
		FunctionPackage: pkg,
		Name:            strings.ToUpper(string(f.Name[0])) + f.Name[1:],
		FuncName:        f.Name,
		Schema:          f.Schema,
		Imports:         imports,
		ReturnsError:    f.ReturnsError,
	}

	bodyTmpl := template.New("body")
	bodyTmpl.Funcs(map[string]interface{}{
		"assignValue": func(s schema.Schema, valueName, resultName string) string {
			return s.AssignValue(params.Imports, valueName, resultName)
		},
		"getType": func(s schema.Schema) string {
			return s.GoType(params.Imports)
		},
		"ensureUniqueGoPackageSelector": utils.EnsureUniqueGoPackageSelector,
	})
	if _, parseErr := bodyTmpl.Parse(interpreterTemplate); parseErr != nil {
		return nil, parseErr
	}

	body := &bytes.Buffer{}
	if err := bodyTmpl.Execute(body, params); err != nil {
		return nil, err
	}

	return body.Bytes(), nil
}
