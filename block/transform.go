// Copyright (c) 2018 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package block

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/identifier"
	"github.com/opsidian/parsley/parsley"
)

// Registry contains a list of block interpreters and behaves as a node transformer registry
type Registry map[string]Interpreter

// NodeTransformer returns with the named node transformer
func (r Registry) NodeTransformer(name string) (parsley.NodeTransformer, bool) {
	interpreter, exists := r[name]
	if !exists {
		return nil, false
	}

	return parsley.NodeTransformFunc(func(userCtx interface{}, node parsley.Node) (parsley.Node, parsley.Error) {
		return transformNode(userCtx, node, interpreter)
	}), true
}

func transformNode(
	userCtx interface{},
	node parsley.Node,
	interpreter Interpreter,
) (parsley.Node, parsley.Error) {
	if interpreter.BlockRegistry() != nil {
		parentCtx := userCtx.(*basil.Context)
		userCtx = basil.NewContext(parentCtx, basil.ContextConfig{
			BlockRegistry: interpreter.BlockRegistry(),
		})
	}

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
						childBlock, err := blockChild.(parsley.Transformable).Transform(userCtx)
						if err != nil {
							return nil, err
						}
						blockNodes = append(blockNodes, childBlock.(*Node))
					} else {
						children := blockChild.(parsley.NonTerminalNode).Children()
						paramNode := children[0]
						paramName, _ := paramNode.Value(nil)
						if _, alreadyExists := paramNodes[paramName.(basil.ID)]; alreadyExists {
							return nil, parsley.NewErrorf(children[0].Pos(), "%q parameter was defined multiple times", paramName)
						}
						valueNode, err := parsley.Transform(userCtx, children[2])
						if err != nil {
							return nil, err
						}

						paramNodes[paramName.(basil.ID)] = NewParamNode(paramNode, valueNode)
					}
				}
			}
		} else { // We have an expression as the value of the block
			valueParamName := interpreter.ValueParamName()
			if valueParamName == "" {
				blockType, _ := typeNode.Value(nil)
				return nil, parsley.NewErrorf(typeNode.Pos(), "%q block does not support short format", blockType)
			}
			valueNode, err := parsley.Transform(userCtx, blockValueNode)
			if err != nil {
				return nil, err
			}
			paramNodes = map[basil.ID]basil.BlockParamNode{
				valueParamName: NewParamNode(
					identifier.NewNode(valueParamName, blockValueNode.Pos(), blockValueNode.Pos()),
					valueNode,
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
