// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blocks

import (
	"github.com/conflowio/parsley/parsley"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/conflow/annotations"
	"github.com/conflowio/conflow/pkg/schema"
)

type Main struct {
	params map[conflow.ID]interface{}
}

func (m *Main) ID() conflow.ID {
	return "main"
}

type MainInterpreter struct {
	BlockTransformerRegistry    parsley.NodeTransformerRegistry
	FunctionTransformerRegistry parsley.NodeTransformerRegistry
}

func (m MainInterpreter) Schema() schema.Schema {
	return &schema.Object{
		Metadata: schema.Metadata{
			Annotations: map[string]string{
				annotations.Type: conflow.BlockTypeMain,
			},
		},
	}
}

func (m MainInterpreter) CreateBlock(conflow.ID, *conflow.BlockContext) conflow.Block {
	return &Main{
		params: map[conflow.ID]interface{}{},
	}
}

func (m MainInterpreter) SetParam(b conflow.Block, name conflow.ID, value interface{}) error {
	b.(*Main).params[name] = value
	return nil
}

func (m MainInterpreter) SetBlock(b conflow.Block, name conflow.ID, key string, value interface{}) error {
	return nil
}

func (m MainInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	return b.(*Main).params[name]
}

func (m MainInterpreter) ValueParamName() conflow.ID {
	return ""
}

func (m MainInterpreter) ParseContext(ctx *conflow.ParseContext) *conflow.ParseContext {
	return ctx.New(conflow.ParseContextOverride{
		BlockTransformerRegistry:    m.BlockTransformerRegistry,
		FunctionTransformerRegistry: m.FunctionTransformerRegistry,
	})
}
