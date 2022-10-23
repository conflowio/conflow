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
type Response struct {
	// @required
	Description string                `json:"description,omitempty"`
	Content     map[string]*MediaType `json:"content,omitempty"`
}

func (r *Response) ParseContextOverride() conflow.ParseContextOverride {
	return conflow.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"content": MediaTypeInterpreter{},
		},
	}
}

func (r *Response) Validate(ctx *schema.Context) error {
	return schema.ValidateMap("content", r.Content)(ctx)
}
