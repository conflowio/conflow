package parser

import (
	"github.com/opsidian/basil/basil/block"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

// Parameter returns with a parameter parser
//   S  -> ID ("="|":=") P
//   ID -> /[a-z][a-z0-9]*(?:_[a-z0-9]+)*/
func Parameter(p parsley.Parser) *combinator.Sequence {
	parameterValue := combinator.Choice(
		Array(p, text.WsSpacesNl),
		Map(p),
		p,
	)
	return combinator.SeqOf(
		ID(),
		text.LeftTrim(combinator.Choice(terminal.Rune('='), terminal.Op(":=")), text.WsSpaces),
		text.LeftTrim(parameterValue, text.WsSpaces),
	).Token(block.TokenParameter)
}
