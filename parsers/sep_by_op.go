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
)

// SepByOp applies the given value parser one or more times separated by the op parser
func SepByOp(p parsley.Parser, op parsley.Parser) *combinator.Sequence {
	ptrim := text.LeftTrim(p, text.WsSpaces)
	optrim := text.LeftTrim(op, text.WsSpaces)

	lookup := func(i int) parsley.Parser {
		if i == 0 {
			return p
		}
		if i%2 == 0 {
			return ptrim
		} else {
			return optrim
		}
	}
	lenCheck := func(len int) bool {
		return len%2 == 1
	}
	return combinator.Seq("SEP_BY", lookup, lenCheck)
}
