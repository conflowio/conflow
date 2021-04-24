// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generator

import (
	"bufio"
	"bytes"
	"fmt"
	goast "go/ast"
	"go/build"
	goparser "go/parser"
	gotoken "go/token"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/opsidian/basil/basil/generator/parser"

	blockgenerator "github.com/opsidian/basil/basil/block/generator"
	functiongenerator "github.com/opsidian/basil/basil/function/generator"
)

// Generate generates code for the given types
func Generate(dir string) error {
	parseCtx := &parser.Context{
		WorkDir: dir,
		FileSet: &gotoken.FileSet{},
	}

	err := filepath.WalkDir(dir, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(d.Name(), ".go") {
			if err := processFile(parseCtx, filePath); err != nil {
				return fmt.Errorf("failed to process %s: %w", filePath, err)
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func processFile(parseCtx *parser.Context, filePath string) error {
	src, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("generate: %s", err)
	}

	input := bufio.NewReader(bytes.NewReader(src))
	lineNum := 0
	var hasDirective bool
	for {
		lineNum++
		var buf []byte
		buf, err = input.ReadSlice('\n')
		if err == bufio.ErrBufferFull {
			for err == bufio.ErrBufferFull {
				_, err = input.ReadSlice('\n')
			}
			if err != nil {
				break
			}
			continue
		}

		if err != nil {
			break
		}

		if bytes.HasPrefix(buf, []byte("// @block")) ||
			bytes.HasPrefix(buf, []byte("//\t@block")) ||
			bytes.HasPrefix(buf, []byte("// @function")) ||
			bytes.HasPrefix(buf, []byte("//\t@function")) {
			hasDirective = true
			break
		}
	}
	if err != nil && err != io.EOF {
		return err
	}

	if !hasDirective {
		return nil
	}

	astFile, err := goparser.ParseFile(parseCtx.FileSet, filePath, src, goparser.AllErrors+goparser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	buildContext := build.Default
	pkg, err := buildContext.Import(".", path.Dir(filePath), 0)
	if err != nil {
		return fmt.Errorf("failed to determine package import path for %s: %w", astFile.Name.Name, err)
	}

	fileParseCtx := parseCtx.WithFile(astFile)

	for _, d := range astFile.Decls {
		switch dt := d.(type) {
		case *goast.GenDecl:
			hasDirective := false
			if dt.Doc != nil {
				for _, c := range dt.Doc.List {
					if strings.Contains(c.Text, "@block") {
						hasDirective = true
						break
					}
				}
			}
			if !hasDirective {
				continue
			}
			name, structType, ok := getStructType(dt)
			if !ok {
				return fmt.Errorf("@block annotation should only be set on a struct type")
			}

			res, data, err := blockgenerator.GenerateInterpreter(fileParseCtx, structType, pkg.ImportPath, name, dt.Doc.List)
			if err != nil {
				return fmt.Errorf("failed to generate interpreter for struct %s: %w", name, err)
			}

			if err := writeFile(path.Join(path.Dir(filePath), data.InterpreterPath), name, res); err != nil {
				return fmt.Errorf("failed to generate interpreter for struct %s: %w", name, err)
			}
		case *goast.FuncDecl:
			hasDirective := false
			if dt.Doc != nil {
				for _, c := range dt.Doc.List {
					if strings.Contains(c.Text, "@function") {
						hasDirective = true
						break
					}
				}
			}
			if !hasDirective {
				continue
			}
			funcType := dt.Type
			name := dt.Name.Name

			res, data, err := functiongenerator.GenerateInterpreter(fileParseCtx, funcType, pkg.ImportPath, name, dt.Doc.List)
			if err != nil {
				return fmt.Errorf("failed to generate interpreter for function %s: %w", name, err)
			}

			if err := writeFile(path.Join(path.Dir(filePath), data.InterpreterPath), name, res); err != nil {
				return fmt.Errorf("failed to generate interpreter for function %s, %w", name, err)
			}
		}
	}

	return nil
}

func getStructType(d *goast.GenDecl) (string, *goast.StructType, bool) {
	if len(d.Specs) == 0 {
		return "", nil, false
	}
	typeSpec, ok := d.Specs[0].(*goast.TypeSpec)
	if !ok {
		return "", nil, false
	}

	structType, ok := typeSpec.Type.(*goast.StructType)
	if !ok {
		return "", nil, false
	}

	return typeSpec.Name.Name, structType, true
}
