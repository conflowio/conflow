// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parameter

import (
	"fmt"

	"github.com/conflowio/parsley/parsley"

	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/schema"
)

var _ conflow.ParameterNode = &Node{}

const (
	Token = "PARAMETER"
)

// Node is a block parameter
type Node struct {
	id            conflow.ID
	nameNode      *conflow.IDNode
	valueNode     parsley.Node
	evalStage     conflow.EvalStage
	isDeclaration bool
	dependencies  conflow.Dependencies
	directives    []conflow.BlockNode
	schema        schema.Schema
}

// NewNode creates a new block parameter node
func NewNode(
	blockID conflow.ID,
	nameNode *conflow.IDNode,
	valueNode parsley.Node,
	isDeclaration bool,
	directives []conflow.BlockNode,
) *Node {
	return &Node{
		id:            conflow.ID(fmt.Sprintf("%s.%s", blockID, nameNode.ID())),
		nameNode:      nameNode,
		valueNode:     valueNode,
		isDeclaration: isDeclaration,
		directives:    directives,
	}
}

// ID returns with the name of the parameter
func (n *Node) ID() conflow.ID {
	return n.id
}

// ID returns with the name of the parameter
func (n *Node) Name() conflow.ID {
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

// Schema returns the schema for the node's value
func (n *Node) Schema() interface{} {
	if n.schema != nil {
		return n.schema
	}

	return n.valueNode.Schema()
}

// EvalStage returns with the evaluation stage
func (n *Node) EvalStage() conflow.EvalStage {
	if n.evalStage == conflow.EvalStageUndefined {
		return conflow.EvalStageMain
	}
	return n.evalStage
}

// Dependencies returns the blocks/parameters this parameter depends on
func (n *Node) Dependencies() conflow.Dependencies {
	if n.dependencies != nil {
		return n.dependencies
	}

	n.dependencies = make(conflow.Dependencies)

	parsley.Walk(n.valueNode, func(node parsley.Node) bool {
		if v, ok := node.(conflow.VariableNode); ok {
			n.dependencies[v.ID()] = v
		}
		return false
	})

	return n.dependencies
}

// Directives returns the parameter directives
func (n *Node) Directives() []conflow.BlockNode {
	return n.directives
}

// Provides returns with nil as a parameter node doesn't define other nodes
func (n *Node) Provides() []conflow.ID {
	return nil
}

// Generates returns with nil as a parameter node doesn't generate other nodes
func (n *Node) Generates() []conflow.ID {
	return nil
}

// IsDeclaration returns true if the parameter was declared in the block
func (n *Node) IsDeclaration() bool {
	return n.isDeclaration
}

// SetSchema sets the schema for the parameter node
func (n *Node) SetSchema(s schema.Schema) {
	n.schema = s
}

// StaticCheck runs a static analysis on the value node
func (n *Node) StaticCheck(ctx interface{}) parsley.Error {
	switch vn := n.valueNode.(type) {
	case parsley.StaticCheckable:
		if err := vn.StaticCheck(ctx); err != nil {
			return err
		}
	case parsley.LiteralNode:
		if n.schema != nil && vn.Value() != nil {
			if err := n.schema.ValidateValue(vn.Value()); err != nil {
				return parsley.NewError(n.valueNode.Pos(), err)
			}
		}
	}

	if n.schema != nil && n.valueNode.Schema().(schema.Schema).Type() != schema.TypeNull {
		if err := n.schema.ValidateSchema(n.valueNode.Schema().(schema.Schema), false); err != nil {
			return parsley.NewError(n.valueNode.Pos(), err)
		}
	}

	return nil
}

// Value returns with the value of the node
func (n *Node) Value(ctx interface{}) (interface{}, parsley.Error) {
	value, err := parsley.EvaluateNode(ctx, n.valueNode)
	if err != nil {
		return nil, err
	}

	if value != nil && n.schema != nil {
		if err := n.schema.ValidateValue(value); err != nil {
			return nil, parsley.NewError(n.valueNode.Pos(), err)
		}
	}

	return value, nil
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
	for _, node := range n.directives {
		if parsley.Walk(node, f) {
			return true
		}
	}

	return parsley.Walk(n.valueNode, f)
}

func (n *Node) CreateContainer(
	ctx *conflow.EvalContext,
	_ conflow.RuntimeConfig,
	parent conflow.BlockContainer,
	value interface{},
	wgs []conflow.WaitGroup,
	pending bool,
) conflow.JobContainer {
	return NewContainer(ctx, n, parent, value, wgs, pending)
}

// String returns with a string representation of the node
func (n *Node) String() string {
	return fmt.Sprintf("%s{%s=%s, %d..%d}", n.Token(), n.nameNode, n.valueNode, n.Pos(), n.ReaderPos())
}
