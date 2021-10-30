// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers

import (
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text/terminal"

	"github.com/opsidian/conflow/basil/block"
	"github.com/opsidian/conflow/conflow"
)

// Name parses a name expression:
//   S  -> ID
//         ID SEP ID
//   ID -> /[a-z][a-z0-9]*(?:_[a-z0-9]+)*/
func Name(sep rune) *combinator.Sequence {
	lookup := func(i int) parsley.Parser {
		switch i {
		case 0, 2:
			return ID()
		case 1:
			return terminal.Rune(sep)
		default:
			return nil
		}
	}
	lenCheck := func(len int) bool {
		return len == 1 || len == 3
	}
	return combinator.Seq(block.TokenName, lookup, lenCheck).
		HandleResult(combinator.SeqResultHandlerFunc(func(_ parsley.Pos, _ string, nodes []parsley.Node, _ parsley.Interpreter) parsley.Node {
			if len(nodes) == 1 {
				return conflow.NewNameNode(nil, nil, nodes[0].(*conflow.IDNode))
			}
			return conflow.NewNameNode(nodes[0].(*conflow.IDNode), nodes[1].(parsley.LiteralNode), nodes[2].(*conflow.IDNode))
		}))
}
