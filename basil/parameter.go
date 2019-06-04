package basil

import "github.com/opsidian/parsley/parsley"

// ParameterDescriptor describes a parameter
type ParameterDescriptor struct {
	Type       string
	EvalStage  EvalStage
	IsRequired bool
	IsOutput   bool
}

// ParameterNode is the AST node for a parameter
//go:generate counterfeiter . ParameterNode
type ParameterNode interface {
	Node
	Name() ID
	ValueNode() parsley.Node
	IsDeclaration() bool
	SetDescriptor(ParameterDescriptor)
}

// ParameterContainer is a parameter container
type ParameterContainer interface {
	Container
	Node() ParameterNode
	Parent() BlockContainer
}
