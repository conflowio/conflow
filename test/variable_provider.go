package test

import "github.com/opsidian/ocl/ocl"

var testVariableProvider = VariableProvider{map[string]interface{}{
	"foo": "bar",
	"testmap": map[string]interface{}{
		"key1": "value1",
		"key2": map[string]interface{}{
			"key3": "value3",
		},
		"key4": []interface{}{
			"value4",
		},
	},
	"arr": []interface{}{
		"value1",
		[]interface{}{
			"value2",
		},
		map[string]interface{}{
			"key1": "value3",
		},
	},
	"intkey": int64(1),
}}

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
func (v VariableProvider) LookupVar(lookup ocl.VariableLookUp) (interface{}, error) {
	return lookup(v)
}
