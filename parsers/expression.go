// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers

import (
	"github.com/opsidian/parsley/ast/interpreter"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

// Expression returns with an expression parser
func Expression() *combinator.Sequence {
	var p combinator.Sequence

	function := Function(&p)
	array := Array(&p)
	variable := Variable()

	arrayIndex := combinator.Choice(
		terminal.Integer(),
		terminal.String(true),
		function,
		variable,
	)

	valueWithIndex := Element(combinator.Choice(
		function,
		array,
		variable,
	), arrayIndex)

	value := combinator.Choice(
		terminal.TimeDuration(),
		terminal.Float(),
		terminal.Integer(),
		MultilineText(),
		terminal.String(true),
		terminal.Bool("true", "false"),
		terminal.Nil("nil"),
		valueWithIndex,
		combinator.SeqOf(
			terminal.Rune('('),
			text.LeftTrim(&p, text.WsSpacesNl),
			text.LeftTrim(terminal.Rune(')'), text.WsSpacesNl),
		).Token("PARENS").Bind(interpreter.Select(1)),
	).Name("value")

	not := Not(value)
	prodMod := ProdMod(not)
	sum := Sum(prodMod)
	compare := Compare(sum)
	and := And(compare)
	or := Or(and)
	p = *TernaryIf(or)

	return &p
}
