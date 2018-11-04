// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/ast/interpreter"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

// Map will match an map expression defined by the following rule, where P is the input parser:
//   S -> "map" "{" "}"
//   S -> "map" "{"
//           (STRING ":" P ",")*
//        "}"
func Map(p parsley.Parser) parser.Func {
	keyValue := combinator.SeqOf(
		terminal.String(false),
		text.LeftTrim(terminal.Rune(':'), text.WsSpaces),
		text.LeftTrim(p, text.WsSpaces),
	).Name("key-value pair")

	emptyMap := combinator.SeqOf(
		terminal.Word("map", "map", basil.TypeString),
		text.LeftTrim(terminal.Rune('{'), text.WsSpaces),
		text.LeftTrim(terminal.Rune('}'), text.WsSpacesNl),
	).Bind(ast.InterpreterFunc(evalEmptyMap))

	nonEmptyMap := combinator.SeqOf(
		terminal.Word("map", "map", basil.TypeString),
		text.LeftTrim(terminal.Rune('{'), text.WsSpaces),
		SepByComma(keyValue, text.WsSpacesForceNl).Bind(interpreter.Object()),
		text.LeftTrim(terminal.Rune('}'), text.WsSpacesForceNl),
	).Bind(interpreter.Select(2))

	return combinator.Choice(
		emptyMap,
		nonEmptyMap,
	).Name("map")
}

func evalEmptyMap(ctx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	return map[string]interface{}{}, nil
}
