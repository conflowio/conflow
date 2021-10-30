// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers

import (
	"fmt"

	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text/terminal"

	"github.com/opsidian/conflow/basil/schema"
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
	).Token("PROD_MOD").Bind(prodModInterpreter{}).HandleResult(combinator.ReturnSingle())
}

type prodModInterpreter struct{}

func (prodModInterpreter) StaticCheck(ctx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	var resultSchema schema.Schema
	var op rune
	var opPos parsley.Pos
	expectsOp := false

	for i, node := range node.Children() {
		if i == 0 {
			resultSchema = node.Schema().(schema.Schema)
		} else if expectsOp {
			op = node.(parsley.LiteralNode).Value().(rune)
			opPos = node.Pos()
		} else {
			s := node.Schema().(schema.Schema)
			switch op {
			case '*', '/':
				if (resultSchema.Type() != schema.TypeInteger && resultSchema.Type() != schema.TypeNumber) ||
					(s.Type() != schema.TypeInteger && s.Type() != schema.TypeNumber) {
					return nil, parsley.NewErrorf(opPos, "unsupported %s operation on %s and %s", string(op), resultSchema.TypeString(), s.TypeString())
				}
				if op == '*' && resultSchema.Type() == schema.TypeInteger && s.Type() == schema.TypeInteger {
					resultSchema = schema.IntegerValue()
				} else {
					resultSchema = schema.NumberValue()
				}
			default: // '%'
				if resultSchema.Type() != schema.TypeInteger || s.Type() != schema.TypeInteger {
					return nil, parsley.NewErrorf(opPos, "unsupported %s operation on %s and %s", string(op), resultSchema.TypeString(), s.TypeString())
				}
				resultSchema = schema.IntegerValue()
			}
		}
		expectsOp = !expectsOp
	}

	return resultSchema, nil
}

func (p prodModInterpreter) Eval(ctx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	nodes := node.Children()
	var res interface{}
	var op rune
	var opPos parsley.Pos
	expectsOp := false
	for i, node := range nodes {
		v, err := parsley.EvaluateNode(ctx, node)
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
						return nil, parsley.NewErrorf(node.Pos(), "division by zero")
					}
					switch rest := res.(type) {
					case int64:
						if rest%vt == 0 {
							res = int64(float64(rest) / float64(vt))
						} else {
							res = float64(rest) / float64(vt)
						}
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
						return nil, parsley.NewErrorf(node.Pos(), "division by zero")
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
