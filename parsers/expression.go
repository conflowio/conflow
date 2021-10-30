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

	"github.com/conflowio/conflow/conflow/schema"
)

// Expression returns with an expression parser
func Expression() *combinator.Sequence {
	var p combinator.Sequence

	function := Function(&p)
	array := Array(&p)
	varp := Variable()

	arrayIndex := combinator.Choice(
		terminal.Integer(schema.IntegerValue()),
		terminal.String(schema.StringValue(), true),
		function,
		varp,
	)

	valueWithIndex := Element(combinator.Choice(
		function,
		array,
		varp,
	), arrayIndex)

	value := combinator.Choice(
		terminal.TimeDuration(schema.TimeDurationValue()),
		terminal.Float(schema.NumberValue()),
		terminal.Integer(schema.IntegerValue()),
		MultilineText(),
		terminal.String(schema.StringValue(), true),
		terminal.Bool(schema.BooleanValue(), "true", "false"),
		terminal.Nil(schema.NullValue(), "null"),
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
