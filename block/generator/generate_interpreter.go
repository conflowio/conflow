package generator

import (
	"bytes"
	"go/ast"
	"strings"
	"text/template"

	"github.com/opsidian/basil/util"
	"github.com/opsidian/basil/variable"
)

// GenerateInterpreter generates an interpreter for the given block
func GenerateInterpreter(str *ast.StructType, file *ast.File, pkgName string, name string) ([]byte, error) {
	params, err := generateTemplateParams(str, file, pkgName, name)
	if err != nil {
		return nil, err
	}

	tmpl := template.New("block_interpreter")
	tmpl.Funcs(map[string]interface{}{
		"trimPrefix": func(s string, prefix string) string {
			return strings.TrimPrefix(s, prefix)
		},
	})
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

func generateTemplateParams(str *ast.StructType, file *ast.File, pkgName string, name string) (*InterpreterTemplateParams, error) {
	fields, err := ParseFields(str, file)
	if err != nil {
		return nil, err
	}
	var idField, valueField *Field
	var hasForeignID bool

	stages := []string{}
	paramTypes := map[string]string{}
	blockTypes := map[string]string{}
	nodeTypes := map[string]string{}
	evalFieldsCnt := 0
	requiredFieldsCnt := 0
	for _, field := range fields {
		if field.Stage != "-" {
			if !util.StringSliceContains(stages, field.Stage) {
				stages = append(stages, field.Stage)
			}
			if !field.IsID && !field.IsBlock {
				evalFieldsCnt++
			}
		}
		if field.Required {
			requiredFieldsCnt++
		}
		if field.IsValue {
			valueField = field
		}

		switch {
		case field.IsID:
			idField = field
			hasForeignID = field.IsReference
		case field.IsParam:
			paramTypes[field.ParamName] = field.Type
		case field.IsNode:
			nodeTypes[field.Type] = field.Name
		case field.IsBlock:
			blockTypes[field.Type] = field.Name
		}
	}

	return &InterpreterTemplateParams{
		Package:                pkgName,
		Name:                   name,
		Stages:                 stages,
		ParamTypes:             paramTypes,
		BlockTypes:             blockTypes,
		NodeTypes:              nodeTypes,
		Fields:                 fields,
		IDField:                idField,
		ValueField:             valueField,
		HasForeignID:           hasForeignID,
		NodeValueFunctionNames: variable.NodeValueFunctionNames,
		EvalFieldsCnt:          evalFieldsCnt,
		RequiredFieldsCnt:      requiredFieldsCnt,
	}, nil
}
