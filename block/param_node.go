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

// Type returns with the value node's type
func (p *ParamNode) Type() string {
	return p.valueNode.Type()
}

// StaticCheck runs a static analysis on the value node
func (p *ParamNode) StaticCheck(ctx interface{}) parsley.Error {
	switch n := p.valueNode.(type) {
	case parsley.StaticCheckable:
		if err := n.StaticCheck(ctx); err != nil {
			return err
		}
	}
	return nil
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

// KeyNode returns with the key node
func (p *ParamNode) KeyNode() parsley.Node {
	return p.keyNode
}

// ValueNode returns with the value node
func (p *ParamNode) ValueNode() parsley.Node {
	return p.valueNode
}
