// Code generated by Conflow. DO NOT EDIT.

package math

import (
	"github.com/conflowio/conflow/pkg/schema"
)

func init() {
	schema.Register(&schema.Function{
		Metadata: schema.Metadata{
			Description: "It returns the integer value of x.",
			ID:          "github.com/conflowio/conflow/pkg/functions/math.Trunc",
		},
		Parameters: schema.Parameters{
			schema.NamedSchema{
				Name:   "x",
				Schema: &schema.Number{},
			},
		},
		Result: &schema.Number{},
	})
}

// TruncInterpreter is the Conflow interpreter for the Trunc function
type TruncInterpreter struct {
}

func (i TruncInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/functions/math.Trunc")
	return s
}

// Eval returns with the result of the function
func (i TruncInterpreter) Eval(ctx interface{}, args []interface{}) (interface{}, error) {
	var val0 = args[0].(float64)
	return Trunc(val0), nil
}