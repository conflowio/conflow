// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers

import (
	"errors"
	"reflect"

	"github.com/conflowio/parsley/combinator"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/text"
	"github.com/conflowio/parsley/text/terminal"

	"github.com/conflowio/conflow/src/schema"
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
	).Bind(ternaryInterpreter{}).HandleResult(combinator.ReturnSingle())
}

type ternaryInterpreter struct{}

func (t ternaryInterpreter) StaticCheck(ctx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	nodes := node.Children()

	if nodes[0].Schema().(schema.Schema).Type() != schema.TypeBoolean {
		return nil, parsley.NewError(nodes[0].Pos(), errors.New("must be boolean"))
	}

	s, err := schema.GetCommonSchema(nodes[2].Schema().(schema.Schema), nodes[4].Schema().(schema.Schema))
	if err != nil {
		return nil, parsley.NewErrorf(
			nodes[2].Pos(),
			"both expressions must have the same type, but got %s and %s",
			nodes[2].Schema().(schema.Schema).TypeString(),
			nodes[4].Schema().(schema.Schema).TypeString(),
		)
	}

	return s, nil
}

func (t ternaryInterpreter) Eval(ctx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	nodes := node.Children()
	cond, err := parsley.EvaluateNode(ctx, nodes[0])
	if err != nil {
		return nil, err
	}
	switch cond {
	case true:
		return parsley.EvaluateNode(ctx, nodes[2])
	case false:
		return parsley.EvaluateNode(ctx, nodes[4])
	default:
		return nil, parsley.NewErrorf(nodes[0].Pos(), "expecting bool, got %s", reflect.ValueOf(cond).Kind())
	}
}
