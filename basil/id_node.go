// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

import (
	"fmt"

	"github.com/opsidian/parsley/parsley"
)

// IDNode contains an identifier
type IDNode struct {
	id        ID
	pos       parsley.Pos
	readerPos parsley.Pos
}

// NewIDNode creates a new id node
func NewIDNode(id ID, pos parsley.Pos, readerPos parsley.Pos) *IDNode {
	if id == "" {
		panic("identifier can not be empty")
	}

	return &IDNode{
		id:        id,
		pos:       pos,
		readerPos: readerPos,
	}
}

// ID returns with the identifier string
func (i *IDNode) ID() ID {
	return i.id
}

// Token returns with the node token
func (i *IDNode) Token() string {
	return "ID"
}

// Type returns with the type of the node
func (i *IDNode) Type() string {
	return "basil.ID"
}

// Value returns with the id of the node
func (i *IDNode) Value(ctx interface{}) (interface{}, parsley.Error) {
	return i.id, nil
}

// Pos returns the position
func (i *IDNode) Pos() parsley.Pos {
	return i.pos
}

// ReaderPos returns the position of the first character immediately after this node
func (i *IDNode) ReaderPos() parsley.Pos {
	return i.readerPos
}

// SetReaderPos changes the reader position
func (i *IDNode) SetReaderPos(f func(parsley.Pos) parsley.Pos) {
	i.readerPos = f(i.readerPos)
}

// String returns with a string representation of the node
func (i *IDNode) String() string {
	return fmt.Sprintf("ID{%v, %d..%d}", i.id, i.pos, i.readerPos)
}
