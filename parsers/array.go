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
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"

	"github.com/opsidian/conflow/basil/schema"
)

// Array will match an array expression defined by the following rule, where P is the input parser:
//   S -> "[" "]"
//   S -> "[" P ("," P)* "]"
func Array(p parsley.Parser) *combinator.Sequence {
	return combinator.SeqOf(
		terminal.Rune('['),
		text.LeftTrim(SepByComma(p), text.WsSpacesNl),
		text.LeftTrim(terminal.Rune(']'), text.WsSpacesNl),
	).Name("array").Token("ARRAY").Bind(arrayInterpreter{})
}

type arrayInterpreter struct{}

func (a arrayInterpreter) Eval(userCtx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	panic("Eval should not be called on a raw array node")
}

func (a arrayInterpreter) TransformNode(userCtx interface{}, node parsley.Node) (parsley.Node, parsley.Error) {
	var nodes []parsley.Node
	if itemsNode, ok := node.(parsley.NonTerminalNode).Children()[1].(parsley.NonTerminalNode); ok {
		nodes = itemsNode.Children()
	}
	items := make([]parsley.Node, (len(nodes)+1)/2)
	var err parsley.Error
	for i := 0; i < len(nodes); i += 2 {
		if items[i/2], err = parsley.Transform(userCtx, nodes[i]); err != nil {
			return nil, err
		}
	}

	return &arrayNode{
		items:     items,
		pos:       node.Pos(),
		readerPos: node.ReaderPos(),
	}, nil
}

type arrayNode struct {
	items     []parsley.Node
	pos       parsley.Pos
	readerPos parsley.Pos
	schema    schema.Schema
}

// Token returns with the node's token
func (a *arrayNode) Token() string {
	return "ARRAY"
}

// Schema returns the schema for the node's value
func (a *arrayNode) Schema() interface{} {
	return a.schema
}

// StaticCheck runs static analysis on the node
func (a *arrayNode) StaticCheck(ctx interface{}) parsley.Error {
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

	a.schema = &schema.Array{Items: s}

	return nil
}

// Value creates a new block
func (a *arrayNode) Value(userCtx interface{}) (interface{}, parsley.Error) {
	if len(a.items) == 0 {
		return []interface{}{}, nil
	}

	res := make([]interface{}, len(a.items))
	for i, item := range a.items {
		value, err := parsley.EvaluateNode(userCtx, item)
		if err != nil {
			return nil, err
		}
		res[i] = value
	}
	return res, nil
}

// Pos returns with the node's position
func (a *arrayNode) Pos() parsley.Pos {
	return a.pos
}

// ReaderPos returns with the reader's position
func (a *arrayNode) ReaderPos() parsley.Pos {
	return a.readerPos
}

// SetReaderPos amends the reader position using the given function
func (a *arrayNode) SetReaderPos(f func(parsley.Pos) parsley.Pos) {
	a.readerPos = f(a.readerPos)
}

// Children returns with the array items
func (a *arrayNode) Children() []parsley.Node {
	return a.items
}

// String returns with a string representation of the node
func (a *arrayNode) String() string {
	return fmt.Sprintf("%s{%s, %d..%d}", a.Token(), a.items, a.pos, a.readerPos)
}
