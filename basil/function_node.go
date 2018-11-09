package basil

import (
	"github.com/opsidian/basil/variable"
	"github.com/opsidian/parsley/parsley"
)

// FunctionNode is the AST node for a function
//go:generate counterfeiter . FunctionNode
type FunctionNode interface {
	parsley.Node
	parsley.StaticCheckable
	Name() variable.ID
	ArgumentNodes() []parsley.Node
}

// FunctionRegistryAware is an interface to get a function node transformer registry
type FunctionRegistryAware interface {
	FunctionRegistry() parsley.NodeTransformerRegistry
}
