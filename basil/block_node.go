package basil

import (
	"github.com/opsidian/parsley/parsley"
)

// BlockNode is the AST node for a block
//go:generate counterfeiter . BlockNode
type BlockNode interface {
	parsley.Node
	parsley.StaticCheckable
	Eval(ctx interface{}, stage string, block Block) parsley.Error
	Identifiable
	ParamNodes() map[ID]BlockParamNode
	BlockNodes() []BlockNode
}

// BlockParamNode is the AST node for a block parameter
type BlockParamNode interface {
	parsley.Node
	KeyNode() parsley.Node
	ValueNode() parsley.Node
}

// BlockRegistryAware is an interface to get a block node transformer registry
type BlockRegistryAware interface {
	BlockRegistry() parsley.NodeTransformerRegistry
}
