// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser

import (
	"fmt"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

// Not will match a logical not expression defined by the following rule, where P is the input parser:
//   S -> "!"? P
func Not(p parsley.Parser) parser.Func {
	notp := combinator.Seq(
		combinator.SuppressError(combinator.Optional(terminal.Rune('!'))),
		text.LeftTrim(p, text.WsSpaces),
	).Bind(ast.InterpreterFunc(evalNot))

	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		res, cp, err := notp.Parse(ctx, leftRecCtx, pos)
		if err != nil {
			return nil, cp, err
		}
		nodes := res.(*ast.NonTerminalNode).Children()
		if _, ok := nodes[0].(ast.NilNode); ok {
			return nodes[1], cp, err
		}

		return res, cp, err
	})
}

func evalNot(ctx interface{}, nodes []parsley.Node) (interface{}, parsley.Error) {
	v, err := nodes[1].Value(ctx)
	if err != nil {
		return nil, err
	}

	switch vt := v.(type) {
	case bool:
		return !vt, nil
	default:
		return nil, parsley.NewErrorf(nodes[0].Pos(), "unsupported ! operation on %s", fmt.Sprintf("%T", v))
	}
}
