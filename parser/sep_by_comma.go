package parser

import (
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

// SepByComma applies the given value parser zero or more times separated by comma
func SepByComma(valueP parsley.Parser, wsMode text.WsMode) *combinator.Sequence {
	commaP := text.LeftTrim(terminal.Rune(','), wsMode)

	lookup := func(i int) parsley.Parser {
		if i == 0 {
			return valueP
		}
		if i%2 == 0 {
			return text.LeftTrim(valueP, wsMode)
		} else {
			return commaP
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
