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

	"github.com/conflowio/conflow/conflow/schema"
)

// Sum will match +, - arithmetic operations defined by the following rule, where P is the input parser:
//   S      -> P (SUM_OP P)*
//   SUM_OP -> "+"
//          -> "-"
func Sum(p parsley.Parser) *combinator.Sequence {
	return SepByOp(
		p,
		combinator.Choice(
			terminal.Rune('+'),
			terminal.Rune('-'),
		),
	).Token("SUM").Bind(sumInterpreter{}).HandleResult(combinator.ReturnSingle())
}

type sumInterpreter struct{}

func (s sumInterpreter) StaticCheck(ctx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
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

			switch resultSchema.Type() {
			case schema.TypeString:
				if s.Type() != schema.TypeString {
					return nil, parsley.NewErrorf(opPos, "unsupported %s operation on %s and %s", string(op), resultSchema.TypeString(), s.TypeString())
				}
				resultSchema = schema.StringValue()
			case schema.TypeInteger, schema.TypeNumber:
				if s.Type() != schema.TypeInteger && s.Type() != schema.TypeNumber {
					return nil, parsley.NewErrorf(opPos, "unsupported %s operation on %s and %s", string(op), resultSchema.TypeString(), s.TypeString())
				}
				if resultSchema.Type() == schema.TypeInteger && s.Type() == schema.TypeInteger {
					resultSchema = schema.IntegerValue()
				} else {
					resultSchema = schema.NumberValue()
				}
			default:
				return nil, parsley.NewErrorf(opPos, "unsupported %s operation on %s and %s", string(op), resultSchema.TypeString(), s.TypeString())
			}
		}
		expectsOp = !expectsOp
	}

	return resultSchema, nil
}

func (s sumInterpreter) Eval(ctx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	nodes := node.Children()
	var res interface{}
	var op rune
	var opPos parsley.Pos
	expectsOp := false
	modifier := int64(1)
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
