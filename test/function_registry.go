package test

import (
	"errors"
	"strings"

	"github.com/opsidian/ocl/ocl"
	"github.com/opsidian/parsley/parsley"
)

type functionRegistry struct{}

func (f *functionRegistry) RegisterFunction(name string, callable ocl.Callable) {
}

func (f *functionRegistry) FunctionExists(name string) bool {
	return name == "rand" || name == "upper" || name == "default"
}

func (f *functionRegistry) CallFunction(ctx interface{}, function parsley.Node, params []parsley.Node) (interface{}, parsley.Error) {
	name, _ := function.Value(ctx)
	switch name {
	case "rand":
		return int64(123), nil
	case "upper":
		value, err := params[0].Value(ctx)
		if err != nil {
			return nil, err
		}
		return strings.ToUpper(value.(string)), nil
	case "default":
		value1, err := params[0].Value(ctx)
		if err != nil {
			return nil, err
		}
		value2, err := params[1].Value(ctx)
		if err != nil {
			return nil, err
		}
		if value1 != nil {
			return value1, nil
		}

		return value2, nil
	}
	return nil, parsley.NewError(function.Pos(), errors.New("unknown function"))
}
