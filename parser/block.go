// Copyright (c) 2018 Opsidian Ltd.
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

// Block return a parser for parsing blocks
//   S     -> ID ID? {
//              (ATTR|S)*
//            }
//         -> ID ID? VALUE
//   ID    -> /[a-z][a-z0-9]*(?:_[a-z0-9]+)*/
//   ATTR  -> ID ("="|":=") P
//   VALUE -> STRING
//         -> INT
//         -> FLOAT
//         -> BOOL
//         -> TIME_DURATION
func Block() *combinator.Sequence {
	var p combinator.Sequence
	expr := Expression()

	parameterValue := combinator.Choice(
		Array(expr, text.WsSpacesNl),
		Map(expr),
		expr,
	)
	parameter := combinator.SeqOf(
		ID(),
		text.LeftTrim(combinator.Choice(terminal.Rune('='), terminal.Op(":=")), text.WsSpaces),
		text.LeftTrim(parameterValue, text.WsSpaces),
	).Token("PARAMETER")

	paramOrBlock := combinator.Choice(
		parameter,
		&p,
	).Name("parameter or block definition")

	emptyBlockValue := combinator.SeqOf(
		terminal.Rune('{'),
		text.LeftTrim(terminal.Rune('}'), text.WsSpacesNl),
	).Token("BLOCK_BODY")

	nonEmptyBlockValue := combinator.SeqOf(
		terminal.Rune('{'),
		combinator.Many(
			text.LeftTrim(paramOrBlock, text.WsSpacesForceNl),
		),
		text.LeftTrim(terminal.Rune('}'), text.WsSpacesForceNl),
	).Token("BLOCK_BODY")

	blockValue := combinator.Choice(
		emptyBlockValue,
		nonEmptyBlockValue,
		terminal.TimeDuration(),
		terminal.Float(),
		terminal.Integer(),
		terminal.String(true),
		terminal.Bool("true", "false"),
		Array(expr, text.WsSpaces),
		Array(expr, text.WsSpacesNl),
		Map(expr),
	).Name("block value")

	p = *combinator.SeqTry(
		combinator.SeqTry(ID(), text.LeftTrim(ID(), text.WsSpaces)),
		text.LeftTrim(blockValue, text.WsSpaces),
	).Name("block definition").Token("BLOCK").Bind(blockInterpreter{})

	return &p
}

type blockInterpreter struct{}

func (b blockInterpreter) Eval(userCtx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	panic("Eval should not be called on a raw block node")
}

func (b blockInterpreter) TransformNode(userCtx interface{}, node parsley.Node) (parsley.Node, parsley.Error) {
	registry := userCtx.(basil.BlockTransformerRegistryAware).BlockTransformerRegistry()

	nodes := node.(parsley.NonTerminalNode).Children()
	blockIDNodes := nodes[0].(parsley.NonTerminalNode).Children()
	typeNode := blockIDNodes[0]
	blockType, _ := typeNode.Value(nil)

	transformer, exists := registry.NodeTransformer(string(blockType.(basil.ID)))
	if !exists {
		return nil, parsley.NewError(typeNode.Pos(), fmt.Errorf("%q type is invalid or not allowed here", blockType))
	}

	return transformer.TransformNode(userCtx, node)
}
