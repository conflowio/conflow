package parser

import (
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

// SepByComma applies the given value parser zero or more times separated by comma
func SepByComma(p parsley.Parser, wsMode text.WsMode) *combinator.Sequence {
	comma := text.LeftTrim(terminal.Rune(','), text.WsSpaces)
	ptrim := text.LeftTrim(p, wsMode)

	lookup := func(i int) parsley.Parser {
		if i == 0 {
			return p
		}
		if i%2 == 0 {
			return ptrim
		} else {
			return comma
		}
	}
	lenCheck := func(len int) bool {
		if wsMode == text.WsNone || wsMode == text.WsSpaces {
			return len == 0 || len%2 == 1
		} else {
			return len%2 == 0
		}
	}
	return combinator.Seq("SEP_BY", lookup, lenCheck)
}
