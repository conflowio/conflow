// Code generated by Conflow. DO NOT EDIT.

package test

import (
	"github.com/conflowio/conflow/src/schema"
)

// TestFunc1Interpreter is the conflow interpreter for the testFunc1 function
type TestFunc1Interpreter struct {
	s schema.Schema
}

func (i TestFunc1Interpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Function{
			Parameters: schema.Parameters{
				schema.NamedSchema{
					Name:   "str",
					Schema: &schema.String{},
				},
			},
			Result: &schema.String{},
		}
	}
	return i.s
}

// Eval returns with the result of the function
func (i TestFunc1Interpreter) Eval(ctx interface{}, args []interface{}) (interface{}, error) {
	var val0 = args[0].(string)
	return testFunc1(val0), nil
}