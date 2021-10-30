// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blocks

import (
	"context"

	"github.com/opsidian/parsley/parsley"

	"github.com/opsidian/conflow/basil"
	"github.com/opsidian/conflow/basil/schema"
)

type module struct {
	id           basil.ID
	interpreter  basil.BlockInterpreter
	node         parsley.Node
	params       map[basil.ID]interface{}
	logger       basil.Logger
	userCtx      interface{}
	jobScheduler basil.JobScheduler
}

func (m *module) ID() basil.ID {
	return m.id
}

func (m *module) Run(ctx context.Context) (basil.Result, error) {
	moduleCtx, moduleCancel := context.WithCancel(ctx)
	defer moduleCancel()

	evalContext := basil.NewEvalContext(
		moduleCtx, m.userCtx, m.logger, m.jobScheduler, nil,
	)
	evalContext.InputParams = m.params

	value, err := parsley.EvaluateNode(evalContext, m.node)
	if err != nil {
		return nil, err
	}

	for propertyName, property := range m.interpreter.Schema().(schema.ObjectKind).GetProperties() {
		if property.GetReadOnly() {
			m.params[basil.ID(propertyName)] = m.interpreter.Param(value.(basil.Block), basil.ID(propertyName))
		}
	}

	return nil, nil
}

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

func (m *moduleInterpreter) CreateBlock(id basil.ID, blockCtx *basil.BlockContext) basil.Block {
	return &module{
		id:           id,
		interpreter:  m.interpreter,
		node:         m.node,
		params:       map[basil.ID]interface{}{},
		logger:       blockCtx.Logger(),
		userCtx:      blockCtx.UserContext(),
		jobScheduler: blockCtx.JobScheduler(),
	}
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
