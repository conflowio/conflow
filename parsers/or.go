// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers

import (
	"errors"

	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text/terminal"

	"github.com/conflowio/conflow/conflow/schema"
)

// Or will match a logical or expression defined by the following rule, where P is the input parser:
//   S -> P ("||" P)*
func Or(p parsley.Parser) *combinator.Sequence {
	return SepByOp(
		p,
		terminal.Op("||"),
	).Token("COMPARE").Bind(orInterpreter{}).HandleResult(combinator.ReturnSingle())
}

type orInterpreter struct{}

func (o orInterpreter) StaticCheck(ctx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	nodes := node.Children()
	for i := 0; i < len(nodes); i += 2 {
		if err := schema.BooleanValue().ValidateSchema(nodes[i].Schema().(schema.Schema), false); err != nil {
			return nil, parsley.NewError(nodes[i].Pos(), err)
		}
	}

	return schema.BooleanValue(), nil
}

func (o orInterpreter) Eval(ctx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	nodes := node.Children()
	res := false
	for i := 0; i < len(nodes); i += 2 {
		v, err := parsley.EvaluateNode(ctx, nodes[i])
		if err != nil {
			return nil, err
		}
		switch vt := v.(type) {
		case bool:
			res = res || vt
		default:
			return nil, parsley.NewError(nodes[i].Pos(), errors.New("was expecting boolean"))
		}
	}
	return res, nil
}
