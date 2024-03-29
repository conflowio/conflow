// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conflow

import (
	"fmt"
	"strings"

	"github.com/conflowio/parsley/parsley"

	"github.com/conflowio/conflow/pkg/internal/utils"
	"github.com/conflowio/conflow/pkg/schema"
)

const (
	ClassifierNone       = rune(0)
	ClassifierAnnotation = '@'
)

// IDNode contains an identifier
type IDNode struct {
	id         ID
	classifier rune
	pos        parsley.Pos
	readerPos  parsley.Pos
}

// NewIDNode creates a new id node
func NewIDNode(id ID, classifier rune, pos parsley.Pos, readerPos parsley.Pos) *IDNode {
	if id == "" {
		panic("identifier can not be empty")
	}

	return &IDNode{
		id:         id,
		classifier: classifier,
		pos:        pos,
		readerPos:  readerPos,
	}
}

// ID returns with the identifier string
func (i *IDNode) ID() ID {
	return i.id
}

// Classifier returns with classifier of the id
// A classifier is a single rune, like '@', '$', '&'.
func (i *IDNode) Classifier() rune {
	return i.classifier
}

// Token returns with the node token
func (i *IDNode) Token() string {
	return "ID"
}

// Schema returns the schema for the node's value
func (i *IDNode) Schema() interface{} {
	return schema.StringValue()
}

// ValueOf returns with the id of the node
func (i *IDNode) Value() interface{} {
	return i.id
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

func (i *IDNode) StaticCheck(userCtx interface{}) parsley.Error {
	if !schema.NameRegExp.MatchString(string(i.id)) {
		// Let's suggest a valid identifier, if we can
		expected := utils.ToSnakeCase(string(i.id))
		expected = strings.ReplaceAll(expected, "__", "_")
		expected = strings.Trim(expected, "_")
		if expected != "" && schema.NameRegExp.MatchString(expected) {
			return parsley.NewErrorf(i.pos, "invalid identifier %q (did you mean %q?)", string(i.id), expected)
		} else {
			return parsley.NewError(i.pos, fmt.Errorf("invalid identifier %q", string(i.id)))
		}
	}

	return nil
}
