package parameter

import (
	"fmt"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/parsley/parsley"
)

var _ basil.BlockParamNode = &Node{}

// Node is a block parameter
type Node struct {
	id            basil.ID
	nameNode      *basil.IDNode
	valueNode     parsley.Node
	evalStage     basil.EvalStage
	isDeclaration bool
	dependencies  []basil.VariableNode
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
	return n.evalStage
}

// Dependencies returns the blocks/parameters this parameter depends on
func (n *Node) Dependencies() []basil.VariableNode {
	if n.dependencies != nil {
		return n.dependencies
	}

	n.dependencies = make([]basil.VariableNode, 0)

	parsley.Walk(n.valueNode, func(node parsley.Node) bool {
		if v, ok := node.(basil.VariableNode); ok {
			n.dependencies = append(n.dependencies, v)
		}
		return false
	})

	return n.dependencies
}

// Provides returns with nil as a parameter node doesn't define any other nodes
func (n *Node) Provides() []basil.ID {
	return nil
}

// IsDeclaration returns true if the parameter was declared in the block
func (n *Node) IsDeclaration() bool {
	return n.isDeclaration
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

// String returns with a string representation of the node
func (n *Node) String() string {
	return fmt.Sprintf("%s{%s=%s, %d..%d}", n.Token(), n.nameNode, n.valueNode, n.Pos(), n.ReaderPos())
}

// Walk runs the given function on all child nodes
func (n *Node) Walk(f func(n parsley.Node) bool) bool {
	return parsley.Walk(n.valueNode, f)
}
