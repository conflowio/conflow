// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generator

import (
	"github.com/conflowio/conflow/pkg/schema"
)

type ObjectMethodsTemplateParams struct {
	Name    string
	Schema  schema.Schema
	Imports map[string]string
}

const objectMethodsTemplate = `
{{ $root := . -}}
{{- $schemaSel := ensureUniqueGoPackageSelector .Imports "github.com/conflowio/conflow/pkg/schema" -}}


func init() {
	{{ $schemaSel }}Register({{ .Schema.GoString .Imports }})
}

// New{{ .Name }}WithDefaults creates a new {{ .Name }} instance with default values
func New{{ .Name }}WithDefaults() *{{ .Name }} {
	b := &{{ .Name }}{}
	{{ range $name, $schema := filterDefaults (filterParams .Schema.Properties) -}}
	b.{{ getFieldName $name }} = {{ printf "%#v" .DefaultValue }}
	{{ end -}}
	{{ range $name, $property := filterInputs (filterBlocks .Schema.Properties) -}}
	{{ if isMap $property -}}
	b.{{ getFieldName $name }} = map[string]{{ getType $property.AdditionalProperties }}{}
	{{ end -}}
	{{ end -}}
	return b
}
`
