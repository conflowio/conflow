// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package block

import (
	"context"
	"fmt"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/parsley"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/dependency"
	"github.com/opsidian/basil/basil/directive"
	"github.com/opsidian/basil/basil/job"
	"github.com/opsidian/basil/basil/parameter"
	"github.com/opsidian/basil/basil/schema"
	"github.com/opsidian/basil/util"
)

func TransformNode(ctx interface{}, node parsley.Node, interpreter basil.BlockInterpreter) (parsley.Node, parsley.Error) {
	parseCtx := interpreter.ParseContext(ctx.(*basil.ParseContext))
	nodes := node.(parsley.NonTerminalNode).Children()

	var directives []basil.BlockNode
	dependencies := make(basil.Dependencies)

	if n, ok := nodes[0].(parsley.NonTerminalNode); ok && len(n.Children()) > 0 {
		var err parsley.Error
		var deps basil.Dependencies
		if directives, deps, err = directive.Transform(parseCtx, n.Children()); err != nil {
			return nil, err
		}
		dependencies.Add(deps)
	}

	var idNode *basil.IDNode
	var nameNode *basil.NameNode
	switch n := nodes[1].(type) {
	case parsley.NonTerminalNode:
		if _, isEmpty := n.Children()[0].(ast.EmptyNode); !isEmpty {
			idNode = n.Children()[0].(*basil.IDNode)
			if err := parseCtx.RegisterID(idNode.ID()); err != nil {
				return nil, parsley.NewError(idNode.Pos(), err)
			}
		}

		nameNode = n.Children()[1].(*basil.NameNode)
	case *basil.NameNode:
		nameNode = n
	case *basil.IDNode:
		nameNode = basil.NewNameNode(nil, nil, n)
	default:
		panic(fmt.Errorf("unexpected identifier node: %T", nodes[1]))
	}

	if idNode == nil {
		id := parseCtx.GenerateID()
		idNode = basil.NewIDNode(id, basil.ClassifierNone, nameNode.Pos(), nameNode.Pos())
	}

	var children []basil.Node
	if len(nodes) > 2 {
		blockValueNode := nodes[2]
		if blockValueNode.Token() == TokenBlockBody {
			blockValueChildren := blockValueNode.(parsley.NonTerminalNode).Children()

			if len(blockValueChildren) > 2 {
				var err parsley.Error
				var deps basil.Dependencies
				children, deps, err = TransformChildren(
					parseCtx,
					idNode.ID(),
					blockValueChildren[1].(parsley.NonTerminalNode).Children(),
					interpreter,
				)
				if err != nil {
					return nil, err
				}
				dependencies.Add(deps)
			}
		} else if _, empty := blockValueNode.(ast.EmptyNode); !empty { // We have an expression as the value of the block
			valueParamName := interpreter.ValueParamName()
			if valueParamName == "" {
				return nil, parsley.NewErrorf(nameNode.Pos(), "%q block does not support short format", nameNode.NameNode().ID())
			}
			valueNode, err := parsley.Transform(parseCtx, blockValueNode)
			if err != nil {
				return nil, err
			}

			paramNode := parameter.NewNode(
				idNode.ID(),
				basil.NewIDNode(valueParamName, basil.ClassifierNone, blockValueNode.Pos(), blockValueNode.Pos()),
				valueNode,
				false,
				nil,
			)
			paramNode.SetSchema(interpreter.Schema().(*schema.Object).Properties[string(valueParamName)])

			children = []basil.Node{paramNode}
		}
	}

	res := NewNode(
		idNode,
		nameNode,
		children,
		node.Token(),
		directives,
		node.ReaderPos(),
		interpreter,
		dependencies,
	)

	if err := parseCtx.AddBlockNode(res); err != nil {
		return nil, parsley.NewError(idNode.Pos(), err)
	}

	return res, nil
}

func TransformMainNode(ctx interface{}, node parsley.Node, id basil.ID, interpreter basil.BlockInterpreter) (parsley.Node, parsley.Error) {
	parseCtx := interpreter.ParseContext(ctx.(*basil.ParseContext))

	children, dependencies, err := TransformChildren(
		parseCtx,
		id,
		node.(parsley.NonTerminalNode).Children(),
		interpreter,
	)
	if err != nil {
		return nil, err
	}

	moduleSchema, err := getModuleSchema(children, interpreter)
	if err != nil {
		return nil, err
	}

	if moduleSchema != interpreter.Schema() {
		interpreter = &mainInterpreter{BlockInterpreter: interpreter, schema: moduleSchema}
	}

	if len(dependencies) > 0 {
		for _, d := range dependencies {
			if _, blockNodeExists := parseCtx.BlockNode(d.ParentID()); blockNodeExists {
				return nil, parsley.NewErrorf(d.Pos(), "unknown parameter: %q", d.ID())
			} else {
				return nil, parsley.NewErrorf(d.Pos(), "unknown block: %q", d.ParentID())
			}
		}
	}

	res := NewNode(
		basil.NewIDNode(id, basil.ClassifierNone, node.Pos(), node.Pos()),
		basil.NewNameNode(nil, nil, basil.NewIDNode(basil.ID("main"), basil.ClassifierNone, node.Pos(), node.Pos())),
		children,
		TokenBlock,
		nil,
		node.ReaderPos(),
		interpreter,
		nil,
	)

	if err := parseCtx.AddBlockNode(res); err != nil {
		panic(fmt.Errorf("failed to register the main block node: %w", err))
	}

	return res, nil
}

func TransformChildren(
	parseCtx *basil.ParseContext,
	blockID basil.ID,
	nodes []parsley.Node,
	interpreter basil.BlockInterpreter,
) ([]basil.Node, basil.Dependencies, parsley.Error) {
	if len(nodes) == 0 {
		return nil, nil, nil
	}

	basilNodes := make([]basil.Node, 0, len(nodes))
	paramNames := make(map[basil.ID]struct{}, len(nodes))

	for _, node := range nodes {
		if node.Token() == TokenBlock {
			res, err := node.(parsley.Transformable).Transform(parseCtx)
			if err != nil {
				return nil, nil, err
			}
			blockNode := res.(basil.BlockNode)

			if blockSchema, ok := interpreter.Schema().(*schema.Object).Properties[string(blockNode.ID())]; ok {
				blockNode.SetSchema(blockSchema)
			} else if blockSchema, ok := interpreter.Schema().(*schema.Object).Properties[string(blockNode.ParameterName())]; ok {
				blockNode.SetSchema(blockSchema)
			}

			basilNodes = append(basilNodes, res.(basil.BlockNode))

			if blockNode.EvalStage() == basil.EvalStageParse {
				if err := evaluateBlock(parseCtx, blockNode); err != nil {
					return nil, nil, err
				}
			}
		} else if node.Token() == parameter.Token {
			paramNode, err := parameter.TransformNode(parseCtx, node, blockID, paramNames)
			if err != nil {
				return nil, nil, err
			}

			if paramSchema, ok := interpreter.Schema().(*schema.Object).Properties[string(paramNode.Name())]; ok {
				paramNode.SetSchema(paramSchema)
			}

			basilNodes = append(basilNodes, paramNode)
		} else {
			panic(fmt.Errorf("invalid block child node: %T", node))
		}
	}

	return dependency.NewResolver(blockID, basilNodes...).Resolve()
}

func getModuleSchema(children []basil.Node, interpreter basil.BlockInterpreter) (schema.Schema, parsley.Error) {
	var s schema.Schema

	for _, c := range children {
		if paramNode, ok := c.(basil.ParameterNode); ok {
			config, err := getParameterConfig(paramNode)
			if err != nil {
				return nil, err
			}

			if util.BoolValue(config.Input) || util.BoolValue(config.Output) {
				if _, exists := interpreter.Schema().(schema.ObjectKind).GetProperties()[string(paramNode.Name())]; exists {
					return nil, parsley.NewErrorf(c.Pos(), "%q parameter already exists.", paramNode.Name())
				}
			} else {
				continue
			}

			if s == nil {
				s = interpreter.Schema().Copy()
			}
			o := s.(*schema.Object)
			if o.Properties == nil {
				o.Properties = map[string]schema.Schema{}
			}

			switch {
			case util.BoolValue(config.Input):
				if util.BoolValue(config.Required) {
					o.Required = append(o.Required, string(paramNode.Name()))
				}

				if config.Schema == nil {
					return nil, parsley.NewErrorf(paramNode.Pos(), "must have a schema")
				}

				config.Schema.(schema.MetadataAccessor).SetAnnotation("eval_stage", basil.EvalStageInit.String())
				config.Schema.(schema.MetadataAccessor).SetAnnotation("user_defined", "true")

				o.Properties[string(paramNode.Name())] = config.Schema
				paramNode.SetSchema(config.Schema)
			case util.BoolValue(config.Output):
				if config.Schema == nil {
					return nil, parsley.NewErrorf(paramNode.Pos(), "must have a schema")
				}

				config.Schema.(schema.MetadataAccessor).SetAnnotation("eval_stage", basil.EvalStageInit.String())
				config.Schema.(schema.MetadataAccessor).SetAnnotation("user_defined", "true")
				config.Schema.(schema.MetadataAccessor).SetReadOnly(true)

				o.Properties[string(paramNode.Name())] = config.Schema
				paramNode.SetSchema(config.Schema)
			}
		}
	}

	if s == nil {
		return interpreter.Schema(), nil
	}

	return s, nil
}

func getParameterConfig(param basil.ParameterNode) (*basil.ParameterConfig, parsley.Error) {
	config := &basil.ParameterConfig{}

	for _, d := range param.Directives() {
		if d.EvalStage() != basil.EvalStageParse {
			continue
		}

		evalCtx := basil.NewEvalContext(context.Background(), nil, nil, job.SimpleScheduler{}, nil)

		block, err := d.Value(evalCtx)
		if err != nil {
			return nil, parsley.NewErrorf(d.Pos(), "failed to evaluate directive %s: %w", d.ID(), err)
		}

		opt, ok := block.(basil.ParameterConfigOption)
		if !ok {
			return nil, parsley.NewError(d.Pos(), fmt.Errorf("%q directive can not be defined on a parameter", d.BlockType()))
		}

		opt.ApplyToParameterConfig(config)
	}

	return config, nil
}

func evaluateBlock(parseCtx *basil.ParseContext, node basil.BlockNode) parsley.Error {
	evalCtx := basil.NewEvalContext(context.Background(), nil, nil, job.SimpleScheduler{}, nil)

	block, err := parsley.EvaluateNode(evalCtx, node)
	if err != nil {
		return err
	}

	if provider, ok := block.(basil.BlockProvider); ok {
		interpreters, err := provider.BlockInterpreters(parseCtx)
		if err != nil {
			return parsley.NewError(node.Pos(), err)
		}

		registry := parseCtx.BlockTransformerRegistry().(InterpreterRegistry)
		for name, interpreter := range interpreters {
			if _, exists := registry[string(name)]; exists {
				return parsley.NewErrorf(node.Pos(), "%q block is already registered, please use an alias", name)
			}
			registry[string(name)] = interpreter
		}
	}

	return nil
}
