// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generator

import (
	"github.com/opsidian/basil/basil/generator/parser"
	"github.com/opsidian/basil/basil/schema"
)

type InterpreterTemplateParams struct {
	Package           string
	NameSelector      string
	Name              string
	Schema            schema.Schema
	IDPropertyName    string
	ValuePropertyName string
	Imports           map[string]string
	Dependencies      []parser.Dependency
}

const interpreterHeaderTemplate = `
{{ $root := . -}}

// Code generated by Basil. DO NOT EDIT.

package {{ .Package }}

import (
	{{ range sortedImportKeys .Imports -}}
	{{ if ne . "." -}}
	{{ if ne (last (index $root.Imports .)) . }}{{ . }} {{ end }}{{ printf "%q" (index $root.Imports .) }}
	{{ end -}}
	{{ end -}}
)
`

const interpreterTemplate = `
{{ $root := . -}}

// {{ .Name }}Interpreter is the basil interpreter for the {{ .Name }} block
type {{ .Name }}Interpreter struct {
	s schema.Schema
}

func (i {{ .Name }}Interpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = {{ .Schema.GoString }}
	}
	return i.s
}

// Create creates a new {{ .Name }} block
func (i {{ .Name }}Interpreter) CreateBlock(id basil.ID, blockCtx *basil.BlockContext) basil.Block {
	return &{{ .NameSelector }}{{ .Name }}{
		{{ if .IDPropertyName -}}
		{{ getPropertyName .IDPropertyName }}: id,
		{{ end -}}
		{{ range $name, $schema := filterDefaults (filterParams .Schema.GetProperties) -}}
		{{ getPropertyName $name }}: {{ printf "%#v" .DefaultValue }},
		{{ end -}}
		{{ range .Dependencies -}}
		{{ .FieldName }}: blockCtx.{{ title .Name }}(),
		{{ end -}}
	}
}

// ValueParamName returns the name of the parameter marked as value field, if there is one set
func (i {{ .Name }}Interpreter) ValueParamName() basil.ID {
	return "{{ .ValuePropertyName }}"
}

// ParseContext returns with the parse context for the block
func (i {{.Name}}Interpreter) ParseContext(ctx *basil.ParseContext) *basil.ParseContext {
	var nilBlock *{{ .NameSelector }}{{.Name}}
	if b, ok := basil.Block(nilBlock).(basil.ParseContextOverrider); ok {
		return ctx.New(b.ParseContextOverride())
	}

	return ctx
}

func (i {{ .Name }}Interpreter) Param(b basil.Block, name basil.ID) interface{} {
	switch name {
	{{ range $name, $property := filterParams .Schema.GetProperties -}}
	case "{{ $name }}":
		return b.(*{{ $root.NameSelector }}{{ $root.Name }}).{{ getPropertyName $name }}
	{{ end -}}
	default:
		panic(fmt.Errorf("unexpected parameter %q in {{ .Name }}", name))
	}
}

func (i {{ .Name }}Interpreter) SetParam(block basil.Block, name basil.ID, value interface{}) error {
	{{ if filterInputs (filterParams .Schema.GetProperties) -}}
	b := block.(*{{ .NameSelector }}{{ .Name }})
	switch name {
	{{ range $name, $property := filterInputs (filterParams .Schema.GetProperties) -}}
	case "{{ $name }}":
		{{ assignValue $property "value" (printf "b.%s" (getPropertyName $name)) }}
	{{ end -}}
	}
	return nil
	{{ else -}}
	return nil
	{{ end -}}
}

func (i {{ .Name }}Interpreter) SetBlock(block basil.Block, name basil.ID, value interface{}) error {
	{{ if filterInputs (filterBlocks .Schema.GetProperties) -}}
	b := block.(*{{ $root.NameSelector }}{{ $root.Name }})
	switch name {
	{{ range $name, $property := filterInputs (filterBlocks .Schema.GetProperties) -}}
	case "{{ $name }}":
		{{ if isArray $property -}}
		b.{{ getPropertyName $name }} = append(b.{{ getPropertyName $name }}, value.({{ getType $property.GetItems }}))
		{{ else -}}
		b.{{ getPropertyName $name }} = value.({{ getType $property }})
		{{ end -}}
	{{ end -}}
	}
	{{ end -}}
	return nil
}
`
