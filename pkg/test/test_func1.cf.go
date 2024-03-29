// Code generated by Conflow. DO NOT EDIT.

package test

import (
	"github.com/conflowio/conflow/pkg/schema"
)

func init() {
	schema.Register(&schema.Function{
		Metadata: schema.Metadata{
			ID: "github.com/conflowio/conflow/pkg/test.testFunc1",
		},
		Parameters: schema.Parameters{
			schema.NamedSchema{
				Name:   "str",
				Schema: &schema.String{},
			},
		},
		Result: &schema.String{},
	})
}

// TestFunc1Interpreter is the Conflow interpreter for the testFunc1 function
type TestFunc1Interpreter struct {
}

func (i TestFunc1Interpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/test.testFunc1")
	return s
}

// Eval returns with the result of the function
func (i TestFunc1Interpreter) Eval(ctx interface{}, args []interface{}) (interface{}, error) {
	var val0 = schema.Value[string](args[0])
	return testFunc1(val0), nil
}
