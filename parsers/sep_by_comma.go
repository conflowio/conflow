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
func SepByComma(p parsley.Parser, wsMode text.WsMode) *combinator.Sequence {
	comma := text.LeftTrim(terminal.Rune(','), text.WsSpaces)
	ptrim := text.LeftTrim(p, wsMode)

	lookup := func(i int) parsley.Parser {
		if i == 0 && (wsMode == text.WsNone || wsMode == text.WsSpaces) {
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
