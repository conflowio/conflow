// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers

import (
	"fmt"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text/terminal"
)

// Sum will match +, - arithmetic operations defined by the following rule, where P is the input parser:
//   S      -> P (SUM_OP P)*
//   SUM_OP -> "+"
//          -> "-"
func Sum(p parsley.Parser) parser.Func {
	return combinator.Single(
		SepByOp(
			p,
			combinator.Choice(
				terminal.Rune('+'),
				terminal.Rune('-'),
			),
		).Bind(ast.InterpreterFunc(evalSum)),
	)
}

func evalSum(ctx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	nodes := node.Children()
	var res interface{}
	var op rune
	var opPos parsley.Pos
	expectsOp := false
	modifier := int64(1)
	for i, node := range nodes {
		v, err := node.Value(ctx)
		if err != nil {
			return nil, err
		}
		if i == 0 {
			res = v
		} else if expectsOp {
			op = v.(rune)
			opPos = node.Pos()
			if op == '+' {
				modifier = int64(1)
			} else {
				modifier = int64(-1)
			}
		} else {
			switch vt := v.(type) {
			case int64:
				switch rest := res.(type) {
				case int64:
					res = rest + modifier*vt
				case float64:
					res = rest + float64(modifier)*float64(vt)
				default:
					return nil, parsley.NewErrorf(opPos, "unsupported %s operation on %s and %s", string(op), fmt.Sprintf("%T", res), fmt.Sprintf("%T", v))
				}
			case float64:
				switch rest := res.(type) {
				case int64:
					res = float64(rest) + float64(modifier)*vt
				case float64:
					res = rest + float64(modifier)*vt
				default:
					return nil, parsley.NewErrorf(opPos, "unsupported %s operation on %s and %s", string(op), fmt.Sprintf("%T", res), fmt.Sprintf("%T", v))
				}
			case string:
				if op != '+' {
					return nil, parsley.NewErrorf(opPos, "unsupported %s operation on %s and %s", string(op), fmt.Sprintf("%T", res), fmt.Sprintf("%T", v))
				}
				switch rest := res.(type) {
				case string:
					res = rest + vt
				default:
					return nil, parsley.NewErrorf(opPos, "unsupported %s operation on %s and %s", string(op), fmt.Sprintf("%T", res), fmt.Sprintf("%T", v))
				}
			default:
				return nil, parsley.NewErrorf(opPos, "unsupported %s operation on %s and %s", string(op), fmt.Sprintf("%T", res), fmt.Sprintf("%T", v))
			}
		}
		expectsOp = !expectsOp
	}
	return res, nil
}
