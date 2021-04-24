// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generator

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"

	"golang.org/x/tools/go/packages"
)

func getType(filename string, line int) (*packages.Package, ast.Node, string, error) {
	fs := token.NewFileSet()

	pkgConf := &packages.Config{
		Mode: packages.NeedImports + packages.NeedDeps + packages.NeedTypes + packages.NeedTypesInfo,
		Fset: fs,
	}
	pkgs, err := packages.Load(pkgConf, fmt.Sprintf("file=%s", filename))
	if err != nil {
		log.Fatal(err)
	}
	pkg := pkgs[0]

	var name string
	for e := range pkg.TypesInfo.Defs {
		if fs.File(e.Pos()).Name() == filename && fs.Position(e.Pos()).Line == line+1 {
			name = e.Name
			break
		}
	}

	for e := range pkg.TypesInfo.Types {
		if fs.File(e.Pos()).Name() == filename && fs.Position(e.Pos()).Line == line+1 {
			switch n := e.(type) {
			case *ast.StructType:
				return pkg, n, name, nil
			case *ast.FuncType:
				return pkg, n, name, nil
			}
		}
	}

	return nil, nil, "", fmt.Errorf("no function or struct found, you must place go generate on the previous line")
}
