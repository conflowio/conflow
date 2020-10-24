// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers

import (
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

// SepByComma applies the given value parser zero or more times separated by comma
func SepByComma(p parsley.Parser) *combinator.Sequence {
	comma := text.LeftTrim(terminal.Rune(','), text.WsSpaces)
	ptrim := text.LeftTrim(p, text.WsSpacesNl)

	lookup := func(i int) parsley.Parser {
		switch {
		case i == 0:
			return p
		case i%2 == 0:
			return ptrim
		default:
			return comma
		}
	}
	lenCheck := func(len int) bool {
		return true
	}
	return combinator.Seq("SEP_BY_COMMA", lookup, lenCheck)
}
