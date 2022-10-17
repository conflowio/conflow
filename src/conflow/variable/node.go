// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package variable

import (
	"fmt"

	"github.com/conflowio/parsley/ast"
	"github.com/conflowio/parsley/parsley"

	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/schema"
)

var _ conflow.VariableNode = &Node{}

// Node is a function node
type Node struct {
	id            conflow.ID
	blockIDNode   *conflow.IDNode
	paramNameNode *conflow.IDNode
	schema        schema.Schema
}

// NewNode creates a new variable node
func NewNode(blockIDNode, paramNameNode *conflow.IDNode) *Node {
	return &Node{
		id:            conflow.ID(fmt.Sprintf("%s.%s", blockIDNode.ID(), paramNameNode.ID())),
		blockIDNode:   blockIDNode,
		paramNameNode: paramNameNode,
	}
}

// ID returns with the variable identifier as "<block id>.<parameter name>"
func (n *Node) ID() conflow.ID {
	return n.id
}

// BlockID returns with the block's id
func (n *Node) ParentID() conflow.ID {
	return n.blockIDNode.ID()
}

// ParamName returns with the parameter's name
func (n *Node) ParamName() conflow.ID {
	return n.paramNameNode.ID()
}

// Token returns with the node's token
func (n *Node) Token() string {
	return "VAR"
}

// Schema returns the schema for the node's value
func (n *Node) Schema() interface{} {
	return n.schema
}

// StaticCheck runs static analysis on the node
func (n *Node) StaticCheck(ctx interface{}) parsley.Error {
	parseCtx := ctx.(*conflow.ParseContext)

	if err := n.blockIDNode.StaticCheck(ctx); err != nil {
		return err
	}

	if err := n.paramNameNode.StaticCheck(ctx); err != nil {
		return err
	}

	blockNode, exists := parseCtx.BlockNode(n.blockIDNode.ID())
	if !exists {
		return parsley.NewErrorf(n.blockIDNode.Pos(), "block %q does not exist", n.blockIDNode.ID())
	}

	paramSchema, paramExists := blockNode.GetPropertySchema(n.paramNameNode.ID())
	if !paramExists {
		return parsley.NewErrorf(n.paramNameNode.Pos(), "parameter %q does not exist", n.paramNameNode.ID())
	}

	n.schema = paramSchema

	return nil
}

// Value returns with the result of the function
func (n *Node) Value(ctx interface{}) (interface{}, parsley.Error) {
	blockContainer, ok := ctx.(*conflow.EvalContext).BlockContainer(n.blockIDNode.ID())
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
	return fmt.Sprintf("%s{%s, %d..%d}", n.Token(), n.id, n.Pos(), n.ReaderPos())
}
