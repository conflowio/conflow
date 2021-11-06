// Code generated by Conflow. DO NOT EDIT.

package json

import (
	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/function"
	"github.com/conflowio/conflow/conflow/schema"
	"github.com/conflowio/parsley/parsley"
)

// DecodeInterpreter is the conflow interpreter for the Decode function
type DecodeInterpreter struct {
	s schema.Schema
}

func (i DecodeInterpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Function{
			Metadata: schema.Metadata{
				Description: "It converts the given json string to a data structure",
			},
			Parameters: schema.Parameters{
				schema.NamedSchema{
					Name:   "jsonStr",
					Schema: &schema.String{},
				},
			},
			Result: &schema.Untyped{},
		}
	}
	return i.s
}

// Eval returns with the result of the function
func (i DecodeInterpreter) Eval(ctx interface{}, node conflow.FunctionNode) (interface{}, parsley.Error) {
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

	res, resErr := Decode(val0)
	if resErr != nil {
		if funcErr, ok := resErr.(*function.Error); ok {
			return nil, parsley.NewError(arguments[funcErr.ArgIndex].Pos(), funcErr.Err)
		}
		return nil, parsley.NewError(node.Pos(), resErr)
	}

	return res, nil
}