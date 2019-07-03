// Code generated by Basil. DO NOT EDIT.
package json

import (
	"fmt"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/function"
	"github.com/opsidian/basil/basil/variable"
	"github.com/opsidian/parsley/parsley"
)

// EncodeInterpreter is an AST node interpreter for the Encode function
type EncodeInterpreter struct{}

// StaticCheck runs a static analysis on the function parameters
func (i EncodeInterpreter) StaticCheck(ctx interface{}, node basil.FunctionNode) (string, parsley.Error) {
	if len(node.ArgumentNodes()) != 1 {
		return "", parsley.NewError(node.Pos(), fmt.Errorf("%s expects 1 arguments", node.Name()))
	}

	arguments := node.ArgumentNodes()

	if err := variable.CheckNodeType(arguments[0], "interface{}"); err != nil {
		return "", err
	}

	return "string", nil

}

// Eval returns with the result of the function
func (i EncodeInterpreter) Eval(ctx interface{}, node basil.FunctionNode) (interface{}, parsley.Error) {
	arguments := node.ArgumentNodes()

	arg0, evalErr := arguments[0].Value(ctx)
	if evalErr != nil {
		return nil, evalErr
	}

	val0, convertErr := variable.AnyValue(arg0)
	if convertErr != nil {
		return nil, parsley.NewError(arguments[0].Pos(), convertErr)
	}

	res, resErr := Encode(val0)
	if resErr != nil {
		if funcErr, ok := resErr.(*function.Error); ok {
			return nil, parsley.NewError(arguments[funcErr.ArgIndex].Pos(), funcErr.Err)
		}
		return nil, parsley.NewError(node.Pos(), resErr)
	}

	return res, nil

}