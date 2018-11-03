package block

import (
	"fmt"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/parsley"
)

// ParamNode is a block parameter
type ParamNode struct {
	keyNode   parsley.Node
	valueNode parsley.Node
}

// NewParamNode creates a new block parameter node
func NewParamNode(keyNode parsley.Node, valueNode parsley.Node) *ParamNode {
	return &ParamNode{
		keyNode:   keyNode,
		valueNode: valueNode,
	}
}

// Token returns with the node token
func (p *ParamNode) Token() string {
	return "BLOCK_PARAM"
}

// Value returns with the value of the node
func (p *ParamNode) Value(ctx interface{}) (interface{}, parsley.Error) {
	return p.valueNode.Value(ctx)
}

// Pos returns the position
func (p *ParamNode) Pos() parsley.Pos {
	return p.keyNode.Pos()
}

// ReaderPos returns the position of the first character immediately after this node
func (p *ParamNode) ReaderPos() parsley.Pos {
	return p.valueNode.ReaderPos()
}

// String returns with a string representation of the node
func (p *ParamNode) String() string {
	return fmt.Sprintf("%s{%s=%s, %d..%d}", p.Token(), p.keyNode, p.valueNode, p.Pos(), p.ReaderPos())
}

// Walk runs the given function on all child nodes
func (p *ParamNode) Walk(f func(n parsley.Node) bool) bool {
	if ast.WalkNode(p.keyNode, f) {
		return true
	}
	if ast.WalkNode(p.valueNode, f) {
		return true
	}

	return false
}
