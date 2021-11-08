// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package block

import (
	"context"
	"fmt"

	"github.com/conflowio/parsley/ast"
	"github.com/conflowio/parsley/parsley"

	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/dependency"
	"github.com/conflowio/conflow/conflow/directive"
	"github.com/conflowio/conflow/conflow/job"
	"github.com/conflowio/conflow/conflow/parameter"
	"github.com/conflowio/conflow/conflow/schema"
	"github.com/conflowio/conflow/util"
)

func TransformNode(ctx interface{}, node parsley.Node, interpreter conflow.BlockInterpreter) (parsley.Node, parsley.Error) {
	parseCtx := interpreter.ParseContext(ctx.(*conflow.ParseContext))
	nodes := node.(parsley.NonTerminalNode).Children()

	var directives []conflow.BlockNode
	dependencies := make(conflow.Dependencies)

	if n, ok := nodes[0].(parsley.NonTerminalNode); ok && len(n.Children()) > 0 {
		var err parsley.Error
		var deps conflow.Dependencies
		if directives, deps, err = directive.Transform(parseCtx, n.Children()); err != nil {
			return nil, err
		}
		dependencies.Add(deps)
	}

	var idNode *conflow.IDNode
	var nameNode *conflow.NameNode
	switch n := nodes[1].(type) {
	case parsley.NonTerminalNode:
		if _, isEmpty := n.Children()[0].(ast.EmptyNode); !isEmpty {
			idNode = n.Children()[0].(*conflow.IDNode)
			if err := parseCtx.RegisterID(idNode.ID()); err != nil {
				return nil, parsley.NewError(idNode.Pos(), err)
			}
		}

		nameNode = n.Children()[1].(*conflow.NameNode)
	case *conflow.NameNode:
		nameNode = n
	case *conflow.IDNode:
		nameNode = conflow.NewNameNode(nil, nil, n)
	default:
		panic(fmt.Errorf("unexpected identifier node: %T", nodes[1]))
	}

	if idNode == nil {
		id := parseCtx.GenerateID()
		idNode = conflow.NewIDNode(id, conflow.ClassifierNone, nameNode.Pos(), nameNode.Pos())
	}

	var children []conflow.Node
	if len(nodes) > 2 {
		blockValueNode := nodes[2]
		if blockValueNode.Token() == TokenBlockBody {
			blockValueChildren := blockValueNode.(parsley.NonTerminalNode).Children()

			if len(blockValueChildren) > 2 {
				var err parsley.Error
				var deps conflow.Dependencies
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
				conflow.NewIDNode(valueParamName, conflow.ClassifierNone, blockValueNode.Pos(), blockValueNode.Pos()),
				valueNode,
				false,
				nil,
			)
			paramNode.SetSchema(interpreter.Schema().(*schema.Object).Properties[string(valueParamName)])

			var deps conflow.Dependencies
			children, deps, err = dependency.NewResolver(idNode.ID(), paramNode).Resolve()
			if err != nil {
				return nil, err
			}
			dependencies.Add(deps)
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

func TransformMainNode(ctx interface{}, node parsley.Node, id conflow.ID, interpreter conflow.BlockInterpreter) (parsley.Node, parsley.Error) {
	parseCtx := interpreter.ParseContext(ctx.(*conflow.ParseContext))

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
		conflow.NewIDNode(id, conflow.ClassifierNone, node.Pos(), node.Pos()),
		conflow.NewNameNode(nil, nil, conflow.NewIDNode(conflow.ID("main"), conflow.ClassifierNone, node.Pos(), node.Pos())),
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
	parseCtx *conflow.ParseContext,
	blockID conflow.ID,
	nodes []parsley.Node,
	interpreter conflow.BlockInterpreter,
) ([]conflow.Node, conflow.Dependencies, parsley.Error) {
	if len(nodes) == 0 {
		return nil, nil, nil
	}

	conflowNodes := make([]conflow.Node, 0, len(nodes))
	paramNames := make(map[conflow.ID]struct{}, len(nodes))

	for _, node := range nodes {
		if node.Token() == TokenBlock {
			res, err := node.(parsley.Transformable).Transform(parseCtx)
			if err != nil {
				return nil, nil, err
			}
			blockNode := res.(conflow.BlockNode)

			if blockSchema, ok := interpreter.Schema().(*schema.Object).Properties[string(blockNode.ID())]; ok {
				blockNode.SetSchema(blockSchema)
			} else if blockSchema, ok := interpreter.Schema().(*schema.Object).Properties[string(blockNode.ParameterName())]; ok {
				blockNode.SetSchema(blockSchema)
			}

			conflowNodes = append(conflowNodes, res.(conflow.BlockNode))

			if blockNode.EvalStage() == conflow.EvalStageParse {
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

			conflowNodes = append(conflowNodes, paramNode)
		} else {
			panic(fmt.Errorf("invalid block child node: %T", node))
		}
	}

	return dependency.NewResolver(blockID, conflowNodes...).Resolve()
}

func getModuleSchema(children []conflow.Node, interpreter conflow.BlockInterpreter) (schema.Schema, parsley.Error) {
	var s schema.Schema

	for _, c := range children {
		if paramNode, ok := c.(conflow.ParameterNode); ok {
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

				config.Schema.(schema.MetadataAccessor).SetAnnotation(conflow.AnnotationEvalStage, conflow.EvalStageInit.String())
				config.Schema.(schema.MetadataAccessor).SetAnnotation(conflow.AnnotationUserDefined, "true")

				o.Properties[string(paramNode.Name())] = config.Schema
				paramNode.SetSchema(config.Schema)
			case util.BoolValue(config.Output):
				if config.Schema == nil {
					return nil, parsley.NewErrorf(paramNode.Pos(), "must have a schema")
				}

				config.Schema.(schema.MetadataAccessor).SetAnnotation(conflow.AnnotationEvalStage, conflow.EvalStageInit.String())
				config.Schema.(schema.MetadataAccessor).SetAnnotation(conflow.AnnotationUserDefined, "true")
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

func evaluateBlock(parseCtx *conflow.ParseContext, node conflow.BlockNode) parsley.Error {
	evalCtx := conflow.NewEvalContext(context.Background(), nil, nil, job.SimpleScheduler{}, nil)

	block, err := parsley.EvaluateNode(evalCtx, node)
	if err != nil {
		return err
	}

	if provider, ok := block.(conflow.BlockProvider); ok {
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
