// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generator

import "github.com/conflowio/conflow/pkg/schema"

type schemasTemplateParams struct {
	Schemas []schema.Schema
	Imports map[string]string
}

var schemaTemplate = `
{{- $root := . -}}
{{- $schemaSel := ensureUniqueGoPackageSelector .Imports "github.com/conflowio/conflow/pkg/schema" -}}
func init() {
{{ range $schema := .Schemas -}}
	{{ $schemaSel }}Register({{ $schema.GoString $root.Imports }})
{{ end -}}
}
`
