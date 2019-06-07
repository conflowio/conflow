// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parameter

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/parsley/parsley"
)

func TransformNode(
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

	return NewNode(blockID, nameNode, valueNode, isDeclaration), nil
}
