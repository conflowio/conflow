// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/parameter"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

// Parameter returns with a parameter parser
// If allowNewAssignment is false then only "=" will be allowed
//   S  -> ID ("="|":=") P
//   ID -> /[a-z][a-z0-9]*(?:_[a-z0-9]+)*/
func Parameter(expr parsley.Parser, allowNewAssignment bool, allowDirectives bool) *combinator.Sequence {
	var directives parsley.Parser
	if allowDirectives {
		directives = combinator.Many(text.RightTrim(Directive(expr), text.WsSpacesForceNl))
	} else {
		directives = parser.Empty()
	}

	var parameterValue parser.FuncWrapper
	parameterValue = parser.FuncWrapper{F: combinator.Choice(
		Array(&parameterValue, text.WsSpacesNl),
		Map(&parameterValue),
		expr,
	)}

	var assignment parsley.Parser
	if allowNewAssignment {
		assignment = combinator.Choice(terminal.Rune('='), terminal.Op(":="))
	} else {
		assignment = terminal.Rune('=')
	}

	return combinator.SeqOf(
		directives,
		ID(basil.IDRegExpPattern),
		text.LeftTrim(assignment, text.WsSpaces),
		text.LeftTrim(parameterValue, text.WsSpaces),
	).Token(parameter.Token)
}
