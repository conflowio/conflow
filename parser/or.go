// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser

import (
	"github.com/opsidian/basil/variable"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text/terminal"
)

// Or will match a logical or expression defined by the following rule, where P is the input parser:
//   S -> P ("||" P)*
func Or(p parsley.Parser) parser.Func {
	return combinator.Single(
		SepByOp(
			p,
			terminal.Op("||"),
		).Bind(orInterpreter{}),
	)
}

type orInterpreter struct{}

func (o orInterpreter) StaticCheck(ctx interface{}, node parsley.NonTerminalNode) (string, parsley.Error) {
	nodes := node.Children()
	for i := 0; i < len(nodes); i += 2 {
		if err := variable.CheckNodeType(nodes[i], variable.TypeBool); err != nil {
			return "", err
		}
	}

	return variable.TypeBool, nil
}

func (o orInterpreter) Eval(ctx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	nodes := node.Children()
	res := false
	for i := 0; i < len(nodes); i += 2 {
		v, err := nodes[i].Value(ctx)
		if err != nil {
			return nil, err
		}
		switch vt := v.(type) {
		case bool:
			res = res || vt
		default:
			return nil, parsley.NewError(nodes[i].Pos(), variable.ErrExpectingBool)
		}
	}
	return res, nil
}
