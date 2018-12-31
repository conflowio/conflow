// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/variable"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text/terminal"
)

// Variable will match a variable expression defined by the following rule, where P is the input parser:
//   S         -> ID ("." ID)?
//   ID        -> /[a-z][a-z0-9]*(?:_[a-z0-9]+)*/
//
// Variable can refer to named block's field or to a field in the root block. In the first case it will be in <block ID>.<field ID> format, otherwise <field ID> (which really a short format for root.<field ID>).
func Variable() *combinator.Sequence {
	return combinator.SeqFirstOrAll(
		ID(),
		terminal.Rune('.'),
		ID(),
	).Token("VAR").Bind(variableInterpreter{})
}

type variableInterpreter struct{}

func (v variableInterpreter) StaticCheck(userCtx interface{}, node parsley.NonTerminalNode) (string, parsley.Error) {
	blockNodeRegistry := userCtx.(basil.BlockNodeRegistryAware).BlockNodeRegistry()

	nodes := node.Children()
	var blockName basil.ID
	if len(nodes) == 1 {
		blockName = basil.ID("root")
	} else {
		blockName, _ = variable.NodeIdentifierValue(nodes[0], userCtx)
	}
	paramName, _ := variable.NodeIdentifierValue(nodes[len(nodes)-1], userCtx)

	blockNode, exists := blockNodeRegistry.BlockNode(blockName)
	if !exists {
		return "", parsley.NewErrorf(node.Pos(), "block %q does not exist", blockName)
	}

	paramType, paramExists := blockNode.ParamType(paramName)
	if !paramExists {
		return "", parsley.NewErrorf(nodes[len(nodes)-1].Pos(), "parameter %q does not exist", paramName)
	}

	return paramType, nil
}

func (v variableInterpreter) Eval(userCtx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	blockContainerRegistry := userCtx.(basil.BlockContainerRegistryAware).BlockContainerRegistry()

	nodes := node.Children()
	var blockName basil.ID
	if len(nodes) == 1 {
		blockName = basil.ID("root")
	} else {
		blockName, _ = variable.NodeIdentifierValue(nodes[0], userCtx)
	}
	paramName, _ := variable.NodeIdentifierValue(nodes[len(nodes)-1], userCtx)

	blockContainer, ok := blockContainerRegistry.BlockContainer(blockName)
	if !ok {
		return nil, parsley.NewErrorf(node.Pos(), "block %q does not exist", blockName)
	}

	return blockContainer.Param(paramName), nil
}
