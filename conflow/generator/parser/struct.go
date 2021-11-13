// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
)

type errTypeNotFound struct {
	msg string
}

func (e errTypeNotFound) Error() string {
	return e.msg
}

func FindType(parseCtx *Context, pkgName, name string) (e ast.Expr, m *Metadata, err error) {
	buildContext := build.Default
	pkg, err := buildContext.Import(pkgName, parseCtx.WorkDir, 0)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find package %s: %w", pkgName, err)
	}

	typeRegex := regexp.MustCompile(`^\s*type\s+` + name + `\s+`)

	for _, goFile := range pkg.GoFiles {
		filePath := path.Join(pkg.Dir, goFile)
		hasType, err := fileHasLine(filePath, typeRegex)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
		}

		if hasType {
			return parseType(parseCtx, pkgName, filePath, name)
		}
	}

	return nil, nil, errTypeNotFound{msg: fmt.Sprintf("type %q not found in package %q", name, pkgName)}
}

func fileHasLine(filePath string, re *regexp.Regexp) (bool, error) {
	src, err := os.ReadFile(filePath)
	if err != nil {
		return false, err
	}

	input := bufio.NewReader(bytes.NewReader(src))
	for {
		var buf []byte
		buf, err = input.ReadSlice('\n')
		if err == bufio.ErrBufferFull {
			if re.Match(buf) {
				return true, nil
			}

			for err == bufio.ErrBufferFull {
				_, err = input.ReadSlice('\n')
			}
			if err != nil {
				break
			}

			continue
		}

		if (err == nil || err == io.EOF) && re.Match(buf) {
			return true, nil
		}

		if err != nil {
			break
		}
	}

	if err != io.EOF {
		return false, err
	}

	return false, nil
}

func parseType(parseCtx *Context, pkgName, filePath, name string) (ast.Expr, *Metadata, error) {
	astFile, err := parser.ParseFile(&token.FileSet{}, filePath, nil, parser.AllErrors+parser.ParseComments)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse file %s: %w", filePath, err)
	}

	for _, d := range astFile.Decls {
		genDecl, ok := d.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			if typeSpec.Name.Name != name {
				continue
			}

			var comments []*ast.Comment
			if genDecl.Doc != nil {
				comments = append(comments, genDecl.Doc.List...)
			}

			metadata, err := ParseMetadataFromComments(comments)
			if err != nil {
				return nil, nil, err
			}

			if strings.HasPrefix(metadata.Description, name+" ") {
				metadata.Description = strings.Replace(metadata.Description, name+" ", "It ", 1)
			}

			switch t := typeSpec.Type.(type) {
			case *ast.StructType:
				return t, metadata, nil
			case *ast.InterfaceType:
				return &ast.StructType{
					Struct: t.Pos(),
					Fields: &ast.FieldList{},
				}, metadata, nil
			case *ast.Ident:
				switch t.Name {
				case "bool", "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "float32",
					"float64", "complex64", "complex128", "string", "int", "uint", "uintptr", "byte", "rune":
					return t, metadata, nil
				}
				str, _, err := FindType(parseCtx, pkgName, t.Name)
				if err != nil {
					if _, notFound := err.(errTypeNotFound); notFound {
						for _, imp := range astFile.Imports {
							if imp.Name != nil && imp.Name.Name == "." {
								str, _, err := FindType(parseCtx, imp.Path.Value[1:len(imp.Path.Value)-1], t.Name)
								if err == nil {
									return str, metadata, nil
								}
							}
						}
					}
					return nil, nil, err
				}
				return str, metadata, nil
			case *ast.SelectorExpr:
				pkgName := GetImportPath(astFile, t.X.(*ast.Ident).Name)
				if pkgName == "" {
					return nil, nil, fmt.Errorf("package alias %s could not be resolved", t.X.(*ast.Ident).Name)
				}

				str, _, err := FindType(parseCtx, pkgName, t.Sel.Name)
				if err != nil {
					return nil, nil, err
				}

				return str, metadata, nil
			default:
				return t, metadata, nil
			}
		}
	}

	return nil, nil, errTypeNotFound{msg: fmt.Sprintf("type %q not found in package %q", name, pkgName)}
}
