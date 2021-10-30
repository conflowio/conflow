// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parameter

import (
	"errors"

	"github.com/opsidian/parsley/parsley"

	"github.com/opsidian/conflow/basil/directive"
	"github.com/opsidian/conflow/conflow"
)

func TransformNode(
	parseCtx interface{},
	node parsley.Node,
	blockID conflow.ID,
	paramNames map[conflow.ID]struct{},
) (conflow.ParameterNode, parsley.Error) {
	nodes := node.(parsley.NonTerminalNode).Children()

	var directives []conflow.BlockNode
	if n, ok := nodes[0].(parsley.NonTerminalNode); ok && len(n.Children()) > 0 {
		var err parsley.Error
		var deps conflow.Dependencies
		if directives, deps, err = directive.Transform(parseCtx, n.Children()); err != nil {
			return nil, err
		}
		if len(deps) > 0 {
			return nil, parsley.NewError(n.Pos(), errors.New("a parameter directive can not have dependencies"))
		}
	}

	nameNode := nodes[1].(*conflow.IDNode)
	if _, exists := paramNames[nameNode.ID()]; exists {
		return nil, parsley.NewErrorf(
			nodes[0].Pos(),
			"%q parameter was defined multiple times", nameNode.ID(),
		)
	}
	paramNames[nameNode.ID()] = struct{}{}

	op := nodes[2].(parsley.LiteralNode).Value()
	isDeclaration := op == ":="

	valueNode, err := parsley.Transform(parseCtx, nodes[3])
	if err != nil {
		return nil, err
	}

	return NewNode(blockID, nameNode, valueNode, isDeclaration, directives), nil
}
