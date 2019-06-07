// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package test

import (
	"github.com/opsidian/basil/basil/variable"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text/terminal"
)

// EvalErrorParser returns with a parser which will read the "ERR" string but the result node evaluation will throw an error
func EvalErrorParser() parser.Func {
	return func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		res, cp, err := terminal.Word("ERR", "ERR", variable.TypeString).Parse(ctx, leftRecCtx, pos)
		if err != nil {
			return nil, cp, err
		}
		node := ast.NewNonTerminalNode("ERR", []parsley.Node{res}, ast.InterpreterFunc(func(ctx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
			return nil, parsley.NewErrorf(pos, "ERR")
		}))
		return node, cp, nil
	}
}
