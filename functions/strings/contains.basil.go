// Code generated by Basil. DO NOT EDIT.

package strings

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/schema"
	"github.com/opsidian/parsley/parsley"
)

// ContainsInterpreter is the basil interpreter for the Contains function
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
func (i ContainsInterpreter) Eval(ctx interface{}, node basil.FunctionNode) (interface{}, parsley.Error) {
	parameters := i.Schema().(*schema.Function).GetParameters()
	arguments := node.ArgumentNodes()

	arg0, evalErr := parsley.EvaluateNode(ctx, arguments[0])
	if evalErr != nil {
		return nil, evalErr
	}
	if err := parameters[0].Schema.ValidateValue(arg0); err != nil {
		return nil, parsley.NewError(arguments[0].Pos(), err)
	}
	var val0 = arg0.(string)

	arg1, evalErr := parsley.EvaluateNode(ctx, arguments[1])
	if evalErr != nil {
		return nil, evalErr
	}
	if err := parameters[1].Schema.ValidateValue(arg1); err != nil {
		return nil, parsley.NewError(arguments[1].Pos(), err)
	}
	var val1 = arg1.(string)

	return Contains(val0, val1), nil
}
