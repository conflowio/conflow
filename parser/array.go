// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser

import (
	"github.com/opsidian/parsley/ast/interpreter"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

// Array will match an array expression defined by the following rule, where P is the input parser:
//   S -> "[" "]"
//   S -> "[" P ("," P)* "]"
func Array(p parsley.Parser, wsMode text.WsMode) *combinator.Recursive {
	return combinator.Seq(
		terminal.Rune('['),
		text.LeftTrim(SepByComma(p, wsMode).Bind(interpreter.Array()), wsMode),
		text.LeftTrim(terminal.Rune(']'), wsMode),
	).Bind(interpreter.Select(1))
}
