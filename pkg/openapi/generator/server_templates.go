// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generator

type ServerTemplateParams struct {
	Operations []Operation
	Imports    map[string]string
}

var serverTemplate = `
{{- $root := . -}}
{{- $ctxSel := ensureUniqueGoPackageSelector .Imports "context" -}}

type Server interface {
	{{ range $op := .Operations -}}
		{{ camelize $op.OperationID }}(ctx {{ $ctxSel }}Context, req {{ camelize $op.OperationID }}Request) ({{ camelize $op.OperationID }}Response, error)
	{{ end -}}
}
`
