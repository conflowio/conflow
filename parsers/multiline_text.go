// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers

import (
	"bytes"
	"errors"

	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"

	"github.com/opsidian/conflow/basil/schema"
)

// MultilineText parses a multiline text in the following format:
//
// '''
//   first line
//   second line
// '''
//
// This will result in "first line\nsecond line\n".
//
// Also works with """ and ``` quotes, but the opening and closing quotes must be the same.
//
// There must be a new line after the opening quotes and the closing quotes must be after a new line
// All text lines must be empty or start with the same whitespace characters as the first text line
// The common prefix whitespace characters will be stripped from all lines
func MultilineText() parser.Func {
	notFoundErr := parsley.NotFoundError("multiline text")
	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		tr := ctx.Reader().(*text.Reader)
		var err parsley.Error
		readerPos, res := tr.Readf(pos, func(b []byte) ([]byte, int) {
			q := b[0]
			if q != '\'' && q != '"' && q != '`' {
				return nil, 0
			}
			if b[0] != q || b[1] != q || b[2] != q {
				return nil, 0
			}

			i := 3
			textInLine := false
			for ; i < len(b)-2; i++ {
				if b[i] == '\n' {
					break
				} else if b[i] != ' ' && b[i] != '\t' {
					textInLine = true
				}
			}
			if textInLine {
				err = parsley.NewError(pos+3, errors.New("text must start in a new line after '''"))
				return nil, 0
			}
			i++

			res := make([]byte, 0, len(b))
			var ws []byte
			wsl := 0
			firstLine := true
			linePos := i
			textPos := linePos
			for ; i < len(b)-2; i++ {
				if b[i] == '\n' {
					res = append(res, b[textPos:i+1]...)
					linePos = i + 1
					textPos = linePos
					firstLine = false
					textInLine = false
				} else if b[i] == q && b[i+1] == q && b[i+2] == q {
					if textInLine {
						err = parsley.NewError(pos+parsley.Pos(i), errors.New("the closing ''' must be on a new line"))
						return nil, 0
					}
					return res, i + 3
				} else if b[i] != ' ' && b[i] != '\t' && b[i] != '\n' {
					if !textInLine {
						if firstLine {
							ws = b[linePos:i]
							wsl = len(ws)
						} else {
							if !bytes.Equal(ws, b[linePos:linePos+wsl]) {
								err = parsley.NewError(
									pos+parsley.Pos(i),
									errors.New("every line must be empty or start with the same whitespace characters as the first line"),
								)
								return nil, 0
							}
						}
						textPos = linePos + wsl
						textInLine = true
					}
				}
			}

			return nil, 0
		})

		if res != nil {
			return terminal.NewStringNode(schema.StringValue(), string(res), pos, readerPos), data.EmptyIntSet, nil
		}

		if err == nil {
			err = parsley.NewError(pos, notFoundErr)
		}

		return nil, data.EmptyIntSet, err
	})
}
