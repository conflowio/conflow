// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

import (
	"fmt"

	"github.com/opsidian/parsley/parsley"

	"github.com/opsidian/conflow/basil/schema"
)

// NameNode contains a name node
// A name node consists of one identifier, or two separated by a single character
type NameNode struct {
	selector *IDNode
	name     *IDNode
	sep      parsley.LiteralNode
}

// NewNameNode creates a new name node type
func NewNameNode(selector *IDNode, sep parsley.LiteralNode, name *IDNode) *NameNode {
	return &NameNode{
		selector: selector,
		sep:      sep,
		name:     name,
	}
}

// NameNode returns with the name node
func (t *NameNode) NameNode() *IDNode {
	return t.name
}

// SelectorNode returns with selector node
// A selector can be e.g. a package name for functions, or the parameter name for block types
func (t *NameNode) SelectorNode() *IDNode {
	return t.selector
}

// Token returns with the node token
func (t *NameNode) Token() string {
	return "NAME"
}

// Schema returns the schema for the node's value
func (t *NameNode) Schema() interface{} {
	return schema.StringValue()
}

// Value returns with the block type
func (t *NameNode) Value() interface{} {
	if t.selector == nil {
		return t.name.ID()
	}

	var sep string
	switch v := t.sep.Value().(type) {
	case rune:
		sep = string(v)
	case string:
		sep = v
	default:
		panic(fmt.Errorf("unexpected name separator type: %T", v))
	}

	return ID(fmt.Sprintf("%s%s%s", t.selector.ID(), sep, t.name.ID()))
}

// Pos returns the position
func (t *NameNode) Pos() parsley.Pos {
	if t.selector != nil {
		return t.selector.Pos()
	}

	return t.name.Pos()
}

// ReaderPos returns the position of the first character immediately after this node
func (t *NameNode) ReaderPos() parsley.Pos {
	return t.name.readerPos
}

// SetReaderPos changes the reader position
func (t *NameNode) SetReaderPos(f func(parsley.Pos) parsley.Pos) {
	t.name.SetReaderPos(f)
}

// String returns with a string representation of the node
func (t *NameNode) String() string {
	return fmt.Sprintf("NAME{%s, %d..%d}", t.Value(), t.Pos(), t.ReaderPos())
}
