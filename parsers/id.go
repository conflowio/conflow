// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers

import (
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"

	"github.com/opsidian/basil/basil"
)

// ID parses an identifier:
//   S -> /[a-z][a-z0-9]*(?:_[a-z0-9]+)*/
//
// An identifier can only contain lowercase letters, numbers and underscore characters.
// It must start with a letter, must end with a letter or number, and no duplicate underscores are allowed.
func ID() parser.Func {
	return id(basil.ClassifierNone)
}

func IDWithClassifier(classifier rune) parser.Func {
	return id(classifier)
}
func id(classifier rune) parser.Func {
	notFoundErr := parsley.NotFoundError("identifier")

	return func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		tr := ctx.Reader().(*text.Reader)

		idPos := pos

		if classifier != basil.ClassifierNone {
			var ok bool
			if pos, ok = tr.ReadRune(pos, classifier); !ok {
				return nil, data.EmptyIntSet, parsley.NewError(pos, parsley.NotFoundError(classifier))
			}
		}

		if readerPos, match := tr.ReadRegexp(pos, basil.IDRegExpPattern); match != nil {
			id := string(match)
			if ctx.IsKeyword(id) {
				return nil, data.EmptyIntSet, parsley.NewErrorf(pos, "%s is a reserved keyword", id)
			}
			return basil.NewIDNode(basil.ID(id), classifier, idPos, readerPos), data.EmptyIntSet, nil
		}
		return nil, data.EmptyIntSet, parsley.NewError(pos, notFoundErr)
	}
}
