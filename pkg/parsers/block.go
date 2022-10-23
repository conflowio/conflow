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

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/conflow/block"
	"github.com/conflowio/conflow/pkg/schema"
)

// Block returns a parser for parsing blocks
//
//	S     -> ID? TYPE KEY? {
//	           (ATTR|S)*
//	         }
//	      -> ID? TYPE KEY? VALUE
//	ID    -> /[a-z][a-z0-9]*(?:_[a-z0-9]+)*/
//	TYPE  -> /[a-z][a-z0-9]*(?:_[a-z0-9]+)*/
//	KEY   -> STRING LITERAL
//	ATTR  -> ID ("="|":=") P
//	VALUE -> EXPRESSION
//	      -> ARRAY
//	      -> MAP
func Block(expr parsley.Parser) *combinator.Sequence {
	return blockWithOptions(expr, true, true, true, true)
}

func blockWithOptions(
	expr parsley.Parser,
	allowID bool,
	allowCustomParameters bool,
	allowDirectives bool,
	allowKey bool,
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

	var idParsers []parsley.Parser
	if allowID {
		idParsers = append(idParsers, combinator.Optional(ID()))
	}
	idParsers = append(idParsers, text.LeftTrim(Name(':'), text.WsSpaces))
	if allowKey {
		idParsers = append(idParsers, combinator.Optional(text.LeftTrim(terminal.String(schema.StringValue(), false), text.WsSpaces)))
	}
	var idName parsley.Parser
	if len(idParsers) == 1 {
		idName = idParsers[0]
	} else {
		idName = combinator.SeqOf(idParsers...)
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
	registry := userCtx.(conflow.BlockTransformerRegistryAware).BlockTransformerRegistry()

	nodes := node.(parsley.NonTerminalNode).Children()

	var nameNode *conflow.NameNode
	switch n := nodes[1].(type) {
	case parsley.NonTerminalNode:
		nameNode = n.Children()[1].(*conflow.NameNode)
	case *conflow.NameNode:
		nameNode = n
	case *conflow.IDNode:
		nameNode = conflow.NewNameNode(nil, nil, n)
	default:
		panic(fmt.Errorf("unexpected node type: %T", nodes[1]))
	}

	if err := nameNode.StaticCheck(userCtx); err != nil {
		return nil, err
	}

	transformer, exists := registry.NodeTransformer(string(nameNode.NameNode().ID()))
	if !exists {
		return nil, parsley.NewError(nameNode.Pos(), fmt.Errorf("%q block is unknown or not allowed", nameNode.NameNode().ID()))
	}

	return transformer.TransformNode(userCtx, node)
}
