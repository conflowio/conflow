// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generator

import (
	"github.com/conflowio/conflow/pkg/conflow/generator/parser"
	"github.com/conflowio/conflow/pkg/schema"
)

type InterpreterTemplateParams struct {
	BlockPackage      string
	Name              string
	Schema            schema.Schema
	IDPropertyName    string
	ValuePropertyName string
	Imports           map[string]string
	Dependencies      []parser.Dependency
}

const interpreterTemplate = `
{{ $root := . -}}
{{- $blockPackageSel := ensureUniqueGoPackageSelector .Imports .BlockPackage -}}
{{- $fmtSel := ensureUniqueGoPackageSelector .Imports "fmt" -}}
{{- $conflowSel := ensureUniqueGoPackageSelector .Imports "github.com/conflowio/conflow/pkg/conflow" -}}
{{- $schemaSel := ensureUniqueGoPackageSelector .Imports "github.com/conflowio/conflow/pkg/schema" -}}

// {{ .Name }}Interpreter is the Conflow interpreter for the {{ .Name }} block
type {{ .Name }}Interpreter struct {
}

func (i {{ .Name }}Interpreter) Schema() {{ $schemaSel }}Schema {
	s, _ := {{ $schemaSel }}Get("{{ .Schema.ID }}")
	return s
}

// Create creates a new {{ .Name }} block
func (i {{ .Name }}Interpreter) CreateBlock(id {{ $conflowSel }}ID, blockCtx *{{ $conflowSel }}BlockContext) {{ $conflowSel }}Block {
	b := {{ $blockPackageSel }}New{{ .Name }}WithDefaults()
	{{ if .IDPropertyName -}}
	b.{{ getFieldName .IDPropertyName }} = id
	{{ end -}}
	{{ range .Dependencies -}}
	b.{{ .FieldName }} = blockCtx.{{ title .Name }}()
	{{ end -}}
	return b
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i {{ .Name }}Interpreter) ValueParamName() {{ $conflowSel }}ID {
	return "{{ .ValuePropertyName }}"
}

// ParseContext returns with the parse context for the block
func (i {{.Name}}Interpreter) ParseContext(ctx *{{ $conflowSel }}ParseContext) *{{ $conflowSel }}ParseContext {
	var nilBlock *{{ $blockPackageSel }}{{.Name}}
	if b, ok := {{ $conflowSel }}Block(nilBlock).({{ $conflowSel }}ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i {{ .Name }}Interpreter) Param(b {{ $conflowSel }}Block, name {{ $conflowSel }}ID) interface{} {
	switch name {
	{{ range $name, $property := filterParams .Schema.Properties -}}
	case "{{ getParameterName $name }}":
		return b.(*{{ $blockPackageSel }}{{ $root.Name }}).{{ getFieldName $name }}
	{{ end -}}
	default:
		panic({{ $fmtSel }}Errorf("unexpected parameter %q in {{ .Name }}", name))
	}
}

func (i {{ .Name }}Interpreter) SetParam(block {{ $conflowSel }}Block, name {{ $conflowSel }}ID, value interface{}) error {
	{{ if filterInputs (filterParams .Schema.Properties) -}}
	b := block.(*{{ $blockPackageSel }}{{ .Name }})
	switch name {
	{{ range $name, $property := filterInputs (filterParams .Schema.Properties) -}}
	case "{{ getParameterName $name }}":
		{{ assignValue $property "value" (printf "b.%s" (getFieldName $name)) }}
	{{ end -}}
	}
	return nil
	{{ else -}}
	return nil
	{{ end -}}
}

func (i {{ .Name }}Interpreter) SetBlock(block {{ $conflowSel }}Block, name {{ $conflowSel }}ID, key string, value interface{}) error {
	{{ if filterInputs (filterBlocks .Schema.Properties) -}}
	b := block.(*{{ $blockPackageSel }}{{ $root.Name }})
	switch name {
	{{ range $name, $property := filterInputs (filterBlocks .Schema.Properties) -}}
	case "{{ getParameterName $name }}":
		{{ if isArray $property -}}
		b.{{ getFieldName $name }} = append(b.{{ getFieldName $name }}, value.({{ getType $property.GetItems }}))
		{{ else if isMap $property -}}
		b.{{ getFieldName $name }}[key] = value.({{ getType $property.GetAdditionalProperties }})
		{{ else -}}
		b.{{ getFieldName $name }} = value.({{ getType $property }})
		{{ end -}}
	{{ end -}}
	}
	{{ end -}}
	return nil
}
`
