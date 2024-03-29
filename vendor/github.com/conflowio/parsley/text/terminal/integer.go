// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package terminal

import (
	"fmt"
	"strconv"

	"github.com/conflowio/parsley/data"
	"github.com/conflowio/parsley/parser"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/text"
)

// IntegerNode is a leaf node in the AST
type IntegerNode struct {
	schema    interface{}
	value     int64
	pos       parsley.Pos
	readerPos parsley.Pos
}

// NewIntegerNode creates a new IntegerNode instance
func NewIntegerNode(schema interface{}, value int64, pos parsley.Pos, readerPos parsley.Pos) *IntegerNode {
	return &IntegerNode{
		schema:    schema,
		value:     value,
		pos:       pos,
		readerPos: readerPos,
	}
}

// Token returns with the node token
func (i *IntegerNode) Token() string {
	return "INTEGER"
}

// Schema returns the schema for the node's value
func (i *IntegerNode) Schema() interface{} {
	return i.schema
}

// Value returns with the value of the node
func (i *IntegerNode) Value() interface{} {
	return i.value
}

// Pos returns the position
func (i *IntegerNode) Pos() parsley.Pos {
	return i.pos
}

// ReaderPos returns the position of the first character immediately after this node
func (i *IntegerNode) ReaderPos() parsley.Pos {
	return i.readerPos
}

// SetReaderPos changes the reader position
func (i *IntegerNode) SetReaderPos(fun func(parsley.Pos) parsley.Pos) {
	i.readerPos = fun(i.readerPos)
}

// String returns with a string representation of the node
func (i *IntegerNode) String() string {
	return fmt.Sprintf("%s{%v, %d..%d}", i.Token(), i.value, i.pos, i.readerPos)
}

// Integer matches all integer numbers and zero with an optional -/+ sign
func Integer(schema interface{}) parser.Func {
	notFoundErr := parsley.NotFoundError("integer value")

	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		tr := ctx.Reader().(*text.Reader)
		if readerPos, result := tr.ReadRegexp(pos, "[-+]?(?:[1-9][0-9]*|0[xX][0-9a-fA-F]+|0[0-7]*)"); result != nil {
			if _, isFloat := tr.ReadRune(readerPos, '.'); isFloat {
				return nil, data.EmptyIntSet, parsley.NewError(pos, notFoundErr)
			}
			intValue, err := strconv.ParseInt(string(result), 0, 64)
			if err != nil {
				// This should never happen
				panic(fmt.Sprintf("Could not convert %s to integer", string(result)))
			}
			return NewIntegerNode(schema, intValue, pos, readerPos), data.EmptyIntSet, nil
		}
		return nil, data.EmptyIntSet, parsley.NewError(pos, notFoundErr)
	})
}
