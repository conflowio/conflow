// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path"
	"regexp"
)

type errStructNotFound error

func FindStruct(parseCtx *Context, pkgName, name string) (*ast.StructType, *Metadata, error) {
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
			return parseStruct(parseCtx, pkgName, filePath, name)
		}
	}

	return nil, nil, errStructNotFound(fmt.Errorf("type %s not found in package %s", name, pkgName))
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

func parseStruct(parseCtx *Context, pkgName, filePath, name string) (*ast.StructType, *Metadata, error) {
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

			metadata, err := ParseMetadataFromComments(name, comments)
			if err != nil {
				return nil, nil, err
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
				str, _, err := FindStruct(parseCtx, pkgName, t.Name)
				if err != nil {
					if _, notFound := err.(errStructNotFound); notFound {
						for _, imp := range astFile.Imports {
							if imp.Name != nil && imp.Name.Name == "." {
								str, _, err := FindStruct(parseCtx, imp.Path.Value[1:len(imp.Path.Value)-1], t.Name)
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

				str, _, err := FindStruct(parseCtx, pkgName, t.Sel.Name)
				if err != nil {
					return nil, nil, err
				}

				return str, metadata, nil
			default:
				return nil, nil, fmt.Errorf("unexpected type: %T", t)
			}
		}
	}

	return nil, nil, errStructNotFound(errors.New("struct not found"))
}
