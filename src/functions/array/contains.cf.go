// Code generated by Conflow. DO NOT EDIT.

package array

import (
	"github.com/conflowio/conflow/src/schema"
)

// ContainsInterpreter is the conflow interpreter for the Contains function
type ContainsInterpreter struct {
	s schema.Schema
}

func (i ContainsInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Function{
			Metadata: schema.Metadata{
				Description: "It returns true if the array contains the given element",
			},
			Parameters: schema.Parameters{
				schema.NamedSchema{
					Name: "arr",
					Schema: &schema.Array{
						Items: &schema.Untyped{},
					},
				},
				schema.NamedSchema{
					Name:   "elem",
					Schema: &schema.Untyped{},
				},
			},
			Result: &schema.Boolean{},
		}
	}
	return i.s
}

// Eval returns with the result of the function
func (i ContainsInterpreter) Eval(ctx interface{}, args []interface{}) (interface{}, error) {
	var val0 = args[0].([]interface{})
	var val1 = args[1]
	return Contains(val0, val1)
}
