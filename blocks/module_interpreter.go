// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blocks

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/schema"
	"github.com/opsidian/parsley/parsley"
)

// NewModuleInterpreter creates a new interpreter for a module
func NewModuleInterpreter(
	interpreter basil.BlockInterpreter,
	node parsley.Node,
) basil.BlockInterpreter {
	s := interpreter.Schema().Copy().(*schema.Object)
	for _, p := range s.Properties {
		p.(schema.MetadataAccessor).SetAnnotation("user_defined", "")
	}

	return &moduleInterpreter{
		interpreter: interpreter,
		node:        node,
		schema:      s,
	}
}

type moduleInterpreter struct {
	node        parsley.Node
	interpreter basil.BlockInterpreter
	schema      schema.Schema
}

func (m *moduleInterpreter) CreateBlock(id basil.ID) basil.Block {
	return NewModule(id, m.interpreter, m.node)
}

func (m *moduleInterpreter) Schema() schema.Schema {
	return m.schema
}

func (m *moduleInterpreter) SetParam(b basil.Block, name basil.ID, value interface{}) error {
	b.(*module).params[name] = value
	return nil
}

func (m *moduleInterpreter) SetBlock(b basil.Block, name basil.ID, value interface{}) error {
	return nil
}

func (m *moduleInterpreter) Param(b basil.Block, name basil.ID) interface{} {
	return b.(*module).params[name]
}

func (m *moduleInterpreter) ValueParamName() basil.ID {
	return ""
}

func (m *moduleInterpreter) ParseContext(context *basil.ParseContext) *basil.ParseContext {
	return context
}
