// Code generated by Conflow. DO NOT EDIT.

package math

import (
	"github.com/conflowio/conflow/src/schema"
)

// RoundInterpreter is the conflow interpreter for the Round function
type RoundInterpreter struct {
	s schema.Schema
}

func (i RoundInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Function{
			Metadata: schema.Metadata{
				Description: "It returns the nearest integer, rounding half away from zero.",
			},
			Parameters: schema.Parameters{
				schema.NamedSchema{
					Name:   "value",
					Schema: &schema.Number{},
				},
			},
			Result: &schema.Number{},
		}
	}
	return i.s
}

// Eval returns with the result of the function
func (i RoundInterpreter) Eval(ctx interface{}, args []interface{}) (interface{}, error) {
	var val0 = args[0].(float64)
	return Round(val0), nil
}
