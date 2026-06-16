// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package bind_test

import (
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/pkg/conflow/bind"
	"github.com/conflowio/conflow/pkg/schema"
	"github.com/conflowio/conflow/pkg/values"
)

var _ = Describe("BindValue", func() {
	var stringSchema schema.Schema = &schema.String{}
	var integerSchema schema.Schema = &schema.Integer{}
	stringArraySchema := &schema.Array{Items: &schema.String{}}
	intMapSchema := &schema.Map{AdditionalProperties: &schema.Integer{}}
	objectSchema := &schema.Object{
		Properties: map[string]schema.Schema{
			"name":  &schema.String{},
			"count": &schema.Integer{},
			"tags":  &schema.Array{Items: &schema.String{}},
		},
	}

	DescribeTable("bind policy",
		func(s schema.Schema, input interface{}, assert func(result interface{})) {
			result, err := bind.BindValue(s, input)
			Expect(err).NotTo(HaveOccurred())
			assert(result)
		},
		Entry("[]interface{} array of string converts to immutable list",
			stringArraySchema,
			[]interface{}{"a", "b"},
			func(result interface{}) {
				list, ok := result.(*values.List[interface{}])
				Expect(ok).To(BeTrue())
				Expect(list.At(0)).To(Equal("a"))
				Expect(list.At(1)).To(Equal("b"))
			},
		),
		Entry("map[string]interface{} converts to immutable map",
			intMapSchema,
			map[string]interface{}{"k": int64(1)},
			func(result interface{}) {
				immutable, ok := result.(*values.Map[string, interface{}])
				Expect(ok).To(BeTrue())
				v, ok := immutable.Get("k")
				Expect(ok).To(BeTrue())
				Expect(v).To(Equal(int64(1)))
			},
		),
		Entry("integer scalar passes through unchanged",
			integerSchema,
			int64(42),
			func(result interface{}) {
				Expect(result).To(Equal(int64(42)))
			},
		),
		Entry("string scalar passes through unchanged",
			stringSchema,
			"hello",
			func(result interface{}) {
				Expect(result).To(Equal("hello"))
			},
		),
		Entry("nested object map deep copies inner collections",
			objectSchema,
			map[string]interface{}{
				"name":  "item",
				"count": int64(1),
				"tags":  []interface{}{"a", "b"},
			},
			func(result interface{}) {
				resultMap, ok := result.(map[string]interface{})
				Expect(ok).To(BeTrue())
				Expect(resultMap["name"]).To(Equal("item"))
				Expect(resultMap["count"]).To(Equal(int64(1)))

				tags, ok := resultMap["tags"].(*values.List[interface{}])
				Expect(ok).To(BeTrue())
				Expect(tags.At(0)).To(Equal("a"))
				Expect(tags.At(1)).To(Equal("b"))
			},
		),
	)

	It("returns the same pointer for immutable list input", func() {
		input := values.NewList("a", "b")
		result, err := bind.BindValue(stringArraySchema, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(BeIdenticalTo(input))
	})

	It("returns the same pointer for immutable map input", func() {
		input := values.NewMapFromGoMap(map[string]int64{"k": 1})
		result, err := bind.BindValue(intMapSchema, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(BeIdenticalTo(input))
	})

	It("does not alias mutable slice input", func() {
		upstream := []interface{}{"a", "b"}
		result, err := bind.BindValue(stringArraySchema, upstream)
		Expect(err).NotTo(HaveOccurred())

		upstream[0] = "mutated"

		list := result.(*values.List[interface{}])
		Expect(list.At(0)).To(Equal("a"))
		Expect(reflect.ValueOf(result).Pointer()).NotTo(Equal(reflect.ValueOf(&upstream).Pointer()))
	})

	It("does not alias mutable map input", func() {
		upstream := map[string]interface{}{"k": int64(1)}
		result, err := bind.BindValue(intMapSchema, upstream)
		Expect(err).NotTo(HaveOccurred())

		upstream["k"] = int64(99)

		immutable := result.(*values.Map[string, interface{}])
		v, ok := immutable.Get("k")
		Expect(ok).To(BeTrue())
		Expect(v).To(Equal(int64(1)))
	})

	It("deep copies nested object inner maps without sharing", func() {
		inner := map[string]interface{}{"nested": "value"}
		upstream := map[string]interface{}{
			"name": "item",
			"meta": inner,
		}
		nestedObjectSchema := &schema.Object{
			Properties: map[string]schema.Schema{
				"name": &schema.String{},
				"meta": &schema.Object{
					Properties: map[string]schema.Schema{
						"nested": &schema.String{},
					},
				},
			},
		}

		result, err := bind.BindValue(nestedObjectSchema, upstream)
		Expect(err).NotTo(HaveOccurred())

		resultMap := result.(map[string]interface{})
		meta := resultMap["meta"].(map[string]interface{})
		Expect(meta).To(Equal(map[string]interface{}{"nested": "value"}))
		Expect(reflect.ValueOf(meta).Pointer()).NotTo(Equal(reflect.ValueOf(inner).Pointer()))

		inner["nested"] = "mutated"
		Expect(meta["nested"]).To(Equal("value"))
	})

	It("returns nil for nil input", func() {
		result, err := bind.BindValue(stringSchema, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(BeNil())
	})

	It("deep copies typed string slices", func() {
		upstream := []string{"a", "b"}
		result, err := bind.BindValue(stringArraySchema, upstream)
		Expect(err).NotTo(HaveOccurred())

		resultSlice, ok := result.([]string)
		Expect(ok).To(BeTrue())
		Expect(resultSlice).To(Equal([]string{"a", "b"}))
		Expect(reflect.ValueOf(resultSlice).Pointer()).NotTo(Equal(reflect.ValueOf(upstream).Pointer()))

		upstream[0] = "mutated"
		Expect(resultSlice[0]).To(Equal("a"))
	})

	It("binds nested mutable collections inside immutable list", func() {
		inner := map[string]interface{}{"nested": "value"}
		input := values.NewList(inner)
		nestedMapArraySchema := &schema.Array{
			Items: &schema.Map{},
		}

		result, err := bind.BindValue(nestedMapArraySchema, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).NotTo(BeIdenticalTo(input))

		list := result.(*values.List[interface{}])
		immutable, ok := list.At(0).(*values.Map[string, interface{}])
		Expect(ok).To(BeTrue())
		v, ok := immutable.Get("nested")
		Expect(ok).To(BeTrue())
		Expect(v).To(Equal("value"))

		inner["nested"] = "mutated"
		v, ok = immutable.Get("nested")
		Expect(ok).To(BeTrue())
		Expect(v).To(Equal("value"))
	})

	It("binds unknown object properties with Any schema", func() {
		inner := map[string]interface{}{"nested": "value"}
		upstream := map[string]interface{}{
			"name":   "item",
			"extra":  inner,
		}
		partialObjectSchema := &schema.Object{
			Properties: map[string]schema.Schema{
				"name": &schema.String{},
			},
		}

		result, err := bind.BindValue(partialObjectSchema, upstream)
		Expect(err).NotTo(HaveOccurred())

		resultMap := result.(map[string]interface{})
		extra := resultMap["extra"].(*values.Map[string, interface{}])
		v, ok := extra.Get("nested")
		Expect(ok).To(BeTrue())
		Expect(v).To(Equal("value"))
		Expect(reflect.ValueOf(extra).Pointer()).NotTo(Equal(reflect.ValueOf(inner).Pointer()))
	})

	It("binds map with Any additional properties", func() {
		anyMapSchema := &schema.Map{}
		upstream := map[string]interface{}{"k": int64(1)}

		result, err := bind.BindValue(anyMapSchema, upstream)
		Expect(err).NotTo(HaveOccurred())

		immutable, ok := result.(*values.Map[string, interface{}])
		Expect(ok).To(BeTrue())
		v, ok := immutable.Get("k")
		Expect(ok).To(BeTrue())
		Expect(v).To(Equal(int64(1)))
	})

	It("does not treat ListBuilder as an immutable handle", func() {
		builder := values.NewListBuilder[string]()
		builder.Append("a")

		result, err := bind.BindValue(stringArraySchema, builder)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).NotTo(BeIdenticalTo(builder))

		list, ok := result.(*values.List[string])
		Expect(ok).To(BeTrue())
		Expect(list.At(0)).To(Equal("a"))
	})
})
