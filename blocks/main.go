// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blocks

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/parsley/parsley"
)

type Main struct {
}

func (m Main) ID() basil.ID {
	return "main"
}

type MainInterpreter struct {
	BlockTransformerRegistry    parsley.NodeTransformerRegistry
	FunctionTransformerRegistry parsley.NodeTransformerRegistry
}

func (m MainInterpreter) CreateBlock(basil.ID) basil.Block {
	return Main{}
}

func (m MainInterpreter) SetParam(b basil.Block, name basil.ID, value interface{}) error {
	return nil
}

func (m MainInterpreter) SetBlock(b basil.Block, name basil.ID, value interface{}) error {
	return nil
}

func (m MainInterpreter) Param(b basil.Block, name basil.ID) interface{} {
	return nil
}

func (m MainInterpreter) Params() map[basil.ID]basil.ParameterDescriptor {
	return nil
}

func (m MainInterpreter) Blocks() map[basil.ID]basil.BlockDescriptor {
	return nil
}

func (m MainInterpreter) ValueParamName() basil.ID {
	return ""
}

func (m MainInterpreter) HasForeignID() bool {
	return false
}

func (m MainInterpreter) ParseContext(ctx *basil.ParseContext) *basil.ParseContext {
	return ctx.New(basil.ParseContextOverride{
		BlockTransformerRegistry:    m.BlockTransformerRegistry,
		FunctionTransformerRegistry: m.FunctionTransformerRegistry,
	})
}
