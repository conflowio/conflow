// Copyright (c) 2018 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser

import (
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

// Block return a parser for parsing blocks
//   S -> ID ID? {
//          (ATTR|S)*
//        }
//     -> ID ID? VALUE
//   ID -> /[a-z][a-z0-9]*(?:_[a-z0-9]+)*/
//   ATTR -> ID "=" P
//   VALUE -> STRING
//         -> INT
//         -> FLOAT
//         -> BOOL
//         -> TIME_DURATION
func Block() *combinator.Sequence {
	var p combinator.Sequence
	expr := Expression()

	parameterValue := combinator.Choice(
		Array(expr, text.WsSpacesNl),
		Map(expr),
		expr,
	)
	parameter := combinator.SeqOf(
		ID(),
		text.LeftTrim(terminal.Rune('='), text.WsSpaces),
		text.LeftTrim(parameterValue, text.WsSpaces),
	).Token("PARAMETER")

	paramOrBlock := combinator.Choice(
		parameter,
		&p,
	).Name("parameter or block definition")

	emptyBlockValue := combinator.SeqOf(
		terminal.Rune('{'),
		text.LeftTrim(terminal.Rune('}'), text.WsSpacesNl),
	).Token("BLOCK_BODY")

	nonEmptyBlockValue := combinator.SeqOf(
		terminal.Rune('{'),
		combinator.Many(
			text.LeftTrim(paramOrBlock, text.WsSpacesForceNl),
		),
		text.LeftTrim(terminal.Rune('}'), text.WsSpacesForceNl),
	).Token("BLOCK_BODY")

	blockValue := combinator.Choice(
		emptyBlockValue,
		nonEmptyBlockValue,
		terminal.TimeDuration(),
		terminal.Float(),
		terminal.Integer(),
		terminal.String(true),
		terminal.Bool("true", "false"),
		Array(expr, text.WsSpaces),
		Array(expr, text.WsSpacesNl),
		Map(expr),
	).Name("block value")

	p = *combinator.SeqTry(
		combinator.SeqTry(ID(), text.LeftTrim(ID(), text.WsSpaces)),
		text.LeftTrim(blockValue, text.WsSpaces),
	).Name("block definition").Token("BLOCK")

	return &p
}
