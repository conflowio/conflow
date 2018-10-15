package test

import (
	"github.com/opsidian/ocl/variable"
)

// VariableProvider is a test variable provider
type VariableProvider struct {
	Vars map[string]interface{}
}

// GetVar returns with the given variable
func (v VariableProvider) GetVar(name string) (interface{}, bool) {
	value, ok := v.Vars[name]
	return value, ok
}

// LookupVar looks up the given variable with a function
func (v VariableProvider) LookupVar(lookup variable.LookUp) (interface{}, error) {
	return lookup(v)
}
