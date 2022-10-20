// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generator

type ResponseTemplateParams struct {
	OperationID string
	Imports     map[string]string
}

var responseTemplate = `
{{- $httpSel := ensureUniqueGoPackageSelector .Imports "net/http" -}}

type {{ camelize .OperationID }}Response interface {
	Write{{ camelize .OperationID }}Response(w {{ $httpSel }}ResponseWriter) error
}

type {{ camelize .OperationID }}ResponseFunc func(w {{ $httpSel }}ResponseWriter) error

func (f {{ camelize .OperationID }}ResponseFunc) Write{{ camelize .OperationID }}Response(w {{ $httpSel }}ResponseWriter) error {
	return f(w)
}
`

type WriteResponseFuncTemplateParams struct {
	OperationID  string
	ResponseCode string
	ContentTypes []string
	Imports      map[string]string
}

var writeResponseFuncTemplate = `
{{- $root := . -}}
{{- $httpSel := ensureUniqueGoPackageSelector .Imports "net/http" -}}

func (r {{ camelize .OperationID }}Response{{ if eq .ResponseCode "default" }}Default{{ else }}{{ .ResponseCode }}{{ end }}) Write{{ camelize .OperationID }}Response(w {{ $httpSel }}ResponseWriter) error {
	{{ if not (eq .ResponseCode "default") -}}
	w.WriteHeader({{ .ResponseCode }})
	{{ else -}}
	w.WriteHeader(int(r.ResponseCode))
	{{ end -}}

	{{ if .ContentTypes -}}
	switch {
		{{ range $contentType := .ContentTypes -}}
		case r.Body{{ contentTypeName $contentType }} != nil:
			return {{ ensureUniqueGoPackageSelector $root.Imports "github.com/conflowio/conflow/src/openapi/server" }}WriteResponse(w, {{ printf "%q" $contentType }}, *r.Body{{ contentTypeName $contentType }})
		{{ end -}}
	}
	{{ end -}}
	return nil
}
`
