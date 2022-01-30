// Code generated by Conflow. DO NOT EDIT.

package strings

import (
	"github.com/conflowio/conflow/src/schema"
)

// TrimSpaceInterpreter is the conflow interpreter for the TrimSpace function
type TrimSpaceInterpreter struct {
	s schema.Schema
}

func (i TrimSpaceInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Function{
			Metadata: schema.Metadata{
				Description: "It returns a slice of the string s, with all leading\nand trailing white space removed, as defined by Unicode.",
			},
			Parameters: schema.Parameters{
				schema.NamedSchema{
					Name:   "s",
					Schema: &schema.String{},
				},
			},
			Result: &schema.String{},
		}
	}
	return i.s
}

// Eval returns with the result of the function
func (i TrimSpaceInterpreter) Eval(ctx interface{}, args []interface{}) (interface{}, error) {
	var val0 = args[0].(string)
	return TrimSpace(val0), nil
}