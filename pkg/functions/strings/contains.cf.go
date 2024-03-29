// Code generated by Conflow. DO NOT EDIT.

package strings

import (
	"github.com/conflowio/conflow/pkg/schema"
)

func init() {
	schema.Register(&schema.Function{
		Metadata: schema.Metadata{
			Description: "It reports whether substr is within s.",
			ID:          "github.com/conflowio/conflow/pkg/functions/strings.Contains",
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
	})
}

// ContainsInterpreter is the Conflow interpreter for the Contains function
type ContainsInterpreter struct {
}

func (i ContainsInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/functions/strings.Contains")
	return s
}

// Eval returns with the result of the function
func (i ContainsInterpreter) Eval(ctx interface{}, args []interface{}) (interface{}, error) {
	var val0 = schema.Value[string](args[0])
	var val1 = schema.Value[string](args[1])
	return Contains(val0, val1), nil
}
