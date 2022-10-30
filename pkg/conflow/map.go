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

func NewMapNode(keys []string, items []parsley.Node, pos, readerPos parsley.Pos, s schema.Schema) *MapNode {
	return &MapNode{
		keys:      keys,
		items:     items,
		pos:       pos,
		readerPos: readerPos,
		schema:    s,
	}
}

type MapNode struct {
	keys      []string
	items     []parsley.Node
	pos       parsley.Pos
	readerPos parsley.Pos
	schema    schema.Schema
}

// Token returns with the node's token
func (a *MapNode) Token() string {
	return "MAP"
}

// Schema returns the schema for the node's value
func (a *MapNode) Schema() interface{} {
	return a.schema
}

// StaticCheck runs static analysis on the node
func (a *MapNode) StaticCheck(ctx interface{}) parsley.Error {
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
func (a *MapNode) Value(userCtx interface{}) (interface{}, parsley.Error) {
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
func (a *MapNode) Pos() parsley.Pos {
	return a.pos
}

// ReaderPos returns with the reader's position
func (a *MapNode) ReaderPos() parsley.Pos {
	return a.readerPos
}

// SetReaderPos amends the reader position using the given function
func (a *MapNode) SetReaderPos(f func(parsley.Pos) parsley.Pos) {
	a.readerPos = f(a.readerPos)
}

// Keys returns the map keys
func (a *MapNode) Keys() []string {
	return a.keys
}

// Children returns with the array items
func (a *MapNode) Children() []parsley.Node {
	return a.items
}

// String returns with a string representation of the node
func (a *MapNode) String() string {
	return fmt.Sprintf("%s{%s, %d..%d}", a.Token(), a.items, a.pos, a.readerPos)
}
