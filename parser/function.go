// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser

import (
	"errors"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

// Function will match a function call defined by the following rule, where P is the input parser:
//   S      -> ID "(" PARAMS ")"
//   ID     -> /[a-z][a-z0-9]*(?:_[a-z0-9]+)*/
//   PARAMS -> EMPTY
//          -> P ("," P)*
func Function(p parsley.Parser) *combinator.Sequence {
	return combinator.SeqOf(
		ID(),
		terminal.Rune('('),
		text.LeftTrim(SepByComma(p, text.WsSpaces), text.WsSpaces),
		text.LeftTrim(terminal.Rune(')'), text.WsSpaces),
	).Bind(ast.InterpreterFunc(evalFunction))
}

func evalFunction(ctx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	nodes := node.Children()
	registry := ctx.(basil.FunctionRegistryAware).GetFunctionRegistry()

	functioNode := nodes[0]
	name, _ := functioNode.Value(ctx)

	if !registry.FunctionExists(name.(basil.ID)) {
		return nil, parsley.NewError(functioNode.Pos(), errors.New("function does not exist"))
	}

	paramsNode := nodes[2].(parsley.NonTerminalNode)
	var params []parsley.Node
	children := paramsNode.Children()
	childrenCount := len(children)
	if childrenCount > 0 {
		params = make([]parsley.Node, childrenCount/2+1)
		for i := 0; i < childrenCount; i += 2 {
			params[i/2] = children[i]
		}
	}
	return registry.CallFunction(ctx, nodes[0], params)
}
