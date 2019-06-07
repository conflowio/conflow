// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser

import (
	"fmt"

	"github.com/opsidian/basil/basil"
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
		ID(basil.FunctionNameRegExpPattern),
		terminal.Rune('('),
		text.LeftTrim(SepByComma(p, text.WsSpaces), text.WsSpaces),
		text.LeftTrim(terminal.Rune(')'), text.WsSpaces),
	).Token("FUNC").Name("function").Bind(functionInterpreter{})
}

type functionInterpreter struct{}

func (f functionInterpreter) Eval(userCtx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	panic("Eval should not be called on a raw function node")
}

func (f functionInterpreter) TransformNode(userCtx interface{}, node parsley.Node) (parsley.Node, parsley.Error) {
	registry := userCtx.(basil.FunctionTransformerRegistryAware).FunctionTransformerRegistry()

	nodes := node.(parsley.NonTerminalNode).Children()
	nameNode := nodes[0]
	name, _ := nameNode.Value(nil)

	transformer, exists := registry.NodeTransformer(string(name.(basil.ID)))
	if !exists {
		return nil, parsley.NewError(nameNode.Pos(), fmt.Errorf("%q function does not exist", name))
	}

	return transformer.TransformNode(userCtx, node)
}
