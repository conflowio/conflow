package basil

import "errors"

// ErrVariableNotFound error is used when the given variable does not exist
var ErrVariableNotFound = errors.New("Variable not found")

// VariableProvider is an interface for looking up variables
//go:generate counterfeiter . VariableProvider
type VariableProvider interface {
	GetVar(name string) (interface{}, bool)
	LookupVar(lookup VariableLookUp) (interface{}, error)
}

// VariableProviderAware defines a function to access a variable provider
type VariableProviderAware interface {
	GetVariableProvider() VariableProvider
}

// VariableLookUp is a variable lookup function
type VariableLookUp func(provider VariableProvider) (interface{}, error)
