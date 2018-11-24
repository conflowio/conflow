package basil

import (
	"github.com/opsidian/parsley/parsley"
)

// FunctionNode is the AST node for a function
//go:generate counterfeiter . FunctionNode
type FunctionNode interface {
	parsley.Node
	parsley.StaticCheckable
	Identifiable
	ArgumentNodes() []parsley.Node
}

// FunctionTransformerRegistryAware is an interface to get a function node transformer registry
type FunctionTransformerRegistryAware interface {
	FunctionTransformerRegistry() parsley.NodeTransformerRegistry
}
