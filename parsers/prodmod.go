// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers

import (
	"fmt"

	"github.com/opsidian/basil/basil/schema"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text/terminal"
)

// ProdMod will match *, /, % arithmetic operations defined by the following rule, where P is the input parser:
//   S           -> P (PROD_MOD_OP P)*
//   PROD_MOD_OP -> "*"
//               -> "/"
//               -> "%"
func ProdMod(p parsley.Parser) *combinator.Sequence {
	return SepByOp(
		p,
		combinator.Choice(
			terminal.Rune('*'),
			terminal.Rune('/'),
			terminal.Rune('%'),
		),
	).Token("PROD_MOD").Bind(ast.InterpreterFunc(evalProdMod)).ReturnSingle()
}

func evalProdMod(ctx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	nodes := node.Children()
	var res interface{}
	var op rune
	var opPos parsley.Pos
	expectsOp := false
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
		} else {
			switch vt := v.(type) {
			case int64:
				if op == '*' {
					switch rest := res.(type) {
					case int64:
						res = rest * vt
					case float64:
						res = rest * float64(vt)
					default:
						return nil, parsley.NewErrorf(opPos, "unsupported %s operation on %s and %s", string(op), fmt.Sprintf("%T", res), fmt.Sprintf("%T", v))
					}
				} else if op == '/' {
					if vt == 0 {
						return nil, parsley.NewErrorf(node.Pos(), "divison by zero")
					}
					switch rest := res.(type) {
					case int64:
						res = rest / vt
					case float64:
						res = rest / float64(vt)
					default:
						return nil, parsley.NewErrorf(opPos, "unsupported %s operation on %s and %s", string(op), fmt.Sprintf("%T", res), fmt.Sprintf("%T", v))
					}
				} else {
					switch rest := res.(type) {
					case int64:
						res = rest % vt
					default:
						return nil, parsley.NewErrorf(opPos, "unsupported %s operation on %s and %s", string(op), fmt.Sprintf("%T", res), fmt.Sprintf("%T", v))
					}
				}
			case float64:
				if op == '*' {
					switch rest := res.(type) {
					case int64:
						res = float64(rest) * vt
					case float64:
						res = rest * vt
					default:
						return nil, parsley.NewErrorf(opPos, "unsupported %s operation on %s and %s", string(op), fmt.Sprintf("%T", res), fmt.Sprintf("%T", v))
					}
				} else if op == '/' {
					if 0.0-vt < schema.Epsilon && vt-0.0 < schema.Epsilon {
						return nil, parsley.NewErrorf(node.Pos(), "divison by zero")
					}
					switch rest := res.(type) {
					case int64:
						res = float64(rest) / vt
					case float64:
						res = rest / vt
					default:
						return nil, parsley.NewErrorf(opPos, "unsupported %s operation on %s and %s", string(op), fmt.Sprintf("%T", res), fmt.Sprintf("%T", v))
					}
				} else {
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
