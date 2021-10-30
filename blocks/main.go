// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blocks

import (
	"github.com/opsidian/parsley/parsley"

	"github.com/opsidian/conflow/basil"
	"github.com/opsidian/conflow/basil/schema"
)

type Main struct {
	params map[basil.ID]interface{}
}

func (m *Main) ID() basil.ID {
	return "main"
}

type MainInterpreter struct {
	BlockTransformerRegistry    parsley.NodeTransformerRegistry
	FunctionTransformerRegistry parsley.NodeTransformerRegistry
}

func (m MainInterpreter) Schema() schema.Schema {
	return &schema.Object{}
}

func (m MainInterpreter) CreateBlock(basil.ID, *basil.BlockContext) basil.Block {
	return &Main{
		params: map[basil.ID]interface{}{},
	}
}

func (m MainInterpreter) SetParam(b basil.Block, name basil.ID, value interface{}) error {
	b.(*Main).params[name] = value
	return nil
}

func (m MainInterpreter) SetBlock(b basil.Block, name basil.ID, value interface{}) error {
	return nil
}

func (m MainInterpreter) Param(b basil.Block, name basil.ID) interface{} {
	return b.(*Main).params[name]
}

func (m MainInterpreter) ValueParamName() basil.ID {
	return ""
}

func (m MainInterpreter) ParseContext(ctx *basil.ParseContext) *basil.ParseContext {
	return ctx.New(basil.ParseContextOverride{
		BlockTransformerRegistry:    m.BlockTransformerRegistry,
		FunctionTransformerRegistry: m.FunctionTransformerRegistry,
	})
}
