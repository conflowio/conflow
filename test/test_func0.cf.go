// Code generated by Conflow. DO NOT EDIT.

package test

import (
	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/schema"
	"github.com/opsidian/parsley/parsley"
)

// TestFunc0Interpreter is the conflow interpreter for the testFunc0 function
type TestFunc0Interpreter struct {
	s schema.Schema
}

func (i TestFunc0Interpreter) Schema() schema.Schema {
	if i.s == nil {
		i.s = &schema.Function{
			Result: &schema.String{},
		}
	}
	return i.s
}

// Eval returns with the result of the function
func (i TestFunc0Interpreter) Eval(ctx interface{}, node conflow.FunctionNode) (interface{}, parsley.Error) {
	return testFunc0(), nil
}
