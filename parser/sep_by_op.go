package parser

import (
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
)

// SepByOp applies the given value parser one or more times separated by the op parser
func SepByOp(valueP parsley.Parser, opP parsley.Parser) *combinator.Recursive {
	lookup := func(i int) parsley.Parser {
		if i == 0 {
			return valueP
		}
		if i%2 == 0 {
			return text.LeftTrim(valueP, text.WsSpaces)
		} else {
			return text.LeftTrim(opP, text.WsSpaces)
		}
	}
	lenCheck := func(len int) bool {
		return len%2 == 1
	}
	return combinator.NewRecursive("SEP_BY", lookup, lenCheck)
}
