// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blocks

import (
	"context"

	"github.com/conflowio/parsley/parsley"

	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/conflow/annotations"
	"github.com/conflowio/conflow/src/schema"
)

type module struct {
	id           conflow.ID
	interpreter  conflow.BlockInterpreter
	node         parsley.Node
	params       map[conflow.ID]interface{}
	logger       conflow.Logger
	userCtx      interface{}
	jobScheduler conflow.JobScheduler
}

func (m *module) ID() conflow.ID {
	return m.id
}

func (m *module) Run(ctx context.Context) (conflow.Result, error) {
	moduleCtx, moduleCancel := context.WithCancel(ctx)
	defer moduleCancel()

	evalContext := conflow.NewEvalContext(
		moduleCtx, m.userCtx, m.logger, m.jobScheduler, nil,
	)
	evalContext.InputParams = m.params

	value, err := parsley.EvaluateNode(evalContext, m.node)
	if err != nil {
		return nil, err
	}

	s := m.interpreter.Schema().(*schema.Object)
	for jsonPropertyName, property := range s.Properties {
		if property.GetReadOnly() {
			parameterName := s.ParameterName(jsonPropertyName)
			m.params[conflow.ID(parameterName)] = m.interpreter.Param(value.(conflow.Block), conflow.ID(parameterName))
		}
	}

	return nil, nil
}

// NewModuleInterpreter creates a new interpreter for a module
func NewModuleInterpreter(
	interpreter conflow.BlockInterpreter,
	node parsley.Node,
) conflow.BlockInterpreter {
	s := interpreter.Schema().Copy().(*schema.Object)
	for _, p := range s.Properties {
		p.(schema.MetadataAccessor).SetAnnotation(annotations.UserDefined, "")
	}

	return &moduleInterpreter{
		interpreter: interpreter,
		node:        node,
		schema:      s,
	}
}

type moduleInterpreter struct {
	node        parsley.Node
	interpreter conflow.BlockInterpreter
	schema      schema.Schema
}

func (m *moduleInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	return &module{
		id:           id,
		interpreter:  m.interpreter,
		node:         m.node,
		params:       map[conflow.ID]interface{}{},
		logger:       blockCtx.Logger(),
		userCtx:      blockCtx.UserContext(),
		jobScheduler: blockCtx.JobScheduler(),
	}
}

func (m *moduleInterpreter) Schema() schema.Schema {
	return m.schema
}

func (m *moduleInterpreter) SetParam(b conflow.Block, name conflow.ID, value interface{}) error {
	b.(*module).params[name] = value
	return nil
}

func (m *moduleInterpreter) SetBlock(b conflow.Block, name conflow.ID, key string, value interface{}) error {
	return nil
}

func (m *moduleInterpreter) Param(b conflow.Block, name conflow.ID) interface{} {
	return b.(*module).params[name]
}

func (m *moduleInterpreter) ValueParamName() conflow.ID {
	return ""
}

func (m *moduleInterpreter) ParseContext(context *conflow.ParseContext) *conflow.ParseContext {
	return context
}
