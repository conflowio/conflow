// Code generated by Conflow. DO NOT EDIT.

package strings

import (
	"github.com/conflowio/conflow/src/schema"
)

// TrimPrefixInterpreter is the conflow interpreter for the TrimPrefix function
type TrimPrefixInterpreter struct {
	s schema.Schema
}

func (i TrimPrefixInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Function{
			Metadata: schema.Metadata{
				Description: "It returns s without the provided leading prefix string.\nIf s doesn't start with prefix, s is returned unchanged.",
			},
			Parameters: schema.Parameters{
				schema.NamedSchema{
					Name:   "s",
					Schema: &schema.String{},
				},
				schema.NamedSchema{
					Name:   "prefix",
					Schema: &schema.String{},
				},
			},
			Result: &schema.String{},
		}
	}
	return i.s
}

// Eval returns with the result of the function
func (i TrimPrefixInterpreter) Eval(ctx interface{}, args []interface{}) (interface{}, error) {
	var val0 = args[0].(string)
	var val1 = args[1].(string)
	return TrimPrefix(val0, val1), nil
}