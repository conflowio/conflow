// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package openapi

import (
	"context"
	"errors"
	"strings"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/schema"
	schemainterpreters "github.com/conflowio/conflow/pkg/schema/interpreters"
	"github.com/conflowio/conflow/pkg/util/ptr"
	"github.com/conflowio/conflow/pkg/util/validation"
)

// @block "configuration"
type Parameter struct {
	// @required
	Name string `json:"name"`
	// $required
	// @enum ["cookie", "header", "path", "query"]
	In          string `json:"in"`
	Description string `json:"description,omitempty"`
	Required    *bool  `json:"required,omitempty"`
	Deprecated  bool   `json:"deprecated,omitempty"`
	// @enum ["deepObject", "form", "label", "matrix", "pipeDelimited", "simple", "spaceDelimited"]
	Style   string `json:"style,omitempty"`
	Explode *bool  `json:"explode,omitempty"`
	// @required
	Schema schema.Schema `json:"schema,omitempty"`
}

func (p *Parameter) ParseContextOverride() conflow.ParseContextOverride {
	return conflow.ParseContextOverride{
		BlockTransformerRegistry: schemainterpreters.Registry(),
	}
}

func (p *Parameter) Validate(ctx context.Context) error {
	switch p.In {
	case "path":
		if p.Required != nil && !*p.Required {
			return validation.NewFieldError("required", errors.New("false is not allowed on path parameters"))
		}
		p.Required = ptr.To(true)

		if p.Style != "" && p.Style != "simple" {
			return validation.NewFieldError("style", errors.New("only 'simple' is supported on path parameters"))
		}
		if ptr.Value(p.Explode) {
			return validation.NewFieldError("explode", errors.New("true is not supported on path parameters"))
		}
	case "query":
		if p.Style != "" && p.Style != "form" {
			return validation.NewFieldError("style", errors.New("only 'form' is supported on query parameters"))
		}
		if p.Explode != nil && !*p.Explode {
			return validation.NewFieldError("explode", errors.New("false is not supported on query parameters"))
		}
	case "header":
		if p.Style != "" && p.Style != "simple" {
			return validation.NewFieldError("style", errors.New("only 'simple' is supported on header parameters"))
		}
		if ptr.Value(p.Explode) {
			return validation.NewFieldError("explode", errors.New("true is not supported on header parameters"))
		}
	case "cookie":
		if p.Style != "" && p.Style != "form" {
			return validation.NewFieldError("style", errors.New("only 'form' is supported on cookie parameters"))
		}
		if ptr.Value(p.Explode) {
			return validation.NewFieldError("explode", errors.New("true is not supported on cookie parameters"))
		}
	}

	return nil
}

func (p *Parameter) GoString(imports map[string]string) string {
	pkg := openapiPkg(imports)
	schemaPkg := schemaPkg(imports)
	b := &strings.Builder{}
	fprintf(b, "&%sParameter{\n", pkg)
	fprintf(b, "\tName: %#v,\n", p.Name)
	fprintf(b, "\tIn: %#v,\n", p.In)
	if p.Description != "" {
		fprintf(b, "\tDescription: %#v,\n", p.Description)
	}
	if p.Required != nil {
		fprintf(b, "\tRequired: %sPointer(%#v),\n", schemaPkg, *p.Required)
	}
	if p.Deprecated {
		fprintf(b, "\tDeprecated: true,\n")
	}
	if p.Style != "" {
		fprintf(b, "\tStyle: %#v,\n", p.Style)
	}
	if p.Explode != nil {
		fprintf(b, "\tExplode: %sPointer(%#v),\n", schemaPkg, *p.Explode)
	}
	fprintf(b, "\tSchema: %s,\n", indent(p.Schema.GoString(imports)))
	b.WriteRune('}')
	return b.String()
}
