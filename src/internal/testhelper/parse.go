// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package testhelper

import (
	goast "go/ast"
	goparser "go/parser"
	gotoken "go/token"
	"regexp"

	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/src/conflow/block/generator"
	"github.com/conflowio/conflow/src/conflow/generator/parser"
	"github.com/conflowio/conflow/src/schema"
)

func ExpectGoStructToHaveSchema(source string, expectedSchema schema.Schema, fs ...func(schema2 schema.Schema)) {
	matches := regexp.MustCompile(`type\s+([a-zA-Z0-9_]*)\s+struct`).FindStringSubmatch(source)
	structName := matches[1]

	for _, f := range fs {
		f(expectedSchema)
	}

	source = "package test \n" + source

	fset := gotoken.NewFileSet()
	file, err := goparser.ParseFile(fset, "testfile", source, goparser.AllErrors+goparser.ParseComments)
	Expect(err).ToNot(HaveOccurred())

	parseCtx := &parser.Context{
		FileSet: fset,
		File:    file,
	}

	metadata, err := parser.ParseMetadataFromComments(file.Decls[(len(file.Decls))-1].(*goast.GenDecl).Doc.List)
	Expect(err).ToNot(HaveOccurred())

	str := file.Decls[(len(file.Decls))-1].(*goast.GenDecl).Specs[0].(*goast.TypeSpec).Type.(*goast.StructType)
	resultStruct, parseErr := generator.ParseStruct(parseCtx, str, "test", structName, metadata)

	Expect(parseErr).ToNot(HaveOccurred())
	Expect(resultStruct.Schema).To(Equal(expectedSchema))
}
