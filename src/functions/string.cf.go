// Code generated by Conflow. DO NOT EDIT.

package functions

import (
	"github.com/conflowio/conflow/src/schema"
)

// StringInterpreter is the conflow interpreter for the String function
type StringInterpreter struct {
	s schema.Schema
}

func (i StringInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Function{
			Metadata: schema.Metadata{
				Description: "It converts the given value to a string",
			},
			Parameters: schema.Parameters{
				schema.NamedSchema{
					Name:   "value",
					Schema: &schema.Untyped{},
				},
			},
			Result: &schema.String{},
		}
	}
	return i.s
}

// Eval returns with the result of the function
func (i StringInterpreter) Eval(ctx interface{}, args []interface{}) (interface{}, error) {
	var val0 = args[0]
	return String(val0)
}
