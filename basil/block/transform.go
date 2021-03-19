// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package block

import (
	"context"
	"errors"
	"fmt"

	"github.com/opsidian/basil/basil/job"

	"github.com/opsidian/basil/util"

	"github.com/opsidian/basil/basil/directive"

	"github.com/opsidian/parsley/ast"

	"github.com/opsidian/basil/basil/dependency"
	"github.com/opsidian/basil/basil/parameter"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/parsley/parsley"
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
	var typeNode *basil.IDNode
	switch n := nodes[1].(type) {
	case parsley.NonTerminalNode:
		idNode = n.Children()[0].(*basil.IDNode)
		typeNode = n.Children()[1].(*basil.IDNode)
		if err := parseCtx.RegisterID(idNode.ID()); err != nil {
			return nil, parsley.NewError(idNode.Pos(), err)
		}
	case *basil.IDNode:
		typeNode = n
		if interpreter.HasForeignID() {
			return nil, parsley.NewError(typeNode.ReaderPos(), errors.New("identifier must be set"))
		}

		id := parseCtx.GenerateID()
		idNode = basil.NewIDNode(id, typeNode.ReaderPos(), typeNode.ReaderPos())
	default:
		panic(fmt.Errorf("unexpected identifier node: %T", nodes[1]))
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
				return nil, parsley.NewErrorf(typeNode.Pos(), "%q block does not support short format", typeNode.ID())
			}
			valueNode, err := parsley.Transform(parseCtx, blockValueNode)
			if err != nil {
				return nil, err
			}
			children = []basil.Node{
				parameter.NewNode(
					idNode.ID(),
					basil.NewIDNode(valueParamName, blockValueNode.Pos(), blockValueNode.Pos()),
					valueNode,
					false,
					nil,
				),
			}
		}
	}

	res := NewNode(
		idNode,
		typeNode,
		children,
		node.Token(),
		directives,
		node.ReaderPos(),
		interpreter,
		dependencies,
	)

	if !interpreter.HasForeignID() {
		if err := parseCtx.AddBlockNode(res); err != nil {
			return nil, parsley.NewError(idNode.Pos(), err)
		}
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

	moduleParams, err := getModuleParams(children, interpreter)
	if err != nil {
		return nil, err
	}

	if len(moduleParams) > 0 {
		interpreter = newMainInterpreter(interpreter, moduleParams)
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
		basil.NewIDNode(id, node.Pos(), node.Pos()),
		basil.NewIDNode(basil.ID("main"), node.Pos(), node.Pos()),
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

	blocks := interpreter.Blocks()
	parameters := interpreter.Params()

	for _, node := range nodes {
		if node.Token() == TokenBlock {
			res, err := node.(parsley.Transformable).Transform(parseCtx)
			if err != nil {
				return nil, nil, err
			}
			blockNode := res.(basil.BlockNode)
			blockNode.SetDescriptor(blocks[blockNode.BlockType()])
			basilNodes = append(basilNodes, blockNode)

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
			if descriptor, ok := parameters[paramNode.Name()]; ok {
				paramNode.SetDescriptor(descriptor)
			}
			basilNodes = append(basilNodes, paramNode)
		} else {
			panic(fmt.Errorf("invalid block child node: %T", node))
		}
	}

	return dependency.NewResolver(blockID, basilNodes...).Resolve()
}

func getModuleParams(children []basil.Node, interpreter basil.BlockInterpreter) (map[basil.ID]basil.ParameterDescriptor, parsley.Error) {
	moduleParams := map[basil.ID]basil.ParameterDescriptor{}

	for _, c := range children {
		if paramNode, ok := c.(basil.ParameterNode); ok {
			config, err := getParameterConfig(paramNode)
			if err != nil {
				return nil, err
			}

			if util.BoolValue(config.Input) || util.BoolValue(config.Output) {
				if _, exists := interpreter.Params()[paramNode.Name()]; exists {
					return nil, parsley.NewErrorf(c.Pos(), "%q parameter already exists.", paramNode.Name())
				}
			}

			switch {
			case util.BoolValue(config.Input):
				moduleParams[paramNode.Name()] = basil.ParameterDescriptor{
					Type:          util.StringValue(config.Type),
					EvalStage:     basil.EvalStageInit,
					IsRequired:    util.BoolValue(config.Required),
					IsUserDefined: true,
				}
			case util.BoolValue(config.Output):
				moduleParams[paramNode.Name()] = basil.ParameterDescriptor{
					Type:          util.StringValue(config.Type),
					EvalStage:     basil.EvalStageInit,
					IsUserDefined: true,
					IsOutput:      true,
				}
			}
		}
	}

	return moduleParams, nil
}

func getParameterConfig(param basil.ParameterNode) (*basil.ParameterConfig, parsley.Error) {
	config := &basil.ParameterConfig{}

	evalCtx := basil.NewEvalContext(context.Background(), nil, nil, job.SimpleScheduler{}, nil)

	for _, d := range param.Directives() {
		if d.EvalStage() != basil.EvalStageParse {
			continue
		}

		block, err := d.Value(evalCtx)
		if err != nil {
			return nil, parsley.NewErrorf(d.Pos(), "failed to evaluate directive %s: %w", d.ID(), err)
		}

		opt, ok := block.(basil.ParameterConfigOption)
		if !ok {
			return nil, parsley.NewError(d.Pos(), fmt.Errorf("%q directive can not be defined on a parameter", d.Type()))
		}

		opt.ApplyToParameterConfig(config)
	}

	return config, nil
}

func evaluateBlock(parseCtx *basil.ParseContext, node basil.BlockNode) parsley.Error {
	evalCtx := basil.NewEvalContext(context.Background(), nil, nil, job.SimpleScheduler{}, nil)

	block, err := node.Value(evalCtx)
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
