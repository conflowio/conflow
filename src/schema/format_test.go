// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema_test

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/internal/testhelper"
	"github.com/conflowio/conflow/src/schema"
	"github.com/conflowio/conflow/src/schema/formats"
	"github.com/conflowio/conflow/src/schema/schemafakes"
)

const testStructTypeFQTN = "github.com/conflowio/conflow/src/internal/testhelper.CustomStruct"

var _ = Describe("Format", func() {
	It("should return a default format for string", func() {
		name, f, found := schema.GetFormatForType("string")
		Expect(name).To(Equal(schema.FormatDefault))
		Expect(f).To(Equal(formats.String{}))
		Expect(found).To(BeTrue())
	})

	When("encountering an unregistered format", func() {
		It("should not return a format", func() {
			_, _, found := schema.GetFormatForType("non.existing")
			Expect(found).To(BeFalse())
		})
	})

	When("registering a new format", func() {
		var subject *schemafakes.FakeFormat

		BeforeEach(func() {
			subject = &schemafakes.FakeFormat{}
		})

		JustBeforeEach(func() {
			schema.RegisterFormat("test", subject)
		})

		AfterEach(func() {
			schema.UnregisterFormat("test")
		})

		When("it's a default type", func() {
			BeforeEach(func() {
				subject.TypeReturns(reflect.TypeOf(testhelper.CustomStruct{}), true)
			})

			It("should be returned for a given type", func() {
				name, f, found := schema.GetFormatForType(testStructTypeFQTN)
				Expect(found).To(BeTrue())
				Expect(name).To(Equal("test"))
				Expect(f).To(Equal(subject))
			})

			It("should be used for the given type", func() {
				source := `
					import "github.com/conflowio/conflow/src/internal/testhelper"
					// @block "configuration"
					type Foo struct {
						v testhelper.CustomStruct
					}
				`
				testhelper.ExpectGoStructToHaveSchema(source, &schema.Object{
					Name: "Foo",
					Metadata: schema.Metadata{
						ID: "test.Foo",
						Annotations: map[string]string{
							conflow.AnnotationType: conflow.BlockTypeConfiguration,
						},
					},
					Parameters: map[string]schema.Schema{
						"v": &schema.String{
							Format: "test",
						},
					},
				})
			})
		})

		When("it's not a default type", func() {
			BeforeEach(func() {
				subject.TypeReturns(reflect.TypeOf(testhelper.CustomStruct{}), false)
			})

			It("should not be returned for a given type", func() {
				_, _, found := schema.GetFormatForType(testStructTypeFQTN)
				Expect(found).To(BeFalse())
			})
		})
	})

	When("registering multiple default formats for the same type", func() {
		var f1, f2 *schemafakes.FakeFormat
		BeforeEach(func() {
			f1 = &schemafakes.FakeFormat{}
			f1.TypeReturns(reflect.TypeOf(testhelper.CustomStruct{}), true)
			schema.RegisterFormat("test1", f1)

			f2 = &schemafakes.FakeFormat{}
			f2.TypeReturns(reflect.TypeOf(testhelper.CustomStruct{}), true)
			schema.RegisterFormat("test2", f2)
		})

		It("should return the last registered format for the type", func() {
			name, f, found := schema.GetFormatForType(testStructTypeFQTN)
			Expect(found).To(BeTrue())
			Expect(name).To(Equal("test2"))
			Expect(f).To(Equal(f2))
		})

		When("f2 is unregistered", func() {
			It("should return f1", func() {
				schema.UnregisterFormat("test2")

				name, f, found := schema.GetFormatForType(testStructTypeFQTN)
				Expect(found).To(BeTrue())
				Expect(name).To(Equal("test1"))
				Expect(f).To(Equal(f1))
			})
		})

		When("both are unregistered", func() {
			It("should return no format", func() {
				schema.UnregisterFormat("test2")
				schema.UnregisterFormat("test1")

				_, _, found := schema.GetFormatForType(testStructTypeFQTN)
				Expect(found).To(BeFalse())
			})
		})

		AfterEach(func() {
			schema.UnregisterFormat("test1")
			schema.UnregisterFormat("test2")
		})
	})

})
