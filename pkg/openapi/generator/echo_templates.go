// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generator

import "github.com/conflowio/conflow/pkg/openapi"

type EchoServerTemplateParams struct {
	Operations  []Operation
	Imports     map[string]string
	RequestBody *openapi.RequestBody
}

var echoServerTemplate = `
{{- $root := . -}}
{{- $echoSel := ensureUniqueGoPackageSelector .Imports "github.com/labstack/echo/v4" -}}
{{- $fmtSel := ensureUniqueGoPackageSelector .Imports "fmt" -}}
{{- $serverSel := ensureUniqueGoPackageSelector .Imports "github.com/conflowio/conflow/pkg/openapi/server" -}}
{{- $stringsSel := ensureUniqueGoPackageSelector .Imports "strings" -}}
{{- $httpSel := ensureUniqueGoPackageSelector .Imports "net/http" -}}
type EchoServer struct {
	Server Server
}
{{ range $op := .Operations -}}
func (e *EchoServer) {{ camelize $op.OperationID }}(ctx {{ $echoSel }}Context) error {
	req := {{ camelize $op.OperationID }}Request{}
	
	{{ if $op.RequestBody }}{{ if $op.RequestBody.Content -}}
	contentType := ctx.Request().Header.Get({{ $echoSel }}HeaderContentType)
	switch {
	{{ range $contentType := (sortedMapKeys $op.RequestBody.Content) }}
	{{- $content := (index $op.RequestBody.Content $contentType) -}}
	case {{ $stringsSel }}HasPrefix(contentType, {{ printf "%q" $contentType }}):
		if err := {{ $serverSel }}BindBody[{{ schemaType $content.Schema $root.Imports }}](
			{{ $op.RequestBody.GoString $root.Imports }},
			ctx,
			&req.Body{{ contentTypeName $contentType }},
		); err != nil {
				return err
		}
	{{ end -}}
	default:
		return {{ $serverSel }}NewHTTPErrorf({{ $httpSel }}StatusBadRequest, "unsupported content type: '%s'", contentType)
	}

	{{ end }}{{ end -}}
	
	{{ range $field, $p := $op.Parameters -}}
	if err := {{ $serverSel }}BindParameter[{{ schemaType $p.Schema $root.Imports }}](
		{{ $p.GoString $root.Imports }},
		ctx,
		&req.{{ $field }},
	); err != nil {
		return err
	}
	
	{{- end }}

	resp, err := e.Server.{{ camelize $op.OperationID }}(ctx.Request().Context(), req)
	if err != nil {
		return err
	}
	if r, ok := resp.({{ camelize $op.OperationID }}Response); ok {
		return r.Write{{ camelize $op.OperationID }}Response(ctx.Response())
	}
	return {{ $fmtSel }}Errorf("unexpected response type: %T", resp)
}

{{ end -}}

// EchoRouter defines a common interface for echo.Echo and echo.Group
type EchoRouter interface {
	CONNECT(string, {{ $echoSel }}HandlerFunc, ...{{ $echoSel }}MiddlewareFunc) *{{ $echoSel }}Route
	DELETE(string, {{ $echoSel }}HandlerFunc, ...{{ $echoSel }}MiddlewareFunc) *{{ $echoSel }}Route
	GET(string, {{ $echoSel }}HandlerFunc, ...{{ $echoSel }}MiddlewareFunc) *{{ $echoSel }}Route
	HEAD(string, {{ $echoSel }}HandlerFunc, ...{{ $echoSel }}MiddlewareFunc) *{{ $echoSel }}Route
	OPTIONS(string, {{ $echoSel }}HandlerFunc, ...{{ $echoSel }}MiddlewareFunc) *{{ $echoSel }}Route
	PATCH(string, {{ $echoSel }}HandlerFunc, ...{{ $echoSel }}MiddlewareFunc) *{{ $echoSel }}Route
	POST(string, {{ $echoSel }}HandlerFunc, ...{{ $echoSel }}MiddlewareFunc) *{{ $echoSel }}Route
	PUT(string, {{ $echoSel }}HandlerFunc, ...{{ $echoSel }}MiddlewareFunc) *{{ $echoSel }}Route
	TRACE(string, {{ $echoSel }}HandlerFunc, ...{{ $echoSel }}MiddlewareFunc) *{{ $echoSel }}Route
}

// RegisterEchoHandlers adds each server route to the Echo router
func RegisterEchoHandlers(r EchoRouter, s Server) {
	RegisterEchoHandlersWithBaseURL(r, s, "")
}

// RegisterEchoHandlersWithBaseURL adds each server route to the Echo router with a base URL
func RegisterEchoHandlersWithBaseURL(r EchoRouter, s Server, u string) {
	e := &EchoServer{Server: s}
	{{ range $op := .Operations -}}
		r.{{ $op.Method }}(u+{{ printf "%q" (convertPath $op.Path) }}, e.{{ camelize $op.OperationID }})
	{{ end -}}
}
`
