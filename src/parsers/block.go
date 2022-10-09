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

// Block returns a parser for parsing blocks
//   S     -> ID? TYPE? {
//              (ATTR|S)*
//            }
//         -> ID? TYPE VALUE
//   ID    -> /[a-z][a-z0-9]*(?:_[a-z0-9]+)*/
//   ATTR  -> ID ("="|":=") P
//   VALUE -> EXPRESSION
//         -> ARRAY
//         -> MAP
func Block(expr parsley.Parser) *combinator.Sequence {
	return blockWithOptions(expr, true, true, true)
}

func blockWithOptions(
	expr parsley.Parser,
	allowID bool,
	allowCustomParameters bool,
	allowDirectives bool,
) *combinator.Sequence {
	var p combinator.Sequence

	var directives parsley.Parser
	if allowDirectives {
		directives = combinator.Many(text.RightTrim(Directive(expr), text.WsSpacesForceNl))
	} else {
		directives = parser.Empty()
	}

	paramOrBlock := combinator.Choice(
		Parameter(expr, allowCustomParameters, allowDirectives),
		&p,
	).Name("parameter or block")

	emptyBody := combinator.SeqOf(
		terminal.Rune('{'),
		text.LeftTrim(terminal.Rune('}'), text.WsSpacesNl),
	).Name("block body").Token(block.TokenBlockBody)

	body := combinator.SeqOf(
		terminal.Rune('{'),
		combinator.Many(text.LeftTrim(paramOrBlock, text.WsSpacesForceNl)),
		text.LeftTrim(terminal.Rune('}'), text.WsSpacesForceNl),
	).Name("block body").Token(block.TokenBlockBody)

	blockValue := combinator.Choice(
		combinator.SuppressError(emptyBody),
		body,
		expr,
		Array(expr),
		Map(expr),
	).Name("block value")

	var idName parsley.Parser
	if allowID {
		idName = combinator.SeqOf(
			combinator.Optional(ID()),
			text.LeftTrim(Name(':'), text.WsSpaces),
		)
	} else {
		idName = Name(':')
	}

	p = *combinator.SeqOf(
		directives,
		idName,
		combinator.Optional(text.LeftTrim(blockValue, text.WsSpaces)),
	).Name("block").Token(block.TokenBlock).Bind(blockInterpreter{})

	return &p
}

type blockInterpreter struct{}

func (b blockInterpreter) Eval(userCtx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	panic("Eval should not be called on a raw block node")
}

func (b blockInterpreter) TransformNode(userCtx interface{}, node parsley.Node) (parsley.Node, parsley.Error) {
	nodes := node.(parsley.NonTerminalNode).Children()
	_, nameNode := block.GetIDAndNameFromNode(nodes[1])

	registry := userCtx.(conflow.BlockTransformerRegistryAware).BlockTransformerRegistry()
	transformer, exists := registry.NodeTransformer(string(nameNode.NameNode().ID()))
	if !exists {
		return nil, parsley.NewError(nameNode.Pos(), fmt.Errorf("%q block is unknown or not allowed", nameNode.NameNode().ID()))
	}

	return transformer.TransformNode(userCtx, node)
}
