// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generator

import (
	"fmt"
	"go/ast"
	"path"
	"regexp"
	"strings"

	blockgenerator "github.com/opsidian/basil/basil/block/generator"
	functiongenerator "github.com/opsidian/basil/basil/function/generator"
)

// Generate generates code for the given types
func Generate(dir string, packageName string, file string, line int) error {
	astFile, astNode, name, err := getType(file, line)
	if err != nil {
		return err
	}

	filename := regexp.MustCompile("[A-Z][a-z0-9_]+").ReplaceAllStringFunc(name, func(str string) string {
		return "_" + strings.ToLower(str)
	})
	filename = strings.ToLower(strings.TrimLeft(filename, "_")) + ".basil.go"
	filePath := path.Join(dir, filename)

	var res []byte
	var genErr error

	switch t := astNode.(type) {
	case *ast.StructType:
		res, genErr = blockgenerator.GenerateInterpreter(t, astFile, packageName, name)
	case *ast.FuncType:
		res, genErr = functiongenerator.GenerateInterpreter(t, astFile, packageName, name)
	default:
		return fmt.Errorf("type is not supported: %T", astNode)
	}

	if genErr != nil {
		return genErr
	}

	return writeFile(name, filename, filePath, res)
}
