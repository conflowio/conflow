package ocl

import "github.com/opsidian/parsley/parsley"

// Callable is an interface for callable objects (general function interface)
//go:generate counterfeiter . Callable
type Callable interface {
	CallFunction(ctx interface{}, function parsley.Node, params []parsley.Node) (interface{}, parsley.Error)
}

// CallableFunc defines a helper to implement the Callable interface with functions
type CallableFunc func(ctx interface{}, function parsley.Node, nodes []parsley.Node) (interface{}, parsley.Error)

// CallFunction calls the function
func (f CallableFunc) CallFunction(ctx interface{}, function parsley.Node, nodes []parsley.Node) (interface{}, parsley.Error) {
	return f(ctx, function, nodes)
}

// FunctionRegistry is an interface for a function registry
//go:generate counterfeiter . FunctionRegistry
type FunctionRegistry interface {
	Callable
	FunctionExists(name string) bool
	RegisterFunction(name string, callable Callable)
}

// FunctionRegistryAware defines a function to access a function registry
type FunctionRegistryAware interface {
	GetFunctionRegistry() FunctionRegistry
}
