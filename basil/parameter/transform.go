// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parameter

import (
	"errors"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/directive"
	"github.com/opsidian/parsley/parsley"
)

func TransformNode(
	parseCtx interface{},
	node parsley.Node,
	blockID basil.ID,
	paramNames map[basil.ID]struct{},
) (basil.ParameterNode, parsley.Error) {
	nodes := node.(parsley.NonTerminalNode).Children()

	var directives []basil.BlockNode
	if n, ok := nodes[0].(parsley.NonTerminalNode); ok && len(n.Children()) > 0 {
		var err parsley.Error
		var deps basil.Dependencies
		if directives, deps, err = directive.Transform(parseCtx, n.Children()); err != nil {
			return nil, err
		}
		if len(deps) > 0 {
			return nil, parsley.NewError(n.Pos(), errors.New("a parameter directive can not have dependencies"))
		}
	}

	nameNode := nodes[1].(*basil.IDNode)
	if _, exists := paramNames[nameNode.ID()]; exists {
		return nil, parsley.NewErrorf(
			nodes[0].Pos(),
			"%q parameter was defined multiple times", nameNode.ID(),
		)
	}
	paramNames[nameNode.ID()] = struct{}{}

	op, _ := nodes[2].Value(nil)
	isDeclaration := op == ":="

	valueNode, err := parsley.Transform(parseCtx, nodes[3])
	if err != nil {
		return nil, err
	}

	return NewNode(blockID, nameNode, valueNode, isDeclaration, directives), nil
}
