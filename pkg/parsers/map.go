// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers

import (
	"github.com/conflowio/parsley/ast"
	"github.com/conflowio/parsley/combinator"
	"github.com/conflowio/parsley/parser"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/text"
	"github.com/conflowio/parsley/text/terminal"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/schema"
)

// Map will match an map expression defined by the following rule, where P is the input parser:
//
//	S -> "map" "{" "}"
//	S -> "map" "{"
//	        (STRING ":" P ",")*
//	     "}"
func Map(p parsley.Parser) parser.Func {
	keyValue := combinator.SeqOf(
		terminal.String(schema.StringValue(), false),
		text.LeftTrim(terminal.Rune(':'), text.WsSpaces),
		text.LeftTrim(p, text.WsSpaces),
	).Name("key-value pair")

	emptyMap := combinator.SeqOf(
		terminal.Word(schema.TypeString, "map", "map"),
		text.LeftTrim(terminal.Rune('{'), text.WsSpaces),
		text.LeftTrim(terminal.Rune('}'), text.WsSpacesNl),
	).Token("MAP").Bind(mapInterpreter{})

	nonEmptyMap := combinator.SeqOf(
		terminal.Word(schema.TypeString, "map", "map"),
		text.LeftTrim(terminal.Rune('{'), text.WsSpaces),
		text.LeftTrim(SepByComma(keyValue), text.WsSpacesNl),
		text.LeftTrim(terminal.Rune('}'), text.WsSpacesNl),
	).Name("map").Token("MAP").Bind(mapInterpreter{})

	return combinator.Choice(
		emptyMap,
		nonEmptyMap,
	).Name("map")
}

type mapInterpreter struct{}

func (a mapInterpreter) Eval(userCtx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	panic("Eval should not be called on a raw array node")
}

func (a mapInterpreter) TransformNode(userCtx interface{}, node parsley.Node) (parsley.Node, parsley.Error) {
	if len(node.(parsley.NonTerminalNode).Children()) == 3 {
		return conflow.NewMapNode(nil, nil, node.Pos(), node.ReaderPos(), schema.NullValue()), nil
	}

	var nodes []parsley.Node
	if itemsNode, ok := node.(parsley.NonTerminalNode).Children()[2].(parsley.NonTerminalNode); ok {
		nodes = itemsNode.Children()
	}

	keys := make([]string, (len(nodes)+1)/2)
	items := make([]parsley.Node, (len(nodes)+1)/2)
	for i := 0; i < len(nodes); i += 2 {
		kvNodes := nodes[i].(*ast.NonTerminalNode).Children()

		keys[i/2] = kvNodes[0].(parsley.LiteralNode).Value().(string)

		var err parsley.Error
		if items[i/2], err = parsley.Transform(userCtx, kvNodes[2]); err != nil {
			return nil, err
		}
	}

	return conflow.NewMapNode(keys, items, node.Pos(), node.ReaderPos(), nil), nil
}
