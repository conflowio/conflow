// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generator_test

import (
	goast "go/ast"
	goparser "go/parser"
	gotoken "go/token"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/pkg/conflow/function/generator"
	"github.com/conflowio/conflow/pkg/conflow/generator/parser"
	"github.com/conflowio/conflow/pkg/schema"
)

var _ = Describe("ParseArguments", func() {

	var source string
	var functionResult *generator.Function
	var parseErr error

	JustBeforeEach(func() {
		fset := gotoken.NewFileSet()

		file, err := goparser.ParseFile(fset, "testfile", source, goparser.AllErrors+goparser.ParseComments)
		Expect(err).ToNot(HaveOccurred())

		parseCtx := &parser.Context{
			FileSet: fset,
			File:    file,
		}

		var comments []*goast.Comment
		if file.Decls[0].(*goast.FuncDecl).Doc != nil {
			comments = file.Decls[0].(*goast.FuncDecl).Doc.List
		}

		metadata, err := parser.ParseMetadataFromComments(comments)
		Expect(err).ToNot(HaveOccurred())

		fun := file.Decls[0].(*goast.FuncDecl).Type
		functionResult, parseErr = generator.ParseFunction(parseCtx, fun, "test", "Foo", metadata)
	})

	Context("when the function has no arguments", func() {
		BeforeEach(func() {
			source = `
				package foo
				func Foo() int64 {
					return 0
				}
			`
		})

		It("should return with empty argument list", func() {
			Expect(parseErr).ToNot(HaveOccurred())
			Expect(functionResult.Schema).To(Equal(&schema.Function{
				Metadata: schema.Metadata{
					ID: "test.Foo",
				},
				Result: &schema.Integer{},
			}))
			Expect(functionResult.ReturnsError).To(BeFalse())
		})
	})

	Context("when the function has a docblock arguments", func() {
		BeforeEach(func() {
			source = `
				package foo
				// Foo is a test function
				func Foo() int64 {
					return 0
				}
			`
		})

		It("should set the description", func() {
			Expect(parseErr).ToNot(HaveOccurred())
			Expect(functionResult.Schema).To(Equal(&schema.Function{
				Metadata: schema.Metadata{
					ID:          "test.Foo",
					Description: "It is a test function",
				},
				Result: &schema.Integer{},
			}))
		})
	})

	Context("when it returns an error", func() {
		BeforeEach(func() {
			source = `
				package foo
				func Foo() (int64, error) {
					return 0, nil
				}
			`
		})

		It("should set ReturnsError to true", func() {
			Expect(parseErr).ToNot(HaveOccurred())
			Expect(functionResult.Schema).To(Equal(&schema.Function{
				Metadata: schema.Metadata{
					ID: "test.Foo",
				},
				Result: &schema.Integer{},
			}))
			Expect(functionResult.ReturnsError).To(BeTrue())
		})
	})

	Context("when it has parameters", func() {
		BeforeEach(func() {
			source = `
				package foo
				func Foo(a int64, b float64) int64 {
					return 0
				}
			`
		})

		It("should parse the arguments", func() {
			Expect(parseErr).ToNot(HaveOccurred())
			Expect(functionResult.Schema).To(Equal(&schema.Function{
				Metadata: schema.Metadata{
					ID: "test.Foo",
				},
				Parameters: []schema.NamedSchema{
					{Name: "a", Schema: &schema.Integer{}},
					{Name: "b", Schema: &schema.Number{}},
				},
				Result: &schema.Integer{},
			}))
		})
	})

	Context("when a parameter is marked as result type", func() {
		BeforeEach(func() {
			source = `
				package foo
				func Foo(
					// @result_type
					a interface{},
				) interface{} {
					return a
				}
			`
		})

		It("should parse the arguments", func() {
			Expect(parseErr).ToNot(HaveOccurred())
			Expect(functionResult.Schema).To(Equal(&schema.Function{
				Metadata: schema.Metadata{
					ID: "test.Foo",
				},
				Parameters: []schema.NamedSchema{
					{Name: "a", Schema: &schema.Any{}},
				},
				Result:         &schema.Any{},
				ResultTypeFrom: "a",
			}))
		})
	})

	Context("when an interface{} parameter has types set", func() {
		BeforeEach(func() {
			source = `
				package foo
				func Foo(
					// @types ["integer", "number"]
					a interface{},
				) interface{} {
					return a
				}
			`
		})

		It("should parse the arguments", func() {
			Expect(parseErr).ToNot(HaveOccurred())
			Expect(functionResult.Schema).To(Equal(&schema.Function{
				Metadata: schema.Metadata{
					ID: "test.Foo",
				},
				Parameters: []schema.NamedSchema{
					{
						Name: "a",
						Schema: &schema.Any{
							Types: []string{"integer", "number"},
						},
					},
				},
				Result: &schema.Any{},
			}))
		})
	})

	Context("when it has a variadic argument", func() {
		BeforeEach(func() {
			source = `
				package foo
				func Foo(a1 int64, a2 int64, rest ...int64) int64 {
					return a1
				}
			`
		})

		It("should parse the arguments", func() {
			Expect(parseErr).ToNot(HaveOccurred())
			Expect(functionResult.Schema).To(Equal(&schema.Function{
				Metadata: schema.Metadata{
					ID: "test.Foo",
				},
				Parameters: []schema.NamedSchema{
					{
						Name:   "a1",
						Schema: &schema.Integer{},
					},
					{
						Name:   "a2",
						Schema: &schema.Integer{},
					},
				},
				AdditionalParameters: &schema.NamedSchema{
					Name:   "rest",
					Schema: &schema.Integer{},
				},
				Result: &schema.Integer{},
			}))
		})
	})

	Context("when the function has no return value", func() {
		BeforeEach(func() {
			source = `
				package foo
				func Foo() {
				}
			`
		})

		It("should return an error", func() {
			Expect(parseErr).To(MatchError("the function must return with a single value, or a single value and an error"))
		})
	})

	Context("when the function has an invalid parameter type", func() {
		BeforeEach(func() {
			source = `
				package foo
				func Foo(val int8) int64 {
					return 0
				}
			`
		})

		It("should return an error", func() {
			Expect(parseErr).To(MatchError("parameter val is invalid: type int8 is not allowed"))
		})
	})

	Context("when the return type is invalid", func() {
		BeforeEach(func() {
			source = `
				package foo
				func Foo(val int64) int8 {
					return 0
				}
			`
		})

		It("should return an error", func() {
			Expect(parseErr).To(MatchError("result value is invalid: type int8 is not allowed"))
		})
	})

})
