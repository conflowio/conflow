// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers

import (
	"errors"

	"github.com/conflowio/parsley/data"
	"github.com/conflowio/parsley/parser"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/text"

	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/schema"
	"github.com/conflowio/conflow/internal/utils"
)

// ID parses an identifier:
//   S -> /[a-z][a-z0-9]*(?:_[a-z0-9]+)*/
//
// An identifier can only contain lowercase letters, numbers and underscore characters.
// It must start with a letter, must end with a letter or number, and no duplicate underscores are allowed.
func ID() parser.Func {
	return id(conflow.ClassifierNone)
}

func IDWithClassifier(classifier rune) parser.Func {
	return id(classifier)
}
func id(classifier rune) parser.Func {
	notFoundErr := parsley.NotFoundError("identifier")

	return func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		tr := ctx.Reader().(*text.Reader)

		idPos := pos

		if classifier != conflow.ClassifierNone {
			var ok bool
			if pos, ok = tr.ReadRune(pos, classifier); !ok {
				return nil, data.EmptyIntSet, parsley.NewError(pos, parsley.NotFoundError(classifier))
			}
		}

		if readerPos, match := tr.ReadRegexp(pos, "[a-zA-Z_][a-zA-Z0-9_]*"); match != nil {
			id := string(match)

			if !schema.NameRegExp.MatchString(id) {
				// Let's suggest a valid identifier, if we can
				expected := utils.ToSnakeCase(id)
				if schema.NameRegExp.MatchString(expected) {
					return nil, data.EmptyIntSet, parsley.NewErrorf(pos, "invalid identifier (did you mean %q?)", expected)
				} else {
					return nil, data.EmptyIntSet, parsley.NewError(pos, errors.New("invalid identifier"))
				}
			}

			if ctx.IsKeyword(id) {
				return nil, data.EmptyIntSet, parsley.NewErrorf(pos, "%s is a reserved keyword", id)
			}
			return conflow.NewIDNode(conflow.ID(id), classifier, idPos, readerPos), data.EmptyIntSet, nil
		}
		return nil, data.EmptyIntSet, parsley.NewError(pos, notFoundErr)
	}
}
