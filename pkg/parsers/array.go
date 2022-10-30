// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers

import (
	"github.com/conflowio/parsley/combinator"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/text"
	"github.com/conflowio/parsley/text/terminal"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/schema"
)

// Array will match an array expression defined by the following rule, where P is the input parser:
//
//	S -> "[" "]"
//	S -> "[" P ("," P)* "]"
func Array(p parsley.Parser) *combinator.Sequence {
	return combinator.SeqOf(
		terminal.Rune('['),
		text.LeftTrim(SepByComma(p), text.WsSpacesNl),
		text.LeftTrim(terminal.Rune(']'), text.WsSpacesNl),
	).Name("array").Token("ARRAY").Bind(arrayInterpreter{})
}

type arrayInterpreter struct{}

func (a arrayInterpreter) Eval(userCtx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	panic("Eval should not be called on a raw array node")
}

func (a arrayInterpreter) TransformNode(userCtx interface{}, node parsley.Node) (parsley.Node, parsley.Error) {
	var nodes []parsley.Node
	if itemsNode, ok := node.(parsley.NonTerminalNode).Children()[1].(parsley.NonTerminalNode); ok {
		nodes = itemsNode.Children()
	}

	if len(nodes) == 0 {
		return conflow.NewArrayNode(nil, node.Pos(), node.ReaderPos(), schema.NullValue()), nil
	}

	items := make([]parsley.Node, (len(nodes)+1)/2)
	var err parsley.Error
	for i := 0; i < len(nodes); i += 2 {
		if items[i/2], err = parsley.Transform(userCtx, nodes[i]); err != nil {
			return nil, err
		}
	}

	return conflow.NewArrayNode(items, node.Pos(), node.ReaderPos(), nil), nil
}
