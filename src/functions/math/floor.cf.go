// Code generated by Conflow. DO NOT EDIT.

package math

import (
	"github.com/conflowio/conflow/src/schema"
)

// FloorInterpreter is the conflow interpreter for the Floor function
type FloorInterpreter struct {
	s schema.Schema
}

func (i FloorInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Function{
			Metadata: schema.Metadata{
				Description: "It returns the greatest integer value less than or equal to x.",
			},
			Parameters: schema.Parameters{
				schema.NamedSchema{
					Name: "number",
					Schema: &schema.Untyped{
						Types: []string{"integer", "number"},
					},
				},
			},
			Result: &schema.Integer{},
		}
	}
	return i.s
}

// Eval returns with the result of the function
func (i FloorInterpreter) Eval(ctx interface{}, args []interface{}) (interface{}, error) {
	var val0 = args[0]
	return Floor(val0), nil
}