// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generator

import (
	"bytes"
	"text/template"

	"github.com/conflowio/conflow/pkg/conflow/annotations"
	"github.com/conflowio/conflow/pkg/schema"
)

// GenerateInterpreter generates an interpreter for the given block
func GenerateInterpreter(
	str *Struct,
	blockPackage string,
	imports map[string]string,
) ([]byte, error) {
	var idPropertyName, valuePropertyName string
	o := str.Schema.(*schema.Object)
	for jsonPropertyName, property := range o.Properties {
		parameterName := o.ParameterName(jsonPropertyName)
		switch {
		case property.GetAnnotation(annotations.ID) == "true":
			idPropertyName = parameterName
		case property.GetAnnotation(annotations.Value) == "true":
			valuePropertyName = parameterName
		}
	}

	params := &InterpreterTemplateParams{
		BlockPackage:      blockPackage,
		Name:              str.Name,
		Schema:            str.Schema,
		IDPropertyName:    idPropertyName,
		ValuePropertyName: valuePropertyName,
		Imports:           imports,
		Dependencies:      str.Dependencies,
	}

	bodyTmpl := template.New("body")
	bodyTmpl.Funcs(commonFuncs(params.Schema, params.Imports))
	if _, parseErr := bodyTmpl.Parse(interpreterTemplate); parseErr != nil {
		return nil, parseErr
	}

	body := &bytes.Buffer{}
	if err := bodyTmpl.Execute(body, params); err != nil {
		return nil, err
	}

	return body.Bytes(), nil
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
