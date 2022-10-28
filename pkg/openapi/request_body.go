// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package openapi

import (
	"strings"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/conflow/block"
	"github.com/conflowio/conflow/pkg/schema"
	"github.com/conflowio/conflow/pkg/util"
)

// @block "configuration"
type RequestBody struct {
	Description string `json:"description,omitempty"`
	// @required
	Content  map[string]*MediaType `json:"content"`
	Required bool                  `json:"required,omitempty"`
}

func (r *RequestBody) ParseContextOverride() conflow.ParseContextOverride {
	return conflow.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"content": MediaTypeInterpreter{},
		},
	}
}

func (r *RequestBody) Validate(ctx *schema.Context) error {
	return schema.ValidateMap("content", r.Content)(ctx)
}

func (r *RequestBody) GoString(imports map[string]string) string {
	pkg := openapiPkg(imports)
	b := &strings.Builder{}
	fprintf(b, "&%sRequestBody{\n", pkg)
	if r.Description != "" {
		fprintf(b, "\tDescription: %#v,\n", r.Description)
	}
	fprintf(b, "\tContent: map[string]*%sMediaType{\n", pkg)
	for _, k := range util.SortedMapKeys(r.Content) {
		fprintf(b, "\t\t%#v: %s,\n", k, indent(indent(r.Content[k].GoString(imports))))
	}
	fprintf(b, "\t},\n")
	if r.Required {
		fprintf(b, "\tRequired: true,\n")
	}
	b.WriteRune('}')
	return b.String()
}
