// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package openapi

import (
	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/schema"
	"github.com/conflowio/conflow/src/schema/blocks"
)

// @block "configuration"
type Parameter struct {
	// @required
	Name string `json:"name"`
	// @enum ["query", "header", "path", "cookie"]
	In              string               `json:"in"`
	Description     string               `json:"description,omitempty"`
	Required        bool                 `json:"required,omitempty"`
	Deprecated      bool                 `json:"deprecated,omitempty"`
	AllowEmptyValue bool                 `json:"allowEmptyValue,omitempty"`
	Style           string               `json:"style,omitempty"`
	Explode         bool                 `json:"explode,omitempty"`
	AllowReserved   bool                 `json:"allowReserved,omitempty"`
	Schema          schema.Schema        `json:"schema,omitempty"`
	Content         map[string]MediaType `json:"content,omitempty"`
}

func (p *Parameter) ParseContextOverride() conflow.ParseContextOverride {
	return conflow.ParseContextOverride{
		BlockTransformerRegistry: blocks.InterpreterRegistry(),
	}
}
