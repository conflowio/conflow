package generator

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"

	"github.com/opsidian/basil/basil"
)

// Argument contains metadata for a function argument
type Argument struct {
	Name string
	Type string
}

// ParseArguments parses all arguments of a given go function
func ParseArguments(fun *ast.FuncType, file *ast.File) ([]*Argument, error) {
	arguments := make([]*Argument, 0, len(fun.Params.List))

	for _, argument := range fun.Params.List {
		parsedArguments, err := parseArgument(argument)
		if err != nil {
			return nil, err
		}

		for _, arg := range parsedArguments {
			if _, validType := basil.VariableTypes[arg.Type]; !validType {
				return nil, fmt.Errorf("invalid argument type for argument %q", arg.Name)
			}
		}

		arguments = append(arguments, parsedArguments...)
	}

	return arguments, nil
}

// ParseResults parses all results of a given go function
func ParseResults(fun *ast.FuncType, file *ast.File) ([]*Argument, error) {
	if len(fun.Results.List) == 0 || len(fun.Results.List) > 2 {
		return nil, fmt.Errorf("the function must return with one or two values")
	}

	arguments := make([]*Argument, 0, len(fun.Results.List))
	for _, argument := range fun.Results.List {
		parsedArguments, err := parseArgument(argument)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, parsedArguments...)
	}

	if arguments[0].Type == "error" {
		return nil, fmt.Errorf("the function must return a non-error value")
	}

	if _, validType := basil.VariableTypes[arguments[0].Type]; !validType {
		return nil, errors.New("the first return value must be a valid type")
	}

	if len(fun.Results.List) == 2 {
		if arguments[1].Type != "error" {
			return nil, fmt.Errorf("the function must return an error as the second return value")
		}
	}

	return arguments, nil
}

func parseArgument(field *ast.Field) ([]*Argument, error) {
	if len(field.Names) == 0 {
		return []*Argument{
			{
				Type: getFieldType(field.Type),
			},
		}, nil
	}

	arguments := make([]*Argument, 0, len(field.Names))
	argType := getFieldType(field.Type)
	for _, name := range field.Names {
		arguments = append(arguments, &Argument{
			Name: name.Name,
			Type: argType,
		})
	}
	return arguments, nil
}

func getFieldType(typeNode ast.Expr) string {
	switch t := typeNode.(type) {
	case *ast.Ident:
		return t.String()
	default:
		b := &bytes.Buffer{}
		format.Node(b, token.NewFileSet(), t)
		return b.String()
	}
}
