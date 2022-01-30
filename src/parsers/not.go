// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers

import (
	"github.com/conflowio/parsley/ast"
	"github.com/conflowio/parsley/combinator"
	"github.com/conflowio/parsley/data"
	"github.com/conflowio/parsley/parser"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/text"
	"github.com/conflowio/parsley/text/terminal"

	"github.com/conflowio/conflow/src/schema"
)

// Not will match a logical not expression defined by the following rule, where P is the input parser:
//   S -> "!"? P
func Not(p parsley.Parser) parser.Func {
	notp := combinator.SeqOf(
		combinator.SuppressError(combinator.Optional(terminal.Rune('!'))),
		text.LeftTrim(p, text.WsSpaces),
	).Bind(notInterpreter{})

	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		res, cp, err := notp.Parse(ctx, leftRecCtx, pos)
		if err != nil {
			return nil, cp, err
		}
		nodes := res.(parsley.NonTerminalNode).Children()
		if _, ok := nodes[0].(ast.EmptyNode); ok {
			return nodes[1], cp, err
		}

		return res, cp, err
	})
}

type notInterpreter struct{}

func (n notInterpreter) StaticCheck(ctx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	nodes := node.Children()
	nodeSchema := nodes[1].Schema().(schema.Schema)
	if nodeSchema.Type() != schema.TypeBoolean {
		return nil, parsley.NewErrorf(nodes[0].Pos(), "unsupported ! operation on %s", string(nodeSchema.Type()))
	}

	return schema.BooleanValue(), nil
}

func (n notInterpreter) Eval(ctx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	nodes := node.Children()
	v, err := parsley.EvaluateNode(ctx, nodes[1])
	if err != nil {
		return nil, err
	}

	return !v.(bool), nil
}
