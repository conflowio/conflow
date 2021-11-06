// Code generated by Conflow. DO NOT EDIT.

package test

import (
	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/schema"
	"github.com/conflowio/parsley/parsley"
)

// TestFunc2Interpreter is the conflow interpreter for the testFunc2 function
type TestFunc2Interpreter struct {
	s schema.Schema
}

func (i TestFunc2Interpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Function{
			Parameters: schema.Parameters{
				schema.NamedSchema{
					Name:   "str1",
					Schema: &schema.String{},
				},
				schema.NamedSchema{
					Name:   "str2",
					Schema: &schema.String{},
				},
			},
			Result: &schema.String{},
		}
	}
	return i.s
}

// Eval returns with the result of the function
func (i TestFunc2Interpreter) Eval(ctx interface{}, node conflow.FunctionNode) (interface{}, parsley.Error) {
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

	return testFunc2(val0, val1), nil
}