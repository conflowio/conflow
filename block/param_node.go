package block

import (
	"fmt"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/parsley/parsley"
)

var _ basil.BlockParamNode = &ParamNode{}

// ParamNode is a block parameter
type ParamNode struct {
	id            basil.ID
	nameNode      *basil.IDNode
	valueNode     parsley.Node
	evalStage     basil.EvalStage
	isDeclaration bool
}

// NewParamNode creates a new block parameter node
func NewParamNode(blockID basil.ID, nameNode *basil.IDNode, valueNode parsley.Node, isDeclaration bool) *ParamNode {
	return &ParamNode{
		id:            basil.ID(fmt.Sprintf("%s.%s", blockID, nameNode.ID())),
		nameNode:      nameNode,
		valueNode:     valueNode,
		isDeclaration: isDeclaration,
	}
}

// ID returns with the name of the parameter
func (p *ParamNode) ID() basil.ID {
	return p.id
}

// ID returns with the name of the parameter
func (p *ParamNode) Name() basil.ID {
	return p.nameNode.ID()
}

// Token returns with the node token
func (p *ParamNode) Token() string {
	return "BLOCK_PARAM"
}

// Type returns with the value node's type
func (p *ParamNode) Type() string {
	return p.valueNode.Type()
}

// EvalStage returns with the evaluation stage
func (p *ParamNode) EvalStage() basil.EvalStage {
	return p.evalStage
}

// Dependencies returns the blocks/parameters this parameter depends on
func (p *ParamNode) Dependencies() []basil.IdentifiableNode {
	var dependencies []basil.IdentifiableNode

	parsley.Walk(p.valueNode, func(n parsley.Node) bool {
		if v, ok := n.(basil.VariableNode); ok {
			dependencies = append(dependencies, v)
		}
		return false
	})

	return dependencies
}

// Provides returns with nil as a parameter node doesn't define any other nodes
func (p *ParamNode) Provides() []basil.ID {
	return nil
}

// IsDeclaration returns true if the parameter was declared in the block
func (p *ParamNode) IsDeclaration() bool {
	return p.isDeclaration
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
	return p.nameNode.Pos()
}

// ReaderPos returns the position of the first character immediately after this node
func (p *ParamNode) ReaderPos() parsley.Pos {
	return p.valueNode.ReaderPos()
}

// String returns with a string representation of the node
func (p *ParamNode) String() string {
	return fmt.Sprintf("%s{%s=%s, %d..%d}", p.Token(), p.nameNode, p.valueNode, p.Pos(), p.ReaderPos())
}

// Walk runs the given function on all child nodes
func (p *ParamNode) Walk(f func(n parsley.Node) bool) bool {
	return parsley.Walk(p.valueNode, f)
}
