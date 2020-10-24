// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers

import (
	"reflect"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

// TernaryIf will match a ternary if expression defined by the following rule, where P is the input parser:
//   S -> P
//     -> P "?" P ":" P
func TernaryIf(p parsley.Parser) *combinator.Sequence {
	parsers := []parsley.Parser{
		p,
		combinator.SuppressError(text.LeftTrim(terminal.Rune('?'), text.WsSpacesNl)),
		text.LeftTrim(p, text.WsSpacesNl),
		text.LeftTrim(terminal.Rune(':'), text.WsSpacesNl),
		text.LeftTrim(p, text.WsSpacesNl),
	}

	lookup := func(i int) parsley.Parser {
		if i < len(parsers) {
			return parsers[i]
		}
		return nil
	}
	l := len(parsers)
	lenCheck := func(len int) bool {
		return len == 1 || len == l
	}
	return combinator.Seq(
		"TERNARY_IF", lookup, lenCheck,
	).Bind(ast.InterpreterFunc(evalTernaryIf)).ReturnSingle()
}

func evalTernaryIf(ctx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	nodes := node.Children()
	cond, err := nodes[0].Value(ctx)
	if err != nil {
		return nil, err
	}
	switch cond {
	case true:
		return nodes[2].Value(ctx)
	case false:
		return nodes[4].Value(ctx)
	default:
		return nil, parsley.NewErrorf(nodes[0].Pos(), "expecting bool, got %s", reflect.ValueOf(cond).Kind())
	}
}
