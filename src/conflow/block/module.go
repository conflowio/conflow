// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package block

import (
	"context"
	"fmt"

	"github.com/conflowio/conflow/src/conflow/job"

	"github.com/conflowio/conflow/src/util"

	"github.com/conflowio/parsley/parsley"

	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/schema"
)

type module struct {
	id       conflow.ID
	node     conflow.BlockNode
	params   map[conflow.ID]interface{}
	blockCtx *conflow.BlockContext
}

func (m *module) ID() conflow.ID {
	return m.id
}

func (m *module) Run(ctx context.Context) (conflow.Result, error) {
	moduleCtx, moduleCancel := context.WithCancel(ctx)
	defer moduleCancel()

	evalContext := conflow.NewEvalContext(
		moduleCtx, m.blockCtx.UserContext(), m.blockCtx.Logger(), m.blockCtx.JobScheduler(), nil,
	)
	evalContext.InputParams = m.params

	value, err := parsley.EvaluateNode(evalContext, m.node)
	if err != nil {
		return nil, err
	}

	for propertyName, property := range m.node.Schema().(schema.ObjectKind).GetParameters() {
		if property.GetReadOnly() {
			m.params[conflow.ID(propertyName)] = m.node.Interpreter().Param(value.(conflow.Block), conflow.ID(propertyName))
		}
	}

	return nil, nil
}

// NewModuleInterpreter creates a new interpreter for a module
func NewModuleInterpreter(node conflow.BlockNode) (conflow.BlockInterpreter, parsley.Error) {
	schema, err := getModuleSchema(node.Children(), node.Interpreter())
	if err != nil {
		return nil, err
	}

	return &moduleInterpreter{
		node:   node,
		schema: schema,
	}, nil
}

type moduleInterpreter struct {
	node   conflow.BlockNode
	schema schema.Schema
}

func (m *moduleInterpreter) CreateBlock(id conflow.ID, blockCtx *conflow.BlockContext) conflow.Block {
	return &module{
		id:       id,
		node:     m.node,
		params:   map[conflow.ID]interface{}{},
		blockCtx: blockCtx,
	}
}

func (m *moduleInterpreter) Schema() schema.Schema {
	return m.schema
}

func (m *moduleInterpreter) SetParam(b conflow.Block, name conflow.ID, value interface{}) error {
	b.(*module).params[name] = value
	return nil
}

func (m *moduleInterpreter) SetBlock(b conflow.Block, name conflow.ID, value interface{}) error {
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

func getModuleSchema(children []conflow.Node, interpreter conflow.BlockInterpreter) (schema.Schema, parsley.Error) {
	var s schema.Schema

	for _, c := range children {
		paramNode, ok := c.(conflow.ParameterNode)
		if !ok {
			continue
		}

		config, err := getParameterConfig(paramNode)
		if err != nil {
			return nil, err
		}

		if !util.BoolValue(config.Input) && !util.BoolValue(config.Output) {
			continue
		}

		if _, exists := interpreter.Schema().(schema.ObjectKind).GetParameters()[string(paramNode.Name())]; exists {
			return nil, parsley.NewErrorf(c.Pos(), "%q parameter already exists.", paramNode.Name())
		}

		if config.Schema == nil {
			return nil, parsley.NewErrorf(paramNode.Pos(), "must have a schema")
		}

		if s == nil {
			s = interpreter.Schema().Copy()
		}
		o := s.(*schema.Object)
		if o.Parameters == nil {
			o.Parameters = map[string]schema.Schema{}
		}

		switch {
		case util.BoolValue(config.Input):
			if util.BoolValue(config.Required) {
				o.Required = append(o.Required, string(paramNode.Name()))
			}
			config.Schema.(schema.MetadataAccessor).SetAnnotation(conflow.AnnotationEvalStage, conflow.EvalStageInit.String())

		case util.BoolValue(config.Output):
			config.Schema.(schema.MetadataAccessor).SetAnnotation(conflow.AnnotationEvalStage, conflow.EvalStageClose.String())
			config.Schema.(schema.MetadataAccessor).SetReadOnly(true)
		}

		o.Parameters[string(paramNode.Name())] = config.Schema
		paramNode.SetSchema(config.Schema)
	}

	if s == nil {
		return interpreter.Schema(), nil
	}

	return s, nil
}

func getParameterConfig(param conflow.ParameterNode) (*conflow.ParameterConfig, parsley.Error) {
	config := &conflow.ParameterConfig{}

	for _, d := range param.Directives() {
		if d.EvalStage() != conflow.EvalStageParse {
			continue
		}

		evalCtx := conflow.NewEvalContext(context.Background(), nil, nil, job.SimpleScheduler{}, nil)

		block, err := d.Value(evalCtx)
		if err != nil {
			return nil, parsley.NewErrorf(d.Pos(), "failed to evaluate directive %s: %w", d.ID(), err)
		}

		opt, ok := block.(conflow.ParameterConfigOption)
		if !ok {
			return nil, parsley.NewError(d.Pos(), fmt.Errorf("%q directive can not be defined on a parameter", d.BlockType()))
		}

		opt.ApplyToParameterConfig(config)
	}

	return config, nil
}
