// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package terminal

import (
	"fmt"

	"github.com/conflowio/parsley/data"
	"github.com/conflowio/parsley/parser"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/text"
)

// NilNode is a leaf node in the AST
type NilNode struct {
	schema    interface{}
	pos       parsley.Pos
	readerPos parsley.Pos
}

// NewNilNode creates a new NilNode instance
func NewNilNode(schema interface{}, pos parsley.Pos, readerPos parsley.Pos) *NilNode {
	return &NilNode{
		schema:    schema,
		pos:       pos,
		readerPos: readerPos,
	}
}

// Token returns with the node token
func (n *NilNode) Token() string {
	return "NIL"
}

// Schema returns the schema for the node's value
func (n *NilNode) Schema() interface{} {
	return n.schema
}

// Value returns with the value of the node
func (n *NilNode) Value() interface{} {
	return nil
}

// Pos returns the position
func (n *NilNode) Pos() parsley.Pos {
	return n.pos
}

// ReaderPos returns the position of the first character immediately after this node
func (n *NilNode) ReaderPos() parsley.Pos {
	return n.readerPos
}

// SetReaderPos changes the reader position
func (n *NilNode) SetReaderPos(fun func(parsley.Pos) parsley.Pos) {
	n.readerPos = fun(n.readerPos)
}

// String returns with a string representation of the node
func (n *NilNode) String() string {
	return fmt.Sprintf("%s{%d..%d}", n.Token(), n.pos, n.readerPos)
}

// Nil matches a nil literal
func Nil(schema interface{}, nilStr string) parser.Func {
	if nilStr == "" {
		panic("Nil() should not be called with an empty nilStr")
	}

	notFoundErr := parsley.NotFoundError(nilStr)

	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		tr := ctx.Reader().(*text.Reader)
		if readerPos, found := tr.MatchWord(pos, nilStr); found {
			return NewNilNode(schema, pos, readerPos), data.EmptyIntSet, nil
		}

		return nil, data.EmptyIntSet, parsley.NewError(pos, notFoundErr)
	})
}
