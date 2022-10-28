// Code generated by Conflow. DO NOT EDIT.

package strings

import (
	"github.com/conflowio/conflow/pkg/schema"
)

func init() {
	schema.Register(&schema.Function{
		Metadata: schema.Metadata{
			Description: "It returns s without the provided trailing suffix string.\nIf s doesn't end with suffix, s is returned unchanged.",
			ID:          "github.com/conflowio/conflow/pkg/functions/strings.TrimSuffix",
		},
		Parameters: schema.Parameters{
			schema.NamedSchema{
				Name:   "s",
				Schema: &schema.String{},
			},
			schema.NamedSchema{
				Name:   "suffix",
				Schema: &schema.String{},
			},
		},
		Result: &schema.String{},
	})
}

// TrimSuffixInterpreter is the Conflow interpreter for the TrimSuffix function
type TrimSuffixInterpreter struct {
}

func (i TrimSuffixInterpreter) Schema() schema.Schema {
	s, _ := schema.Get("github.com/conflowio/conflow/pkg/functions/strings.TrimSuffix")
	return s
}

// Eval returns with the result of the function
func (i TrimSuffixInterpreter) Eval(ctx interface{}, args []interface{}) (interface{}, error) {
	var val0 = schema.Value[string](args[0])
	var val1 = schema.Value[string](args[1])
	return TrimSuffix(val0, val1), nil
}
