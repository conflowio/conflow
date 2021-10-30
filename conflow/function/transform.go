// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package function

import (
	"github.com/opsidian/conflow/conflow"
	"github.com/opsidian/parsley/parsley"
)

func transformNode(
	userCtx interface{},
	node parsley.Node,
	interpreter conflow.FunctionInterpreter,
) (parsley.Node, parsley.Error) {
	nodes := node.(parsley.NonTerminalNode).Children()
	nameNode := nodes[0]

	var argumentNodes []parsley.Node

	if argumentsNode, ok := nodes[2].(parsley.NonTerminalNode); ok {
		children := argumentsNode.Children()
		childrenCount := len(children)
		if childrenCount > 0 {
			argumentNodes = make([]parsley.Node, childrenCount/2+1)
			var err parsley.Error
			for i := 0; i < childrenCount; i += 2 {
				if argumentNodes[i/2], err = parsley.Transform(userCtx, children[i]); err != nil {
					return nil, err
				}
			}
		}
	}

	return &Node{
		nameNode:      nameNode,
		argumentNodes: argumentNodes,
		readerPos:     node.ReaderPos(),
		interpreter:   interpreter,
	}, nil
}
