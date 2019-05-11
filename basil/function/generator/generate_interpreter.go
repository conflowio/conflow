package generator

import (
	"bytes"
	"go/ast"
	"strings"
	"text/template"

	"github.com/opsidian/basil/basil/variable"
)

// GenerateInterpreter generates an interpreter for the given function
func GenerateInterpreter(fun *ast.FuncType, file *ast.File, pkgName string, name string) ([]byte, error) {
	params, err := generateTemplateParams(fun, file, pkgName, name)
	if err != nil {
		return nil, err
	}

	tmpl := template.New("block_interpreter")
	if _, parseErr := tmpl.Parse(interpreterTemplate); parseErr != nil {
		return nil, parseErr
	}

	res := &bytes.Buffer{}
	err = tmpl.Execute(res, params)
	if err != nil {
		return nil, err
	}

	return res.Bytes(), nil
}

func generateTemplateParams(fun *ast.FuncType, file *ast.File, pkgName string, name string) (*InterpreterTemplateParams, error) {
	arguments, err := ParseArguments(fun, file)
	if err != nil {
		return nil, err
	}

	results, err := ParseResults(fun, file)
	if err != nil {
		return nil, err
	}

	returnNodeType := false

	if len(arguments) > 0 {
		// If the first argument is a union type and the first return value is the same type
		// then we will return with the node's type
		for unionType := range variable.UnionTypes {
			if arguments[0].Type == unionType {
				if results[0].Type == unionType {
					returnNodeType = true
				}
				break
			}
		}
	}

	return &InterpreterTemplateParams{
		Package:                pkgName,
		Name:                   strings.ToUpper(string(name[0])) + name[1:],
		FuncName:               name,
		Arguments:              arguments,
		Results:                results,
		ResultType:             results[0].Type,
		ReturnNodeType:         returnNodeType,
		NodeValueFunctionNames: variable.NodeValueFunctionNames,
	}, nil
}
