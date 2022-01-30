// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers

import (
	"errors"

	"github.com/conflowio/parsley/combinator"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/text/terminal"

	"github.com/conflowio/conflow/src/conflow/schema"
)

// And will match a logical and expression defined by the following rule, where P is the input parser:
//   S -> P ("&&" P)*
func And(p parsley.Parser) *combinator.Sequence {
	return SepByOp(
		p,
		terminal.Op("&&"),
	).Token("AND").Bind(andInterpreter{}).HandleResult(combinator.ReturnSingle())
}

type andInterpreter struct{}

func (a andInterpreter) StaticCheck(ctx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	nodes := node.Children()
	for i := 0; i < len(nodes); i += 2 {
		if err := schema.BooleanValue().ValidateSchema(nodes[i].Schema().(schema.Schema), false); err != nil {
			return nil, parsley.NewError(nodes[i].Pos(), err)
		}
	}

	return schema.BooleanValue(), nil
}

func (a andInterpreter) Eval(ctx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	nodes := node.Children()
	res := true
	for i := 0; i < len(nodes); i += 2 {
		v, err := parsley.EvaluateNode(ctx, nodes[i])
		if err != nil {
			return nil, err
		}
		switch vt := v.(type) {
		case bool:
			res = res && vt
		default:
			return nil, parsley.NewError(nodes[i].Pos(), errors.New("was expecting boolean"))
		}
	}
	return res, nil
}
