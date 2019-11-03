// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parameter

import (
	"fmt"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/parsley/parsley"
)

var _ basil.ParameterNode = &Node{}

const (
	Token = "PARAMETER"
)

// Node is a block parameter
type Node struct {
	id            basil.ID
	nameNode      *basil.IDNode
	valueNode     parsley.Node
	evalStage     basil.EvalStage
	isDeclaration bool
	dependencies  basil.Dependencies
}

// NewNode creates a new block parameter node
func NewNode(blockID basil.ID, nameNode *basil.IDNode, valueNode parsley.Node, isDeclaration bool) *Node {
	return &Node{
		id:            basil.ID(fmt.Sprintf("%s.%s", blockID, nameNode.ID())),
		nameNode:      nameNode,
		valueNode:     valueNode,
		isDeclaration: isDeclaration,
	}
}

// ID returns with the name of the parameter
func (n *Node) ID() basil.ID {
	return n.id
}

// ID returns with the name of the parameter
func (n *Node) Name() basil.ID {
	return n.nameNode.ID()
}

// ValueNode returns with the value node
func (n *Node) ValueNode() parsley.Node {
	return n.valueNode
}

// Token returns with the node token
func (n *Node) Token() string {
	return "BLOCK_PARAM"
}

// Type returns with the value node's type
func (n *Node) Type() string {
	return n.valueNode.Type()
}

// EvalStage returns with the evaluation stage
func (n *Node) EvalStage() basil.EvalStage {
	if n.evalStage == basil.EvalStageUndefined {
		return basil.EvalStageMain
	}
	return n.evalStage
}

// Dependencies returns the blocks/parameters this parameter depends on
func (n *Node) Dependencies() basil.Dependencies {
	if n.dependencies != nil {
		return n.dependencies
	}

	n.dependencies = make(basil.Dependencies)

	parsley.Walk(n.valueNode, func(node parsley.Node) bool {
		if v, ok := node.(basil.VariableNode); ok {
			n.dependencies[v.ID()] = v
		}
		return false
	})

	return n.dependencies
}

// Directives returns nil as currently parameters don't support directives
func (n *Node) Directives() []basil.BlockNode {
	return nil
}

// Provides returns with nil as a parameter node doesn't define other nodes
func (n *Node) Provides() []basil.ID {
	return nil
}

// Generates returns with nil as a parameter node doesn't generate other nodes
func (n *Node) Generates() []basil.ID {
	return nil
}

// IsDeclaration returns true if the parameter was declared in the block
func (n *Node) IsDeclaration() bool {
	return n.isDeclaration
}

// SetDescriptor applies the descriptor parameters to the node
func (n *Node) SetDescriptor(descriptor basil.ParameterDescriptor) {
	if descriptor.EvalStage != basil.EvalStageUndefined {
		n.evalStage = descriptor.EvalStage
	}
}

// Generated returns true if the parameter's value contains a generator function
func (n *Node) Generated() bool {
	// TODO: implement generator functions
	return false
}

// StaticCheck runs a static analysis on the value node
func (n *Node) StaticCheck(ctx interface{}) parsley.Error {
	switch n := n.valueNode.(type) {
	case parsley.StaticCheckable:
		if err := n.StaticCheck(ctx); err != nil {
			return err
		}
	}

	return nil
}

// Value returns with the value of the node
func (n *Node) Value(ctx interface{}) (interface{}, parsley.Error) {
	return n.valueNode.Value(ctx)
}

// Pos returns the position
func (n *Node) Pos() parsley.Pos {
	return n.nameNode.Pos()
}

// ReaderPos returns the position of the first character immediately after this node
func (n *Node) ReaderPos() parsley.Pos {
	return n.valueNode.ReaderPos()
}

// Walk runs the given function on all child nodes
func (n *Node) Walk(f func(n parsley.Node) bool) bool {
	return parsley.Walk(n.valueNode, f)
}

func (n *Node) CreateContainer(
	ctx *basil.EvalContext,
	parent basil.BlockContainer,
	value interface{},
	wgs []basil.WaitGroup,
	pending bool,
) basil.JobContainer {
	return NewContainer(ctx, n, parent, value, wgs, pending)
}

// String returns with a string representation of the node
func (n *Node) String() string {
	return fmt.Sprintf("%s{%s=%s, %d..%d}", n.Token(), n.nameNode, n.valueNode, n.Pos(), n.ReaderPos())
}
