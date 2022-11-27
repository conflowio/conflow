// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"fmt"
	"sort"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/schema"
	schemainterpreters "github.com/conflowio/conflow/pkg/schema/interpreters"
)

// @block "directive"
type OneOf struct {
	// @id
	id conflow.ID
	// @name "schema"
	// @value
	schemas []schema.Schema
}

func (o *OneOf) ParseContextOverride() conflow.ParseContextOverride {
	return conflow.ParseContextOverride{
		BlockTransformerRegistry: schemainterpreters.Registry(),
	}
}

func (o *OneOf) ID() conflow.ID {
	return o.id
}

func (o *OneOf) ApplyToSchema(s schema.Schema) error {
	return nil
}

func (o *OneOf) ReplaceSchema(s schema.Schema) (schema.Schema, error) {
	if _, ok := s.(*schema.Any); !ok {
		return nil, fmt.Errorf("@one_of can only be used on an interface{} type")
	}

	sort.Slice(o.schemas, func(i, j int) bool {
		return o.schemas[i].TypeString() < o.schemas[j].TypeString()
	})

	return &schema.OneOf{Schemas: o.schemas}, nil
}
