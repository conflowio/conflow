// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsers

import (
	"fmt"

	"github.com/conflowio/parsley/combinator"
	"github.com/conflowio/parsley/parser"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/text"
	"github.com/conflowio/parsley/text/terminal"

	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/conflow/block"
)

// Directive returns a parser for parsing directives
//   S     -> "@" ID {
//              (PARAMETER|BLOCK)*
//            }
//         -> ID VALUE
//   VALUE -> EXPRESSION
//         -> ARRAY
//         -> MAP
func Directive(expr parsley.Parser) *combinator.Sequence {
	paramOrBlock := combinator.Choice(
		Parameter(expr, false, false),
		blockWithOptions(expr, false, false, false, false),
	).Name("parameter or block")

	emptyBody := combinator.SeqOf(
		terminal.Rune('{'),
		text.LeftTrim(terminal.Rune('}'), text.WsSpacesNl),
	).Token(block.TokenBlockBody)

	body := combinator.SeqOf(
		terminal.Rune('{'),
		combinator.Many(text.LeftTrim(paramOrBlock, text.WsSpacesForceNl)),
		text.LeftTrim(terminal.Rune('}'), text.WsSpacesForceNl),
	).Token(block.TokenBlockBody)

	blockValue := combinator.Choice(
		combinator.SuppressError(emptyBody),
		body,
		expr,
		Array(expr),
		Map(expr),
	).Name("block value")

	return combinator.SeqOf(
		parser.Empty(), // no directives for a directive
		IDWithClassifier(conflow.ClassifierAnnotation),
		combinator.Choice(
			text.LeftTrim(blockValue, text.WsSpaces),
			parser.Empty(),
		),
	).Name("directive").Token(block.TokenDirective).Bind(directiveInterpreter{})
}

type directiveInterpreter struct{}

func (d directiveInterpreter) Eval(userCtx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	panic("Eval should not be called on a raw block node")
}

func (d directiveInterpreter) TransformNode(userCtx interface{}, node parsley.Node) (parsley.Node, parsley.Error) {
	registry := userCtx.(conflow.DirectiveTransformerRegistryAware).DirectiveTransformerRegistry()

	nodes := node.(parsley.NonTerminalNode).Children()

	typeNode := nodes[1].(*conflow.IDNode)
	transformer, exists := registry.NodeTransformer(string(typeNode.ID()))
	if !exists {
		return nil, parsley.NewError(typeNode.Pos(), fmt.Errorf(
			"%s%s directive is unknown or not allowed",
			string(typeNode.Classifier()),
			typeNode.ID()),
		)
	}

	return transformer.TransformNode(userCtx, node)
}
