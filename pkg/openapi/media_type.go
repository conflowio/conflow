// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package openapi

import (
	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/schema"
	schemainterpreters "github.com/conflowio/conflow/pkg/schema/interpreters"
)

// @block "configuration"
type MediaType struct {
	// @required
	Schema schema.Schema `json:"schema,omitempty"`
}

func (m *MediaType) ParseContextOverride() conflow.ParseContextOverride {
	return conflow.ParseContextOverride{
		BlockTransformerRegistry: schemainterpreters.Registry(),
	}
}

func (m *MediaType) Validate(ctx *schema.Context) error {
	return schema.Validate("schema", m.Schema)(ctx)
}
