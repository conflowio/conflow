// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers

import (
	"fmt"

	"github.com/conflowio/parsley/combinator"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/text"
	"github.com/conflowio/parsley/text/terminal"

	"github.com/conflowio/conflow/pkg/conflow"
)

// Function will match a function call defined by the following rule, where P is the input parser:
//
//	S      -> ID "(" PARAMS ")"
//	ID     -> /[a-z][a-z0-9]*(?:_[a-z0-9]+)*/
//	PARAMS -> EMPTY
//	       -> P ("," P)*
func Function(p parsley.Parser) *combinator.Sequence {
	return combinator.SeqOf(
		Name('.'),
		terminal.Rune('('),
		text.LeftTrim(SepByComma(p), text.WsSpacesNl),
		text.LeftTrim(terminal.Rune(')'), text.WsSpaces),
	).Name("function").Token("FUNC").Bind(functionInterpreter{})
}

type functionInterpreter struct{}

func (f functionInterpreter) Eval(userCtx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	panic("Eval should not be called on a raw function node")
}

func (f functionInterpreter) TransformNode(userCtx interface{}, node parsley.Node) (parsley.Node, parsley.Error) {
	registry := userCtx.(conflow.FunctionTransformerRegistryAware).FunctionTransformerRegistry()

	nodes := node.(parsley.NonTerminalNode).Children()
	nameNode := nodes[0].(*conflow.NameNode)
	if err := nameNode.StaticCheck(userCtx); err != nil {
		return nil, err
	}
	name := nameNode.Value().(conflow.ID)

	transformer, exists := registry.NodeTransformer(string(name))
	if !exists {
		return nil, parsley.NewError(nameNode.Pos(), fmt.Errorf("%q function does not exist", string(name)))
	}

	return transformer.TransformNode(userCtx, node)
}
