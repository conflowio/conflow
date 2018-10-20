package function

// Registry is an interface for a function registry
//go:generate counterfeiter . Registry
type Registry interface {
	Callable
	FunctionExists(name string) bool
	RegisterFunction(name string, callable Callable)
}
