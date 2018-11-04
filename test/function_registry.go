package test

import (
	"errors"
	"strings"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/parsley/parsley"
)

type functionRegistry struct{}

func (f *functionRegistry) RegisterFunction(name basil.ID, callable basil.Callable) {
}

func (f *functionRegistry) FunctionExists(name basil.ID) bool {
	return name == basil.ID("rand") || name == basil.ID("upper") || name == basil.ID("default")
}

func (f *functionRegistry) CallFunction(ctx interface{}, function parsley.Node, params []parsley.Node) (interface{}, parsley.Error) {
	name, _ := function.Value(ctx)
	switch name {
	case basil.ID("rand"):
		return int64(123), nil
	case basil.ID("upper"):
		value, err := params[0].Value(ctx)
		if err != nil {
			return nil, err
		}
		return strings.ToUpper(value.(string)), nil
	case basil.ID("default"):
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
