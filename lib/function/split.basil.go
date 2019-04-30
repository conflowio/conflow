// Code generated by Basil. DO NOT EDIT.
package function

import (
	"fmt"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/variable"
	"github.com/opsidian/parsley/parsley"
)

// SplitInterpreter is an AST node interpreter for the Split function
type SplitInterpreter struct{}

// StaticCheck runs a static analysis on the function parameters
func (i SplitInterpreter) StaticCheck(ctx interface{}, node basil.FunctionNode) (string, parsley.Error) {
	if len(node.ArgumentNodes()) != 2 {
		return "", parsley.NewError(node.Pos(), fmt.Errorf("%s expects 2 arguments", node.Name()))
	}

	arguments := node.ArgumentNodes()

	if err := variable.CheckNodeType(arguments[0], "string"); err != nil {
		return "", err
	}

	if err := variable.CheckNodeType(arguments[1], "string"); err != nil {
		return "", err
	}

	return "[]interface{}", nil

}

// Eval returns with the result of the function
func (i SplitInterpreter) Eval(ctx interface{}, node basil.FunctionNode) (interface{}, parsley.Error) {
	arguments := node.ArgumentNodes()

	arg0, err := variable.NodeStringValue(arguments[0], ctx)
	if err != nil {
		return nil, err
	}

	arg1, err := variable.NodeStringValue(arguments[1], ctx)
	if err != nil {
		return nil, err
	}

	return Split(
		arg0,
		arg1,
	), nil

}
