// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generator_test

import (
	"fmt"
	goast "go/ast"
	goparser "go/parser"
	gotoken "go/token"

	"github.com/opsidian/conflow/basil/generator/parser"

	"github.com/opsidian/conflow/basil/block/generator"
	"github.com/opsidian/conflow/basil/schema"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("ParseStruct", func() {

	var expectedIDAnnotations = map[string]string{
		"id": "true",
	}

	basicTemplate := func(fields string) string {
		return fmt.Sprintf(`
			package foo
			import (
				"time"
				"github.com/opsidian/conflow/basil"
			)
			// Foo is a test struct
			type Foo struct {
				// @id
				id basil.ID
				%s
			}
		`, fields)
	}

	DescribeTable("a struct with valid fields",
		func(fields string, f func(schema.Schema)) {
			source := basicTemplate(fields)

			expectedSchema := &schema.Object{
				Metadata: schema.Metadata{
					Description: "It is a test struct",
				},
				Name: "Foo",
				Properties: map[string]schema.Schema{
					"id": &schema.String{
						Metadata: schema.Metadata{
							Annotations: expectedIDAnnotations,
							ReadOnly:    true,
						},
						Format: "basil.ID",
					},
				},
			}
			f(expectedSchema)

			fset := gotoken.NewFileSet()
			file, err := goparser.ParseFile(fset, "testfile", source, goparser.AllErrors+goparser.ParseComments)
			Expect(err).ToNot(HaveOccurred())

			parseCtx := &parser.Context{
				FileSet: fset,
				File:    file,
			}

			metadata, err := parser.ParseMetadataFromComments("Foo", file.Decls[(len(file.Decls))-1].(*goast.GenDecl).Doc.List)
			Expect(err).ToNot(HaveOccurred())

			str := file.Decls[(len(file.Decls))-1].(*goast.GenDecl).Specs[0].(*goast.TypeSpec).Type.(*goast.StructType)
			resultStruct, parseErr := generator.ParseStruct(parseCtx, str, "test", "Foo", metadata)

			Expect(parseErr).ToNot(HaveOccurred())
			Expect(resultStruct.Schema).To(Equal(expectedSchema))
		},
		Entry("valid id field", "", func(schema.Schema) {}),

		Entry("string field", "field string", func(s schema.Schema) {
			s.(*schema.Object).Properties["field"] = &schema.String{}
		}),

		Entry("bool field", "field bool", func(s schema.Schema) {
			s.(*schema.Object).Properties["field"] = &schema.Boolean{}
		}),

		Entry("integer field", "field int64", func(s schema.Schema) {
			s.(*schema.Object).Properties["field"] = &schema.Integer{}
		}),

		Entry("number field", "field float64", func(s schema.Schema) {
			s.(*schema.Object).Properties["field"] = &schema.Number{}
		}),

		Entry("time duration field", "field time.Duration", func(s schema.Schema) {
			s.(*schema.Object).Properties["field"] = &schema.TimeDuration{}
		}),

		Entry("string array", "field []string", func(s schema.Schema) {
			s.(*schema.Object).Properties["field"] = &schema.Array{
				Items: &schema.String{},
			}
		}),

		Entry("bool array", "field []bool", func(s schema.Schema) {
			s.(*schema.Object).Properties["field"] = &schema.Array{
				Items: &schema.Boolean{},
			}
		}),

		Entry("integer array", "field []int64", func(s schema.Schema) {
			s.(*schema.Object).Properties["field"] = &schema.Array{
				Items: &schema.Integer{},
			}
		}),

		Entry("number array", "field []float64", func(s schema.Schema) {
			s.(*schema.Object).Properties["field"] = &schema.Array{
				Items: &schema.Number{},
			}
		}),

		Entry("time duration array", "field []time.Duration", func(s schema.Schema) {
			s.(*schema.Object).Properties["field"] = &schema.Array{
				Items: &schema.TimeDuration{},
			}
		}),

		Entry("arrays of arrays", "field [][]string", func(s schema.Schema) {
			s.(*schema.Object).Properties["field"] = &schema.Array{
				Items: &schema.Array{
					Items: &schema.String{},
				},
			}
		}),

		Entry("string map", "field map[string]string", func(s schema.Schema) {
			s.(*schema.Object).Properties["field"] = &schema.Map{
				AdditionalProperties: &schema.String{},
			}
		}),

		Entry("integer map", "field map[string]int64", func(s schema.Schema) {
			s.(*schema.Object).Properties["field"] = &schema.Map{
				AdditionalProperties: &schema.Integer{},
			}
		}),

		Entry("number map", "field map[string]float64", func(s schema.Schema) {
			s.(*schema.Object).Properties["field"] = &schema.Map{
				AdditionalProperties: &schema.Number{},
			}
		}),

		Entry("boolean map", "field map[string]bool", func(s schema.Schema) {
			s.(*schema.Object).Properties["field"] = &schema.Map{
				AdditionalProperties: &schema.Boolean{},
			}
		}),

		Entry("time duration map", "field map[string]time.Duration", func(s schema.Schema) {
			s.(*schema.Object).Properties["field"] = &schema.Map{
				AdditionalProperties: &schema.TimeDuration{},
			}
		}),

		Entry("maps of maps", "field map[string]map[string]string", func(s schema.Schema) {
			s.(*schema.Object).Properties["field"] = &schema.Map{
				AdditionalProperties: &schema.Map{
					AdditionalProperties: &schema.String{},
				},
			}
		}),

		Entry(
			"valid json field name should be used as property name",
			"field string `json:\"custom_field_name\"`",
			func(s schema.Schema) {
				s.(*schema.Object).Properties["custom_field_name"] = &schema.String{}
				s.(*schema.Object).PropertyNames = map[string]string{
					"custom_field_name": "field",
				}
			},
		),

		Entry(
			"valid @name directive should be used as property name",
			"// @name \"custom_field_name\"\nfield string",
			func(s schema.Schema) {
				s.(*schema.Object).Properties["custom_field_name"] = &schema.String{}
				s.(*schema.Object).PropertyNames = map[string]string{
					"custom_field_name": "field",
				}
			},
		),

		Entry(
			"invalid json field name should not be used as property name",
			"field string `json:\"customFieldName\"`",
			func(s schema.Schema) {
				s.(*schema.Object).Properties["field"] = &schema.String{}
			},
		),

		Entry(
			"valid property name should be generated",
			"fieldName string",
			func(s schema.Schema) {
				s.(*schema.Object).Properties["field_name"] = &schema.String{}
				s.(*schema.Object).PropertyNames = map[string]string{
					"field_name": "fieldName",
				}
			},
		),

		Entry(
			"field should be ignored if a JSON annotation ignores it",
			"fieldName string `json:\"-\"`",
			func(s schema.Schema) {
			},
		),

		Entry(
			"JSON ignore should work on an unsupported field type",
			"fieldName int8 `json:\"-\"`",
			func(s schema.Schema) {
			},
		),

		Entry(
			"JSON ignore should work on an unsupported array field type",
			"fieldName []int8 `json:\"-\"`",
			func(s schema.Schema) {
			},
		),

		Entry(
			"JSON ignore should work on an unsupported map field type",
			"fieldName map[int8]int8 `json:\"-\"`",
			func(s schema.Schema) {
			},
		),

		Entry(
			"field should be ignored if there is an @ignore annotation",
			"// @ignore\nfieldName string",
			func(s schema.Schema) {
			},
		),

		Entry(
			"@ignore should work on an unsupported field type",
			"// @ignore\nfieldName int8",
			func(s schema.Schema) {
			},
		),

		Entry(
			"@ignore should work on an unsupported array field type",
			"// @ignore\nfieldName []int8",
			func(s schema.Schema) {
			},
		),

		Entry(
			"@ignore should work on an unsupported map field type",
			"// @ignore\nfieldName map[int8]int8",
			func(s schema.Schema) {
			},
		),
	)

	Context("special cases", func() {
		var source string
		var resultStruct *generator.Struct
		var parseErr error

		JustBeforeEach(func() {
			fset := gotoken.NewFileSet()
			file, err := goparser.ParseFile(fset, "testfile", source, goparser.AllErrors+goparser.ParseComments)
			Expect(err).ToNot(HaveOccurred())

			parseCtx := &parser.Context{
				FileSet: fset,
				File:    file,
			}

			str := file.Decls[(len(file.Decls))-1].(*goast.GenDecl).Specs[0].(*goast.TypeSpec).Type.(*goast.StructType)
			resultStruct, parseErr = generator.ParseStruct(parseCtx, str, "test", "Foo", &parser.Metadata{})
		})

		Context("when the basil package has an alias", func() {
			BeforeEach(func() {
				source = `
				package foo
				import (
					basilalias "github.com/opsidian/conflow/basil"
				)
				type Foo struct {
					// @id
					id basilalias.ID
				}`
			})

			It("should return with the parsed fields", func() {
				Expect(parseErr).ToNot(HaveOccurred())
				Expect(resultStruct.Schema).To(Equal(&schema.Object{
					Name: "Foo",
					Properties: map[string]schema.Schema{
						"id": &schema.String{
							Metadata: schema.Metadata{
								Annotations: expectedIDAnnotations,
								ReadOnly:    true,
							},
							Format: "basil.ID",
						},
					},
				}))
			})
		})

		Context("when a non basil.ID field is marked as id", func() {
			BeforeEach(func() {
				source = `
				package foo
				type Foo struct {
					// @id
					foo string
				}`
			})

			It("should return with error", func() {
				Expect(parseErr).To(MatchError("failed to parse field \"foo\": id annotation can only be set on a basil.ID field"))
			})
		})

		Context("when there are multiple id fields", func() {
			BeforeEach(func() {
				source = `
				package foo
				import (
					"github.com/opsidian/conflow/basil"
				)
				type Foo struct {
					// @id
					id1 basil.ID
					// @id
					id2 basil.ID
				}`
			})

			It("should return with error", func() {
				Expect(parseErr).To(MatchError("multiple id fields were found: id1, id2"))
			})
		})

		Context("when there are multiple value fields", func() {
			BeforeEach(func() {
				source = `
				package foo
				import (
					"github.com/opsidian/conflow/basil"
				)
				type Foo struct {
					// @id
					id basil.ID
					// @value
					value1 string
					// @value
					value2 string
				}`
			})

			It("should return with error", func() {
				Expect(parseErr).To(MatchError("multiple value fields were found: value1, value2"))
			})
		})

		Context("when there are required fields other than the value field", func() {
			BeforeEach(func() {
				source = `
				package foo
				import (
					"github.com/opsidian/conflow/basil"
				)
				type Foo struct {
					// @id
					id basil.ID
					// @value
					value string
					// @required
					foo string
				}`
			})

			It("should return with error", func() {
				Expect(parseErr).To(MatchError("when setting a value field then no other fields can be required"))
			})
		})

		Context("when there is an unknown directive", func() {
			BeforeEach(func() {
				source = `
				package foo
				import (
					"github.com/opsidian/conflow/basil"
				)
				type Foo struct {
					// @id
					// @nonexisting
					id basil.ID
				}`
			})

			It("should return with error", func() {
				Expect(parseErr).To(MatchError("failed to parse field \"id\": @nonexisting directive is unknown or not allowed at 2:1"))
			})
		})

	})

})
