// Copyright (c) 2018 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package variable

import (
	"fmt"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/parsley"
)

var _ basil.VariableNode = &Node{}

// Node is a function node
type Node struct {
	id            basil.ID
	blockIDNode   *basil.IDNode
	paramNameNode *basil.IDNode
	variableType  string
}

// NewNode creates a new variable node
func NewNode(blockIDNode, paramNameNode *basil.IDNode) *Node {
	return &Node{
		id:            basil.ID(fmt.Sprintf("%s.%s", blockIDNode.ID(), paramNameNode.ID())),
		blockIDNode:   blockIDNode,
		paramNameNode: paramNameNode,
	}
}

// ID returns with the variable identifier as "<block id>.<parameter name>"
func (n *Node) ID() basil.ID {
	return n.id
}

// BlockID returns with the block's id
func (n *Node) ParentID() basil.ID {
	return n.blockIDNode.ID()
}

// ParamName returns with the parameter's name
func (n *Node) ParamName() basil.ID {
	return n.paramNameNode.ID()
}

// Token returns with the node's token
func (n *Node) Token() string {
	return "VAR"
}

// Type returns with the node's type
func (n *Node) Type() string {
	return n.variableType
}

// StaticCheck runs static analysis on the node
func (n *Node) StaticCheck(ctx interface{}) parsley.Error {
	parseCtx := ctx.(*basil.ParseContext)

	blockNode, exists := parseCtx.BlockNode(n.blockIDNode.ID())
	if !exists {
		return parsley.NewErrorf(n.blockIDNode.Pos(), "block %q does not exist", n.blockIDNode.ID())
	}

	paramType, paramExists := blockNode.ParamType(n.paramNameNode.ID())
	if !paramExists {
		return parsley.NewErrorf(n.paramNameNode.Pos(), "parameter %q does not exist", n.paramNameNode.ID())
	}

	n.variableType = paramType

	return nil
}

// Value returns with the result of the function
func (n *Node) Value(ctx interface{}) (interface{}, parsley.Error) {
	blockContainer, ok := ctx.(basil.EvalContext).BlockContainer(n.blockIDNode.ID())
	if !ok {
		panic(parsley.NewErrorf(n.Pos(), "%q was referenced before it was evaluated", n.blockIDNode.ID()))
	}
	return blockContainer.Param(n.paramNameNode.ID()), nil
}

// Pos returns with the node's position
func (n *Node) Pos() parsley.Pos {
	return n.blockIDNode.Pos()
}

// ReaderPos returns with the reader's position
func (n *Node) ReaderPos() parsley.Pos {
	return n.paramNameNode.ReaderPos()
}

// SetReaderPos amends the reader position using the given function
func (n *Node) SetReaderPos(f func(parsley.Pos) parsley.Pos) {
	ast.SetReaderPos(n.paramNameNode, f)
}

func (n *Node) String() string {
	if n.variableType == "" {
		return fmt.Sprintf("%s{%s, %d..%d}", n.Token(), n.id, n.Pos(), n.ReaderPos())
	}
	return fmt.Sprintf("%s{<%s> %s, %d..%d}", n.Token(), n.variableType, n.id, n.Pos(), n.ReaderPos())
}
