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

// NewModuleInterpreter creates a new interpreter for a module
func NewModuleInterpreter(
	interpreter basil.BlockInterpreter,
	node parsley.Node,
) basil.BlockInterpreter {
	params := make(map[basil.ID]basil.ParameterDescriptor, len(interpreter.Params()))
	for k, v := range interpreter.Params() {
		v.IsUserDefined = false
		params[k] = v
	}

	return &moduleInterpreter{
		interpreter: interpreter,
		node:        node,
		params:      params,
	}
}

type moduleInterpreter struct {
	node        parsley.Node
	interpreter basil.BlockInterpreter
	params      map[basil.ID]basil.ParameterDescriptor
}

func (m *moduleInterpreter) CreateBlock(id basil.ID) basil.Block {
	return NewModule(id, m.interpreter, m.node)
}

func (m moduleInterpreter) SetParam(b basil.Block, name basil.ID, value interface{}) error {
	b.(*module).params[name] = value
	return nil
}

func (m moduleInterpreter) SetBlock(b basil.Block, name basil.ID, value interface{}) error {
	return nil
}

func (m moduleInterpreter) Param(b basil.Block, name basil.ID) interface{} {
	return b.(*module).params[name]
}

func (m *moduleInterpreter) Params() map[basil.ID]basil.ParameterDescriptor {
	return m.params
}

func (m moduleInterpreter) Blocks() map[basil.ID]basil.BlockDescriptor {
	return nil
}

func (m moduleInterpreter) ValueParamName() basil.ID {
	return ""
}

func (m moduleInterpreter) HasForeignID() bool {
	return false
}

func (m moduleInterpreter) ParseContext(context *basil.ParseContext) *basil.ParseContext {
	return context
}

func (m moduleInterpreter) EvalStage() basil.EvalStage {
	return basil.EvalStageUndefined
}
