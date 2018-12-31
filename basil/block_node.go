package basil

import (
	"github.com/opsidian/parsley/parsley"
)

// DependencyAware defines an interface to return with dependencies
type DependencyAware interface {
	Dependencies() []parsley.Node
}

// BlockNode is the AST node for a block
//go:generate counterfeiter . BlockNode
type BlockNode interface {
	parsley.NonTerminalNode
	Identifiable
	DependencyAware
	EvalStageAware
	IDNode() parsley.Node
	TypeNode() parsley.Node
	ParamType(ID) (string, bool)
}

// BlockNodeRegistry is an interface for looking up named blocks
//go:generate counterfeiter . BlockNodeRegistry
type BlockNodeRegistry interface {
	BlockNode(ID) (BlockNode, bool)
	AddBlockNode(BlockNode) error
}

// BlockNodeRegistryAware defines an interface to access a block node registry
type BlockNodeRegistryAware interface {
	BlockNodeRegistry() BlockNodeRegistry
}

// BlockParamNode is the AST node for a block parameter
type BlockParamNode interface {
	parsley.Node
	Identifiable
	DependencyAware
	EvalStageAware
	KeyNode() parsley.Node
	ValueNode() parsley.Node
}

// BlockChildNode is an interface for block child nodes (blocks or params)
type BlockChildNode interface {
	parsley.Node
	Identifiable
	DependencyAware
	EvalStageAware
}

// BlockTransformerRegistryAware is an interface to get a block node transformer registry
type BlockTransformerRegistryAware interface {
	BlockTransformerRegistry() parsley.NodeTransformerRegistry
}
