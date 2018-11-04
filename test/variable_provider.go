package test

import "github.com/opsidian/basil/basil"

var testVariableProvider = variableProvider{map[basil.ID]interface{}{
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
	Vars map[basil.ID]interface{}
}

func (v variableProvider) GetVar(name basil.ID) (interface{}, bool) {
	value, ok := v.Vars[name]
	return value, ok
}

func (v variableProvider) LookupVar(lookup basil.VariableLookUp) (interface{}, error) {
	return lookup(v)
}
