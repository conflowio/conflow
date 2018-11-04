package parser

import (
	"github.com/opsidian/parsley/ast/interpreter"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

// Expression returns with an expression parser
func Expression() parser.Func {
	var p parser.Func

	function := Function(&p)
	array := Array(&p, text.WsSpaces)
	variable := Variable(&p)

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
		terminal.String(true),
		terminal.Bool("true", "false"),
		terminal.Nil("nil"),
		valueWithIndex,
		combinator.SeqOf(terminal.Rune('('), &p, terminal.Rune(')')).Bind(interpreter.Select(1)),
	).Name("value")

	not := Not(value)
	prodMod := ProdMod(not)
	sum := Sum(prodMod)
	compare := Compare(sum)
	and := And(compare)
	or := Or(and)
	p = TernaryIf(or)

	return p
}
