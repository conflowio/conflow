package test

import (
	"github.com/opsidian/basil/variable"
)

var testVariableProvider = variableProvider{map[variable.ID]interface{}{
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

type variableProvider struct {
	Vars map[variable.ID]interface{}
}

func (v variableProvider) GetVar(name variable.ID) (interface{}, bool) {
	value, ok := v.Vars[name]
	return value, ok
}

func (v variableProvider) LookupVar(lookup variable.LookUp) (interface{}, error) {
	return lookup(v)
}
