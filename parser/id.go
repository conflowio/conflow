package parser

import (
	"errors"

	"github.com/opsidian/ocl/identifier"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
)

// ID parses an identifier:
//   S -> /[a-z][a-z0-9]*(?:_[a-z0-9]+)*/
//
// An ID can only contain lowercase letters, numbers and underscore characters.
// It must start with a letter and no duplicate underscores are allowed.
//
func ID() parser.Func {
	notFoundErr := errors.New("was expecting identifier")

	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		tr := ctx.Reader().(*text.Reader)
		if readerPos, match := tr.ReadRegexp(pos, identifier.RegExpPattern); match != nil {
			id := string(match)
			if ctx.IsKeyword(id) {
				return nil, data.EmptyIntSet, parsley.NewErrorf(pos, "%s is a reserved keyword", id)
			}
			return identifier.NewNode(id, false, pos, readerPos), data.EmptyIntSet, nil
		}
		return nil, data.EmptyIntSet, parsley.NewError(pos, notFoundErr)
	})
}
