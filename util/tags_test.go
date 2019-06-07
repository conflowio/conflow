// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package util_test

import (
	"sort"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/opsidian/basil/util"
)

var _ = Describe("GetTags", func() {

	DescribeTable("it parser the valid tags correctly",
		func(obj interface{}, expected util.StructTags) {
			tags := util.GetTags(obj, "basil")
			Expect(tags).To(Equal(expected))
		},
		Entry("no fields", struct{}{}, nil),
		Entry("one field, no tags", struct{ f string }{}, util.StructTags{"f": nil}),
		Entry("key-value",
			struct {
				f string `basil:"a=b"`
			}{},
			util.StructTags{
				"f": map[string]string{
					"a": "b",
				},
			},
		),
		Entry("multiple key-values",
			struct {
				f string `basil:"a=b,c=d"`
			}{},
			util.StructTags{
				"f": map[string]string{
					"a": "b",
					"c": "d",
				},
			},
		),
		Entry("multiple fields",
			struct {
				f1 string `basil:"a=b"`
				f2 string `basil:"c=d"`
			}{},
			util.StructTags{
				"f1": map[string]string{
					"a": "b",
				},
				"f2": map[string]string{
					"c": "d",
				},
			},
		),
		Entry("different tag",
			struct {
				f string `foo:"a=b"`
			}{},
			util.StructTags{
				"f": nil,
			},
		),
		Entry("empty tag",
			struct {
				f string `basil:""`
			}{},
			util.StructTags{
				"f": nil,
			},
		),
		Entry("only whitespace in tag",
			struct {
				f string `basil:" "`
			}{},
			util.StructTags{
				"f": nil,
			},
		),
		Entry("strips whitespaces",
			struct {
				f string `basil:" a = b , c = d "`
			}{},
			util.StructTags{
				"f": map[string]string{
					"a": "b",
					"c": "d",
				},
			},
		),
		Entry("no value should set the tag with true",
			struct {
				f string `basil:"foo"`
			}{},
			util.StructTags{
				"f": map[string]string{
					"foo": "true",
				},
			},
		),
		Entry("extra comma should be ignored",
			struct {
				f string `basil:",a=b,"`
			}{},
			util.StructTags{
				"f": map[string]string{
					"a": "b",
				},
			},
		),
		Entry("pointer",
			&struct {
				f string `basil:"a=b"`
			}{},
			util.StructTags{
				"f": map[string]string{
					"a": "b",
				},
			},
		),
	)

	Describe("when input is not a struct", func() {
		It("should panic", func() {
			var obj = "this is not a struct"
			Expect(func() { util.GetTags(obj, "basil") }).To(Panic())
		})
	})

	Describe("when tag key is empty", func() {
		It("should panic", func() {
			obj := struct {
				f1 string `basil:"=value"`
			}{}
			Expect(func() { util.GetTags(obj, "basil") }).To(Panic())
		})
	})
})

var _ = Describe("FilterFieldsByTags", func() {
	var obj = struct {
		f1 string
		f2 string `basil:"key=a"`
		f3 string `basil:"key=a,key2=b"`
		f4 string `basil:"key2=a"`
		f5 string `basil:"key=b"`
		f6 string `basil:"k`
	}{}

	It("should return fields with specific key-value pairs", func() {
		tags := util.GetTags(obj, "basil")
		actual := util.FilterFieldsByTags(tags, "key", "a")
		sort.Strings(actual)
		Expect(actual).To(Equal([]string{"f2", "f3"}))
	})

	It("should return fields for all filtered key values", func() {
		tags := util.GetTags(obj, "basil")
		actual := util.FilterFieldsByTags(tags, "key", "a", "b")
		sort.Strings(actual)
		Expect(actual).To(Equal([]string{"f2", "f3", "f5"}))
	})

	It("should return nil if no key matches", func() {
		tags := util.GetTags(obj, "basil")
		actual := util.FilterFieldsByTags(tags, "x", "foo")
		sort.Strings(actual)
		Expect(actual).To(BeNil())
	})

	It("should return nil if no value matches", func() {
		tags := util.GetTags(obj, "basil")
		actual := util.FilterFieldsByTags(tags, "key", "x")
		sort.Strings(actual)
		Expect(actual).To(BeNil())
	})

	It("should panic if no value was given", func() {
		tags := util.GetTags(obj, "basil")
		Expect(func() { util.FilterFieldsByTags(tags, "key") }).To(Panic())
	})
})
