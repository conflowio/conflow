// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package openapi

import (
	"errors"
	"fmt"
	"strings"

	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/schema"
	schemainterpreters "github.com/conflowio/conflow/src/schema/interpreters"
	"github.com/conflowio/conflow/src/util/ptr"
)

// @block "configuration"
type Parameter struct {
	// @required
	Name string `json:"name"`
	// $required
	// @enum ["query", "header", "path", "cookie"]
	In            string `json:"in"`
	Description   string `json:"description,omitempty"`
	Required      *bool  `json:"required,omitempty"`
	Deprecated    bool   `json:"deprecated,omitempty"`
	Style         string `json:"style,omitempty"`
	Explode       *bool  `json:"explode,omitempty"`
	AllowReserved bool   `json:"allowReserved,omitempty"`
	// @required
	Schema schema.Schema `json:"schema,omitempty"`
}

func (p *Parameter) ParseContextOverride() conflow.ParseContextOverride {
	return conflow.ParseContextOverride{
		BlockTransformerRegistry: schemainterpreters.Registry(),
	}
}

func (p *Parameter) Validate(ctx *schema.Context) error {
	return schema.ValidateAll(
		ctx,
		func(ctx *schema.Context) error {
			switch p.In {
			case "path":
				if p.Required != nil && !*p.Required {
					return fmt.Errorf("required=false is not allowed on path parameters")
				}
				p.Required = ptr.To(true)

				if p.Style != "" && p.Style != "simple" {
					return errors.New("only style=simple is supported on path parameters")
				}
				if ptr.Value(p.Explode) {
					return errors.New("explode=true is not supported on path parameters")
				}
			case "query":
				if p.Style != "" && p.Style != "form" {
					return errors.New("only style=form is supported on query parameters")
				}
				if p.Explode != nil && !*p.Explode {
					return errors.New("explode=false is not supported on query parameters")
				}
			case "header":
				if p.Style != "" && p.Style != "simple" {
					return errors.New("only style=simple is allowed on header parameters")
				}
				if ptr.Value(p.Explode) {
					return errors.New("explode=true is not supported on header parameters")
				}
			case "cookie":
				if p.Style != "" && p.Style != "form" {
					return errors.New("only style=form is allowed on header parameters")
				}
				if ptr.Value(p.Explode) {
					return errors.New("explode=true is not supported on header parameters")
				}
			}

			return nil
		},
		schema.Validate("schema", p.Schema),
	)
}

func (p *Parameter) GoString(imports map[string]string, minimal bool) string {
	pkg := openapiPkg(imports)
	schemaPkg := schemaPkg(imports)
	b := &strings.Builder{}
	fprintf(b, "&%sParameter{\n", pkg)
	fprintf(b, "\tName: %#v,\n", p.Name)
	fprintf(b, "\tIn: %#v,\n", p.In)
	if !minimal && p.Description != "" {
		fprintf(b, "\tDescription: %#v,\n", p.Description)
	}
	if p.Required != nil {
		fprintf(b, "\tRequired: %sPointer(%#v),\n", schemaPkg, *p.Required)
	}
	if !minimal && p.Deprecated {
		fprintf(b, "\tDeprecated: true,\n")
	}
	if p.Style != "" {
		fprintf(b, "\tStyle: %#v,\n", p.Style)
	}
	if p.Explode != nil {
		fprintf(b, "\tExplode: %sPointer(%#v),\n", schemaPkg, *p.Explode)
	}
	if p.AllowReserved {
		fprintf(b, "\tAllowReserved: true,\n")
	}
	fprintf(b, "\tSchema: %s,\n", indent(p.Schema.GoString(imports)))
	b.WriteRune('}')
	return b.String()
}
