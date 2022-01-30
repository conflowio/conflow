// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package test

import (
	"github.com/conflowio/parsley/ast"
	"github.com/conflowio/parsley/data"
	"github.com/conflowio/parsley/parser"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/text/terminal"

	"github.com/conflowio/conflow/src/schema"
)

// EvalErrorParser returns with a parser which will read the "ERR" string but the result node evaluation will throw an error
func EvalErrorParser(s schema.Schema, word string) parser.Func {
	return func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		res, cp, err := terminal.Word(s, word, word).Parse(ctx, leftRecCtx, pos)
		if err != nil {
			return nil, cp, err
		}
		node := ast.NewNonTerminalNode("ERR", []parsley.Node{res}, errInterpreter{schema: s})
		return node, cp, nil
	}
}

type errInterpreter struct {
	schema schema.Schema
}

func (e errInterpreter) StaticCheck(userCtx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	return e.schema, nil
}

func (e errInterpreter) Eval(userCtx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	return nil, parsley.NewErrorf(node.Pos(), "ERR")
}
