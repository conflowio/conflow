// Code generated by Conflow. DO NOT EDIT.

package test

import (
	"github.com/conflowio/conflow/pkg/schema"
)

func init() {
	schema.Register(&schema.Function{
		Metadata: schema.Metadata{
			ID: "github.com/conflowio/conflow/pkg/test.testFunc2",
		},
		Parameters: schema.Parameters{
			schema.NamedSchema{
				Name:   "str1",
				Schema: &schema.String{},
			},
			schema.NamedSchema{
				Name:   "str2",
				Schema: &schema.String{},
			},
		},
		Result: &schema.String{},
	})
}

// TestFunc2Interpreter is the Conflow interpreter for the testFunc2 function
type TestFunc2Interpreter struct {
}

func (i TestFunc2Interpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/test.testFunc2")
	return s
}

// Eval returns with the result of the function
func (i TestFunc2Interpreter) Eval(ctx interface{}, args []interface{}) (interface{}, error) {
	var val0 = args[0].(string)
	var val1 = args[1].(string)
	return testFunc2(val0, val1), nil
}