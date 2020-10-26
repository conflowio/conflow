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
	ptrim := text.LeftTrim(p, text.WsSpacesNl)
	optrim := combinator.SuppressError(text.LeftTrim(op, text.WsSpacesNl))

	lookup := func(i int) parsley.Parser {
		switch {
		case i == 0:
			return p
		case i%2 == 0:
			return ptrim
		default:
			return optrim
		}
	}
	lenCheck := func(len int) bool {
		return len%2 == 1
	}
	return combinator.Seq("SEP_BY_OP", lookup, lenCheck)
}
