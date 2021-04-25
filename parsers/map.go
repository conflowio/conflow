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
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"

	"github.com/opsidian/basil/basil/schema"
)

// Map will match an map expression defined by the following rule, where P is the input parser:
//   S -> "map" "{" "}"
//   S -> "map" "{"
//           (STRING ":" P ",")*
//        "}"
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
		return &mapNode{
			keys:      nil,
			items:     nil,
			pos:       node.Pos(),
			readerPos: node.ReaderPos(),
			schema:    schema.NullValue(),
		}, nil
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

	return &mapNode{
		keys:      keys,
		items:     items,
		pos:       node.Pos(),
		readerPos: node.ReaderPos(),
	}, nil
}

type mapNode struct {
	keys      []string
	items     []parsley.Node
	pos       parsley.Pos
	readerPos parsley.Pos
	schema    schema.Schema
}

// Token returns with the node's token
func (a *mapNode) Token() string {
	return "MAP"
}

// Schema returns the schema for the node's value
func (a *mapNode) Schema() interface{} {
	return a.schema
}

// StaticCheck runs static analysis on the node
func (a *mapNode) StaticCheck(ctx interface{}) parsley.Error {
	if len(a.items) == 0 {
		a.schema = schema.NullValue()
		return nil
	}

	s, err := schema.GetSchemaForValues(len(a.items), func(i int) (schema.Schema, error) {
		return a.items[i].Schema().(schema.Schema), nil
	})

	if err != nil {
		return parsley.NewError(a.Pos(), err)
	}

	a.schema = &schema.Map{AdditionalProperties: s}

	return nil
}

// Value creates a new block
func (a *mapNode) Value(userCtx interface{}) (interface{}, parsley.Error) {
	if len(a.items) == 0 {
		return map[string]interface{}{}, nil
	}

	res := make(map[string]interface{}, len(a.items))
	for i, item := range a.items {
		value, err := parsley.EvaluateNode(userCtx, item)
		if err != nil {
			return nil, err
		}
		res[a.keys[i]] = value
	}
	return res, nil
}

// Pos returns with the node's position
func (a *mapNode) Pos() parsley.Pos {
	return a.pos
}

// ReaderPos returns with the reader's position
func (a *mapNode) ReaderPos() parsley.Pos {
	return a.readerPos
}

// SetReaderPos amends the reader position using the given function
func (a *mapNode) SetReaderPos(f func(parsley.Pos) parsley.Pos) {
	a.readerPos = f(a.readerPos)
}

// Children returns with the array items
func (a *mapNode) Children() []parsley.Node {
	return a.items
}

// String returns with a string representation of the node
func (a *mapNode) String() string {
	return fmt.Sprintf("%s{%s, %d..%d}", a.Token(), a.items, a.pos, a.readerPos)
}
