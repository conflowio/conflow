// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers

import (
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text/terminal"

	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/variable"
)

// Variable will match a variable expression defined by the following rule, where P is the input parser:
//   S         -> ID "." ID
//   ID        -> /[a-z][a-z0-9]*(?:_[a-z0-9]+)*/
//
// Variable refers to a named block's parameter, in the format of `<block ID>.<parameter ID>`.
func Variable() *combinator.Sequence {
	return combinator.SeqOf(
		ID(),
		terminal.Rune('.'),
		ID(),
	).Name("variable").Token("VAR").Bind(variableInterpreter{})
}

type variableInterpreter struct{}

func (v variableInterpreter) Eval(userCtx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	panic("Eval should not be called on a raw variable node")
}

func (v variableInterpreter) TransformNode(userCtx interface{}, node parsley.Node) (parsley.Node, parsley.Error) {
	nodes := node.(parsley.NonTerminalNode).Children()
	return variable.NewNode(nodes[0].(*conflow.IDNode), nodes[2].(*conflow.IDNode)), nil
}
