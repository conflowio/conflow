// Code generated by Conflow. DO NOT EDIT.

package strings

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
				Description: "It reports whether substr is within s.",
			},
			Parameters: schema.Parameters{
				schema.NamedSchema{
					Name:   "s",
					Schema: &schema.String{},
				},
				schema.NamedSchema{
					Name:   "substr",
					Schema: &schema.String{},
				},
			},
			Result: &schema.Boolean{},
		}
	}
	return i.s
}

// Eval returns with the result of the function
func (i ContainsInterpreter) Eval(ctx interface{}, args []interface{}) (interface{}, error) {
	var val0 = args[0].(string)
	var val1 = args[1].(string)
	return Contains(val0, val1), nil
}