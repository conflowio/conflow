// Copyright (c) 2018 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package block

import (
	"errors"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/identifier"
	"github.com/opsidian/parsley/parsley"
)

func transformNode(
	parseCtx interface{},
	node parsley.Node,
	interpreter Interpreter,
) (parsley.Node, parsley.Error) {
	parseCtx = interpreter.ParseContext(parseCtx.(*basil.ParseContext))
	idRegistry := parseCtx.(basil.IDRegistryAware).IDRegistry()

	nodes := node.(parsley.NonTerminalNode).Children()
	blockIDNodes := nodes[0].(parsley.NonTerminalNode).Children()
	typeNode := blockIDNodes[0].(*identifier.Node)
	var idNode *identifier.Node
	if len(blockIDNodes) == 2 {
		idNode = blockIDNodes[1].(*identifier.Node)
		id, _ := idNode.Value(nil)
		if err := idRegistry.RegisterID(id.(basil.ID)); err != nil {
			return nil, parsley.NewError(idNode.Pos(), err)
		}
	} else {
		if interpreter.HasForeignID() {
			return nil, parsley.NewError(typeNode.ReaderPos(), errors.New("identifier must be set"))
		}

		id := idRegistry.GenerateID()
		idNode = identifier.NewNode(id, typeNode.ReaderPos(), typeNode.ReaderPos())
	}

	var children []parsley.Node

	if len(nodes) > 1 {
		blockValueNode := nodes[1]

		if blockValueNode.Token() == "BLOCK_BODY" {
			blockValueChildren := blockValueNode.(parsley.NonTerminalNode).Children()

			if len(blockValueChildren) > 2 {
				err := transformChildren(
					parseCtx,
					blockValueChildren[1].(parsley.NonTerminalNode).Children(),
				)
				if err != nil {
					return nil, err
				}
			}
		} else { // We have an expression as the value of the block
			valueParamName := interpreter.ValueParamName()
			if valueParamName == "" {
				blockType, _ := typeNode.Value(nil)
				return nil, parsley.NewErrorf(typeNode.Pos(), "%q block does not support short format", blockType)
			}
			valueNode, err := parsley.Transform(parseCtx, blockValueNode)
			if err != nil {
				return nil, err
			}
			children = []parsley.Node{
				NewParamNode(
					identifier.NewNode(valueParamName, blockValueNode.Pos(), blockValueNode.Pos()),
					valueNode,
					false,
				),
			}
		}
	}

	// TODO: calculate dependencies, eval stages and child order?

	res := &Node{
		idNode:      idNode,
		typeNode:    typeNode,
		children:    children,
		interpreter: interpreter,
		readerPos:   node.ReaderPos(),
	}

	if !interpreter.HasForeignID() {
		blockNodeRegistry := parseCtx.(basil.BlockNodeRegistryAware).BlockNodeRegistry()
		if err := blockNodeRegistry.AddBlockNode(res); err != nil {
			return nil, parsley.NewError(idNode.Pos(), err)
		}
	}

	return res, nil
}

func transformChildren(parseCtx interface{}, nodes []parsley.Node) parsley.Error {
	if len(nodes) == 0 {
		return nil
	}

	paramNames := make(map[basil.ID]struct{}, len(nodes))

	for i, node := range nodes {
		if node.Token() == "BLOCK" {
			blockNode, err := node.(parsley.Transformable).Transform(parseCtx)
			if err != nil {
				return err
			}
			nodes[i] = blockNode.(*Node)
		} else {
			paramChildren := node.(parsley.NonTerminalNode).Children()

			paramNode := paramChildren[0]
			paramName, _ := paramNode.Value(nil)
			if _, exists := paramNames[paramName.(basil.ID)]; exists {
				return parsley.NewErrorf(paramChildren[0].Pos(), "%q parameter was defined multiple times", paramName)
			}
			paramNames[paramName.(basil.ID)] = struct{}{}

			op, _ := paramChildren[1].Value(nil)
			isDeclaration := op == ":="

			valueNode, err := parsley.Transform(parseCtx, paramChildren[2])
			if err != nil {
				return err
			}

			nodes[i] = NewParamNode(paramNode, valueNode, isDeclaration)
		}
	}

	return nil
}
