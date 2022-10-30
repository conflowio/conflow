// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conflow

import (
	"fmt"

	"github.com/conflowio/parsley/parsley"

	"github.com/conflowio/conflow/pkg/schema"
)

func NewArrayNode(items []parsley.Node, pos, readerPos parsley.Pos, schema schema.Schema) *ArrayNode {
	return &ArrayNode{
		items:     items,
		pos:       pos,
		readerPos: readerPos,
		schema:    schema,
	}
}

type ArrayNode struct {
	items     []parsley.Node
	pos       parsley.Pos
	readerPos parsley.Pos
	schema    schema.Schema
}

// Token returns with the node's token
func (a *ArrayNode) Token() string {
	return "ARRAY"
}

// Schema returns the schema for the node's value
func (a *ArrayNode) Schema() interface{} {
	return a.schema
}

// StaticCheck runs static analysis on the node
func (a *ArrayNode) StaticCheck(ctx interface{}) parsley.Error {
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
func (a *ArrayNode) Value(userCtx interface{}) (interface{}, parsley.Error) {
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
func (a *ArrayNode) Pos() parsley.Pos {
	return a.pos
}

// ReaderPos returns with the reader's position
func (a *ArrayNode) ReaderPos() parsley.Pos {
	return a.readerPos
}

// SetReaderPos amends the reader position using the given function
func (a *ArrayNode) SetReaderPos(f func(parsley.Pos) parsley.Pos) {
	a.readerPos = f(a.readerPos)
}

// Children returns with the array items
func (a *ArrayNode) Children() []parsley.Node {
	return a.items
}

// String returns with a string representation of the node
func (a *ArrayNode) String() string {
	return fmt.Sprintf("%s{%s, %d..%d}", a.Token(), a.items, a.pos, a.readerPos)
}
