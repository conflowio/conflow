// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generator

import "github.com/conflowio/conflow/src/schema"

type InterpreterTemplateParams struct {
	Package               string
	Name                  string
	FuncNameSelector      string
	FuncName              string
	Schema                schema.Schema
	Imports               map[string]string
	ReturnsError          bool
	SchemaPackageSelector string
}

const interpreterTemplate = `
{{ $root := . -}}

func init() {
	{{ .SchemaPackageSelector }}Register({{ .Schema.GoString .Imports }})
}

// {{ .Name }}Interpreter is the Conflow interpreter for the {{ .FuncName }} function
type {{ .Name }}Interpreter struct {
}

func (i {{ .Name }}Interpreter) Schema() {{ .SchemaPackageSelector }}Schema {
	s, _ := {{ .SchemaPackageSelector }}Get("{{ .Schema.ID }}")
	return s
}

// Eval returns with the result of the function
func (i {{ .Name }}Interpreter) Eval(ctx interface{}, args []interface{}) (interface{}, error) {
	{{ range $i, $property := .Schema.GetParameters -}}
	var {{ assignValue $property.Schema (printf "args[%d]" $i) (printf "val%d" $i) }}
	{{ end -}}
	{{ if .Schema.GetAdditionalParameters -}}
	var variadicArgs []{{ getType .Schema.GetAdditionalParameters.Schema }}
	for p := {{ len .Schema.GetParameters }}; p < len(args); p++ {
		var {{ assignValue .Schema.GetAdditionalParameters.Schema "args[p]" "val" }}
		variadicArgs = append(variadicArgs, val)
	}
	{{ end -}}
	return {{ .FuncNameSelector }}{{ .FuncName }}(
		{{- range $i, $property := .Schema.GetParameters -}}
		val{{ $i }},
		{{- end -}}
		{{- if .Schema.GetAdditionalParameters -}}
		variadicArgs...,
		{{- end -}}
	){{ if not .ReturnsError }}, nil{{ end }}
}
`
