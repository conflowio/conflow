// Copyright (c) 2018 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package block

import (
	"errors"

	"github.com/opsidian/basil/basil/dependency"
	"github.com/opsidian/basil/basil/parameter"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/parsley/parsley"
)

func TransformNode(ctx interface{}, node parsley.Node, interpreter basil.BlockInterpreter) (parsley.Node, parsley.Error) {
	parseCtx := interpreter.ParseContext(ctx.(*basil.ParseContext))

	nodes := node.(parsley.NonTerminalNode).Children()
	blockIDNodes := nodes[0].(parsley.NonTerminalNode).Children()
	typeNode := blockIDNodes[0].(*basil.IDNode)
	var idNode *basil.IDNode
	if len(blockIDNodes) == 2 {
		idNode = blockIDNodes[1].(*basil.IDNode)
		if err := parseCtx.RegisterID(idNode.ID()); err != nil {
			return nil, parsley.NewError(idNode.Pos(), err)
		}
	} else {
		if interpreter.HasForeignID() {
			return nil, parsley.NewError(typeNode.ReaderPos(), errors.New("identifier must be set"))
		}

		// TODO: there is a chance the id generator will generate an existing, manually defined id
		id := parseCtx.GenerateID()
		idNode = basil.NewIDNode(id, typeNode.ReaderPos(), typeNode.ReaderPos())
	}

	var children []basil.Node

	var dependencies []basil.VariableNode
	if len(nodes) > 1 {
		blockValueNode := nodes[1]

		if blockValueNode.Token() == TokenBlockBody {
			blockValueChildren := blockValueNode.(parsley.NonTerminalNode).Children()

			if len(blockValueChildren) > 2 {
				var err parsley.Error
				children, dependencies, err = TransformChildren(
					parseCtx,
					idNode.ID(),
					blockValueChildren[1].(parsley.NonTerminalNode).Children(),
				)
				if err != nil {
					return nil, err
				}
			}
		} else { // We have an expression as the value of the block
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
				),
			}
		}
	}

	res := &Node{
		idNode:       idNode,
		typeNode:     typeNode,
		children:     children,
		interpreter:  interpreter,
		readerPos:    node.ReaderPos(),
		dependencies: dependencies,
	}

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
	)
	if err != nil {
		return nil, err
	}

	if len(dependencies) > 0 {
		d := dependencies[0]
		if _, blockNodeExists := parseCtx.BlockNode(d.ParentID()); blockNodeExists {
			return nil, parsley.NewErrorf(dependencies[0].Pos(), "unknown parameter: %q", d.ID())
		} else {
			return nil, parsley.NewErrorf(dependencies[0].Pos(), "unknown block: %q", d.ParentID())
		}
	}

	res := &Node{
		idNode:      basil.NewIDNode(id, node.Pos(), node.Pos()),
		typeNode:    basil.NewIDNode(basil.ID("main"), node.Pos(), node.Pos()),
		children:    children,
		interpreter: interpreter,
		readerPos:   node.ReaderPos(),
	}

	if err := parseCtx.AddBlockNode(res); err != nil {
		panic("failed to register the main block node")
	}

	return res, nil
}

func TransformChildren(
	parseCtx interface{},
	blockID basil.ID,
	nodes []parsley.Node,
) ([]basil.Node, []basil.VariableNode, parsley.Error) {
	if len(nodes) == 0 {
		return nil, nil, nil
	}

	res := make([]basil.Node, 0, len(nodes))
	paramNames := make(map[basil.ID]struct{}, len(nodes))

	for _, node := range nodes {
		if node.Token() == TokenBlock {
			blockNode, err := node.(parsley.Transformable).Transform(parseCtx)
			if err != nil {
				return nil, nil, err
			}
			res = append(res, blockNode.(basil.Node))
		} else if node.Token() == TokenParameter {
			paramNode, err := TransformParamNode(parseCtx, node, blockID, paramNames)
			if err != nil {
				return nil, nil, err
			}
			res = append(res, paramNode)
		}
	}

	return dependency.NewResolver(blockID, res...).Resolve()
}

func TransformParamNode(
	parseCtx interface{},
	node parsley.Node,
	blockID basil.ID,
	paramNames map[basil.ID]struct{},
) (basil.ParameterNode, parsley.Error) {
	paramChildren := node.(parsley.NonTerminalNode).Children()

	nameNode := paramChildren[0].(*basil.IDNode)
	if _, exists := paramNames[nameNode.ID()]; exists {
		return nil, parsley.NewErrorf(
			paramChildren[0].Pos(),
			"%q parameter was defined multiple times", nameNode.ID(),
		)
	}
	paramNames[nameNode.ID()] = struct{}{}

	op, _ := paramChildren[1].Value(nil)
	isDeclaration := op == ":="

	valueNode, err := parsley.Transform(parseCtx, paramChildren[2])
	if err != nil {
		return nil, err
	}

	return parameter.NewNode(blockID, nameNode, valueNode, isDeclaration), nil
}