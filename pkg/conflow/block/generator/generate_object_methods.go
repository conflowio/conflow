// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generator

import (
	"bytes"
	"text/template"
)

func GenerateObjectMethods(
	str *Struct,
	imports map[string]string,
) ([]byte, error) {
	params := &ObjectMethodsTemplateParams{
		Name:    str.Name,
		Schema:  str.Schema,
		Imports: imports,
	}

	bodyTmpl := template.New("body")
	bodyTmpl.Funcs(commonFuncs(params.Schema, params.Imports))
	if _, parseErr := bodyTmpl.Parse(objectMethodsTemplate); parseErr != nil {
		return nil, parseErr
	}

	body := &bytes.Buffer{}
	if err := bodyTmpl.Execute(body, params); err != nil {
		return nil, err
	}

	return body.Bytes(), nil
}
