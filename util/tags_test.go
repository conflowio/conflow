package util_test

import (
	"sort"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/opsidian/ocl/util"
)

var _ = Describe("GetTags", func() {

	DescribeTable("it parser the valid tags correctly",
		func(obj interface{}, expected util.StructTags) {
			tags := util.GetTags(obj, "ocl")
			Expect(tags).To(Equal(expected))
		},
		Entry("no fields", struct{}{}, nil),
		Entry("one field, no tags", struct{ f string }{}, util.StructTags{"f": nil}),
		Entry("key-value",
			struct {
				f string `ocl:"a=b"`
			}{},
			util.StructTags{
				"f": map[string]string{
					"a": "b",
				},
			},
		),
		Entry("multiple key-values",
			struct {
				f string `ocl:"a=b,c=d"`
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
				f1 string `ocl:"a=b"`
				f2 string `ocl:"c=d"`
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
				f string `ocl:""`
			}{},
			util.StructTags{
				"f": nil,
			},
		),
		Entry("only whitespace in tag",
			struct {
				f string `ocl:" "`
			}{},
			util.StructTags{
				"f": nil,
			},
		),
		Entry("strips whitespaces",
			struct {
				f string `ocl:" a = b , c = d "`
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
				f string `ocl:"foo"`
			}{},
			util.StructTags{
				"f": map[string]string{
					"foo": "true",
				},
			},
		),
		Entry("extra comma should be ignored",
			struct {
				f string `ocl:",a=b,"`
			}{},
			util.StructTags{
				"f": map[string]string{
					"a": "b",
				},
			},
		),
		Entry("pointer",
			&struct {
				f string `ocl:"a=b"`
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
			Expect(func() { util.GetTags(obj, "ocl") }).To(Panic())
		})
	})

	Describe("when tag key is empty", func() {
		It("should panic", func() {
			obj := struct {
				f1 string `ocl:"=value"`
			}{}
			Expect(func() { util.GetTags(obj, "ocl") }).To(Panic())
		})
	})
})

var _ = Describe("FilterFieldsByTags", func() {
	var obj = struct {
		f1 string
		f2 string `ocl:"key=a"`
		f3 string `ocl:"key=a,key2=b"`
		f4 string `ocl:"key2=a"`
		f5 string `ocl:"key=b"`
		f6 string `ocl:"k`
	}{}

	It("should return fields with specific key-value pairs", func() {
		tags := util.GetTags(obj, "ocl")
		actual := util.FilterFieldsByTags(tags, "key", "a")
		sort.Strings(actual)
		Expect(actual).To(Equal([]string{"f2", "f3"}))
	})

	It("should return fields for all filtered key values", func() {
		tags := util.GetTags(obj, "ocl")
		actual := util.FilterFieldsByTags(tags, "key", "a", "b")
		sort.Strings(actual)
		Expect(actual).To(Equal([]string{"f2", "f3", "f5"}))
	})

	It("should return nil if no key matches", func() {
		tags := util.GetTags(obj, "ocl")
		actual := util.FilterFieldsByTags(tags, "x", "foo")
		sort.Strings(actual)
		Expect(actual).To(BeNil())
	})

	It("should return nil if no value matches", func() {
		tags := util.GetTags(obj, "ocl")
		actual := util.FilterFieldsByTags(tags, "key", "x")
		sort.Strings(actual)
		Expect(actual).To(BeNil())
	})

	It("should panic if no value was given", func() {
		tags := util.GetTags(obj, "ocl")
		Expect(func() { util.FilterFieldsByTags(tags, "key") }).To(Panic())
	})
})
