// Code generated by Conflow. DO NOT EDIT.

package array

import (
	"github.com/conflowio/conflow/pkg/schema"
)

func init() {
	schema.Register(&schema.Function{
		Metadata: schema.Metadata{
			Description: "It returns true if the array contains the given element",
			ID:          "github.com/conflowio/conflow/pkg/functions/array.Contains",
		},
		Parameters: schema.Parameters{
			schema.NamedSchema{
				Name: "arr",
				Schema: &schema.Array{
					Items: &schema.Any{},
				},
			},
			schema.NamedSchema{
				Name:   "elem",
				Schema: &schema.Any{},
			},
		},
		Result: &schema.Boolean{},
	})
}

// ContainsInterpreter is the Conflow interpreter for the Contains function
type ContainsInterpreter struct {
}

func (i ContainsInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/functions/array.Contains")
	return s
}

// Eval returns with the result of the function
func (i ContainsInterpreter) Eval(ctx interface{}, args []interface{}) (interface{}, error) {
	var val0 = args[0].([]interface{})
	var val1 = args[1]
	return Contains(val0, val1)
}
