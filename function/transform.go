// Copyright (c) 2018 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package function

import (
	"fmt"

	"github.com/opsidian/basil/variable"
	"github.com/opsidian/parsley/parsley"
)

// Registry contains a list of function interpreters and behaves as a node transformer registry
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

// TransformNode returns with a node transformer function for a function
func TransformNode(registry parsley.NodeTransformerRegistry) parsley.NodeTransformFunc {
	return parsley.NodeTransformFunc(func(userCtx interface{}, node parsley.Node) (parsley.Node, parsley.Error) {
		nodes := node.(parsley.NonTerminalNode).Children()
		nameNode := nodes[0]
		name, _ := nameNode.Value(nil)

		transformer, exists := registry.NodeTransformer(string(name.(variable.ID)))
		if !exists {
			return nil, parsley.NewError(nameNode.Pos(), fmt.Errorf("%q function does not exist", name))
		}

		return transformer.TransformNode(userCtx, node)
	})
}

func transformNode(
	userCtx interface{},
	node parsley.Node,
	interpreter Interpreter,
) (parsley.Node, parsley.Error) {
	nodes := node.(parsley.NonTerminalNode).Children()
	nameNode := nodes[0]

	argumentsNode := nodes[2].(parsley.NonTerminalNode)
	var argumentNodes []parsley.Node
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
	return &Node{
		nameNode:      nameNode,
		argumentNodes: argumentNodes,
		readerPos:     node.ReaderPos(),
		interpreter:   interpreter,
	}, nil
}
