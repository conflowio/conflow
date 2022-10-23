// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package openapi

import (
	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/conflow/block"
	"github.com/conflowio/conflow/pkg/schema"
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