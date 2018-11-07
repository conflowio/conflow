// Copyright (c) 2018 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package block

import (
	"fmt"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/identifier"
	"github.com/opsidian/parsley/parsley"
)

// NodeTransformerRegistry contains named node transformers
type NodeTransformerRegistry interface {
	NodeTransformer(name string) (parsley.NodeTransformer, bool)
}

// TransformNode returns with a node transformer function for a block
func TransformNode(registry NodeTransformerRegistry) parsley.NodeTransformFunc {
	var f parsley.NodeTransformFunc
	f = parsley.NodeTransformFunc(func(node parsley.Node) (parsley.Node, parsley.Error) {
		var err parsley.Error
		switch n := node.(type) {
		case parsley.NonTerminalNode:
			if node.Token() == "BLOCK" {
				nodes := node.(parsley.NonTerminalNode).Children()
				blockIDNodes := nodes[0].(parsley.NonTerminalNode).Children()
				typeNode := blockIDNodes[0]
				blockType, _ := typeNode.Value(nil)

				transformer, exists := registry.NodeTransformer(string(blockType.(basil.ID)))
				if !exists {
					return nil, parsley.NewError(typeNode.Pos(), fmt.Errorf("%q type is invalid or not allowed here", blockType))
				}

				return transformer.TransformNode(node)
			}

			children := n.Children()
			for i, childNode := range children {
				if children[i], err = f(childNode); err != nil {
					return nil, err
				}
			}
			return n, nil
		}

		return node, nil
	})

	return f
}

func transformNode(
	node parsley.Node,
	interpreter Interpreter,
) (parsley.Node, parsley.Error) {
	nodes := node.(parsley.NonTerminalNode).Children()
	blockIDNodes := nodes[0].(parsley.NonTerminalNode).Children()
	typeNode := blockIDNodes[0].(*identifier.Node)
	var idNode *identifier.Node
	if len(blockIDNodes) == 2 {
		idNode = blockIDNodes[1].(*identifier.Node)
	}

	var paramNodes map[basil.ID]basil.BlockParamNode
	var blockNodes []basil.BlockNode

	if len(nodes) > 1 {
		blockValueNode := nodes[1]

		if blockValueNode.Token() == "BLOCK_BODY" {
			blockValueChildren := blockValueNode.(parsley.NonTerminalNode).Children()

			if len(blockValueChildren) > 2 {
				blockChildren := blockValueChildren[1].(parsley.NonTerminalNode).Children()

				paramCnt := 0
				blockCnt := 0
				for _, blockChild := range blockChildren {
					if blockChild.Token() == "BLOCK" {
						blockCnt++
					} else {
						paramCnt++
					}
				}

				if paramCnt > 0 {
					paramNodes = make(map[basil.ID]basil.BlockParamNode, paramCnt)
				}
				if blockCnt > 0 {
					blockNodes = make([]basil.BlockNode, 0, blockCnt)
				}

				for _, blockChild := range blockChildren {
					if blockChild.Token() == "BLOCK" {
						childBlock, err := TransformNode(interpreter)(blockChild)
						if err != nil {
							return nil, err
						}
						blockNodes = append(blockNodes, childBlock.(*Node))
					} else {
						children := blockChild.(parsley.NonTerminalNode).Children()
						paramName, _ := children[0].Value(nil)
						if _, alreadyExists := paramNodes[paramName.(basil.ID)]; alreadyExists {
							return nil, parsley.NewErrorf(children[0].Pos(), "%q parameter was defined multiple times", paramName)
						}
						paramNodes[paramName.(basil.ID)] = NewParamNode(children[0], children[2])
					}
				}
			}
		} else { // We have an expression as the value of the block
			valueParamName := interpreter.ValueParamName()
			if valueParamName == "" {
				blockType, _ := typeNode.Value(nil)
				return nil, parsley.NewErrorf(typeNode.Pos(), "%q block does not support short format", blockType)
			}
			paramNodes = map[basil.ID]basil.BlockParamNode{
				valueParamName: NewParamNode(
					identifier.NewNode(valueParamName, blockValueNode.Pos(), blockValueNode.Pos()),
					blockValueNode,
				),
			}
		}
	}

	return &Node{
		idNode:      idNode,
		typeNode:    typeNode,
		paramNodes:  paramNodes,
		blockNodes:  blockNodes,
		interpreter: interpreter,
		readerPos:   node.ReaderPos(),
	}, nil
}
