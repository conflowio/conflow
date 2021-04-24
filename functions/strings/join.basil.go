// Code generated by Basil. DO NOT EDIT.

package strings

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/schema"
	"github.com/opsidian/parsley/parsley"
)

// JoinInterpreter is the basil interpreter for the Join function
type JoinInterpreter struct {
	s schema.Schema
}

func (i JoinInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Function{
			Metadata: schema.Metadata{
				Description: "It concatenates the elements of a to create a single string. The separator string\nsep is placed between elements in the resulting string.",
			},
			Parameters: schema.Parameters{
				schema.NamedSchema{
					Name: "a",
					Schema: &schema.Array{
						Items: &schema.String{},
					},
				},
				schema.NamedSchema{
					Name:   "sep",
					Schema: &schema.String{},
				},
			},
			Result: &schema.String{},
		}
	}
	return i.s
}

// Eval returns with the result of the function
func (i JoinInterpreter) Eval(ctx interface{}, node basil.FunctionNode) (interface{}, parsley.Error) {
	parameters := i.Schema().(*schema.Function).GetParameters()
	arguments := node.ArgumentNodes()

	arg0, evalErr := parsley.EvaluateNode(ctx, arguments[0])
	if evalErr != nil {
		return nil, evalErr
	}
	if err := parameters[0].Schema.ValidateValue(arg0); err != nil {
		return nil, parsley.NewError(arguments[0].Pos(), err)
	}
	var val0 = make([]string, len(arg0.([]interface{})))
	for arg0k, arg0v := range arg0.([]interface{}) {
		val0[arg0k] = arg0v.(string)
	}

	arg1, evalErr := parsley.EvaluateNode(ctx, arguments[1])
	if evalErr != nil {
		return nil, evalErr
	}
	if err := parameters[1].Schema.ValidateValue(arg1); err != nil {
		return nil, parsley.NewError(arguments[1].Pos(), err)
	}
	var val1 = arg1.(string)

	return Join(val0, val1), nil
}
