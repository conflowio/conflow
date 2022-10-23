// Code generated by Conflow. DO NOT EDIT.

package strings

import (
	"github.com/conflowio/conflow/pkg/schema"
)

func init() {
	schema.Register(&schema.Function{
		Metadata: schema.Metadata{
			Description: "It returns a copy of the string s with all\nnon-overlapping instances of old replaced by new.",
			ID:          "github.com/conflowio/conflow/pkg/functions/strings.Replace",
		},
		Parameters: schema.Parameters{
			schema.NamedSchema{
				Name:   "s",
				Schema: &schema.String{},
			},
			schema.NamedSchema{
				Name:   "old",
				Schema: &schema.String{},
			},
			schema.NamedSchema{
				Name:   "new",
				Schema: &schema.String{},
			},
		},
		Result: &schema.String{},
	})
}

// ReplaceInterpreter is the Conflow interpreter for the Replace function
type ReplaceInterpreter struct {
}

func (i ReplaceInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/functions/strings.Replace")
	return s
}

// Eval returns with the result of the function
func (i ReplaceInterpreter) Eval(ctx interface{}, args []interface{}) (interface{}, error) {
	var val0 = args[0].(string)
	var val1 = args[1].(string)
	var val2 = args[2].(string)
	return Replace(val0, val1, val2), nil
}