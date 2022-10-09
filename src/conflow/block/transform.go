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

	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/conflow/dependency"
	"github.com/conflowio/conflow/src/conflow/directive"
	"github.com/conflowio/conflow/src/conflow/job"
	"github.com/conflowio/conflow/src/conflow/parameter"
	"github.com/conflowio/conflow/src/schema"
)

func GetIDAndNameFromNode(node parsley.Node) (*conflow.IDNode, *conflow.NameNode) {
	var idNode *conflow.IDNode
	var nameNode *conflow.NameNode
	switch n := node.(type) {
	case parsley.NonTerminalNode:
		if _, isEmpty := n.Children()[0].(ast.EmptyNode); !isEmpty {
			idNode = n.Children()[0].(*conflow.IDNode)
		}
		nameNode = n.Children()[1].(*conflow.NameNode)
	case *conflow.NameNode:
		nameNode = n
	case *conflow.IDNode:
		nameNode = conflow.NewNameNode(nil, nil, n)
	default:
		panic(fmt.Errorf("unexpected identifier node: %T", node))
	}

	return idNode, nameNode
}

func TransformNode(ctx interface{}, node parsley.Node, interpreter conflow.BlockInterpreter) (parsley.Node, parsley.Error) {
	parseCtx := interpreter.ParseContext(ctx.(*conflow.ParseContext))
	nodes := node.(parsley.NonTerminalNode).Children()

	idNode, nameNode := GetIDAndNameFromNode(nodes[1])
	isModule := nameNode.Value() == "module"

	if idNode != nil {
		if err := parseCtx.RegisterID(idNode.ID()); err != nil {
			return nil, parsley.NewError(idNode.Pos(), err)
		}
	} else {
		if isModule {
			return nil, parsley.NewErrorf(nameNode.Pos(), "a module must have an identifier")
		}

		id := parseCtx.GenerateID()
		idNode = conflow.NewIDNode(id, conflow.ClassifierNone, nameNode.Pos(), nameNode.Pos())
	}

	if isModule {
		parseCtx = parseCtx.NewForModule()
	}

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
			paramNode.SetSchema(interpreter.Schema().(*schema.Object).Parameters[string(valueParamName)])

			var deps conflow.Dependencies
			children, deps, err = dependency.NewResolver(idNode.ID(), paramNode).Resolve()
			if err != nil {
				return nil, err
			}
			dependencies.Add(deps)
		}
	}

	if isModule && len(dependencies) > 0 {
		for _, d := range dependencies {
			if _, blockNodeExists := parseCtx.BlockNode(d.ParentID()); blockNodeExists {
				return nil, parsley.NewErrorf(d.Pos(), "unknown parameter: %q", d.ID())
			} else {
				return nil, parsley.NewErrorf(d.Pos(), "unknown block: %q", d.ParentID())
			}
		}
	}

	return NewNode(
		idNode,
		nameNode,
		children,
		node.Token(),
		directives,
		node.ReaderPos(),
		interpreter,
		dependencies,
	), nil
}

func TransformRootNode(ctx interface{}, node parsley.Node, id conflow.ID, interpreter conflow.BlockInterpreter) (parsley.Node, parsley.Error) {
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

	if len(dependencies) > 0 {
		for _, d := range dependencies {
			if _, blockNodeExists := parseCtx.BlockNode(d.ParentID()); blockNodeExists {
				return nil, parsley.NewErrorf(d.Pos(), "unknown parameter: %q", d.ID())
			} else {
				return nil, parsley.NewErrorf(d.Pos(), "unknown block: %q", d.ParentID())
			}
		}
	}

	var hasModules bool
	for _, child := range children {
		if b, ok := child.(conflow.BlockNode); ok {
			if b.BlockType() == "module" {
				hasModules = true
				break
			}
		}
	}

	blockName := "root"
	if !hasModules {
		blockName = "main"
	}

	blockNode := NewNode(
		conflow.NewIDNode(id, conflow.ClassifierNone, node.Pos(), node.Pos()),
		conflow.NewNameNode(nil, nil, conflow.NewIDNode(blockName, conflow.ClassifierNone, node.Pos(), node.Pos())),
		children,
		TokenBlock,
		nil,
		node.ReaderPos(),
		interpreter,
		nil,
	)

	if !hasModules {
		moduleInterpreter, err := NewModuleInterpreter(blockNode)
		if err != nil {
			return nil, nil, err
		}
		registry := parseCtx.BlockTransformerRegistry().(InterpreterRegistry)
		registry[string(blockNode.ID())] = moduleInterpreter
	}

	if err := parseCtx.AddBlockNode(blockNode); err != nil {
		panic(fmt.Errorf("failed to register the root block node: %w", err))
	}

	return blockNode, nil
}

func TransformChildren(
	parseCtx *conflow.ParseContext,
	blockID conflow.ID,
	children []parsley.Node,
	interpreter conflow.BlockInterpreter,
) ([]conflow.Node, conflow.Dependencies, parsley.Error) {
	if len(children) == 0 {
		return nil, nil, nil
	}

	nodes := make([]conflow.Node, 0, len(children))
	paramNames := make(map[conflow.ID]struct{}, len(children))

	for _, node := range children {
		if node.Token() == TokenBlock {
			res, err := node.(parsley.Transformable).Transform(parseCtx)
			if err != nil {
				return nil, nil, err
			}
			blockNode := res.(conflow.BlockNode)

			if blockNode.BlockType() == "module" {
				moduleInterpreter, err := NewModuleInterpreter(blockNode)
				if err != nil {
					return nil, nil, err
				}
				registry := parseCtx.BlockTransformerRegistry().(InterpreterRegistry)
				registry[string(blockNode.ID())] = moduleInterpreter

				continue
			}

			if err := parseCtx.AddBlockNode(blockNode); err != nil {
				return nil, nil, parsley.NewError(blockNode.Pos(), err)
			}

			// TODO: what's the use case for these?
			if blockSchema, ok := interpreter.Schema().(*schema.Object).Parameters[string(blockNode.ID())]; ok {
				blockNode.SetSchema(blockSchema)
			} else if blockSchema, ok := interpreter.Schema().(*schema.Object).Parameters[string(blockNode.ParameterName())]; ok {
				blockNode.SetSchema(blockSchema)
			}

			nodes = append(nodes, res.(conflow.BlockNode))

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

			if paramSchema, ok := interpreter.Schema().(*schema.Object).Parameters[string(paramNode.Name())]; ok {
				paramNode.SetSchema(paramSchema)
			}

			nodes = append(nodes, paramNode)
		} else {
			panic(fmt.Errorf("invalid block child node: %T", node))
		}
	}

	return dependency.NewResolver(blockID, nodes...).Resolve()
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
