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

	value := combinator.Choice(
		terminal.Float(),
		terminal.Integer(),
		terminal.String(true),
		terminal.Bool("true", "false"),
		terminal.Word("nil", nil),
		Function(&p),
		Array(&p, text.WsSpaces),
		Variable(&p),
		combinator.Seq(terminal.Rune('('), &p, terminal.Rune(')')).Bind(interpreter.Select(1)),
	).ReturnError("was expecting value")

	valueWithIndex := Element(value)
	not := Not(valueWithIndex)
	prodMod := ProdMod(not)
	sum := Sum(prodMod)
	compare := Compare(sum)
	and := And(compare)
	or := Or(and)
	p = TernaryIf(or)

	return p
}
