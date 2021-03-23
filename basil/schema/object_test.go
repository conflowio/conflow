// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema_test

import (
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/opsidian/basil/basil/schema"
)

var _ schema.Schema = &schema.Object{}
var _ schema.ObjectKind = &schema.Object{}

var _ = Describe("Object", func() {
	defaultSchema := func() *schema.Object {
		return &schema.Object{
			Properties: map[string]schema.Schema{
				"foo": &schema.Integer{},
				"bar": &schema.String{},
				"baz": &schema.Object{
					Properties: map[string]schema.Schema{
						"qux": &schema.Boolean{},
					},
				},
			},
		}
	}

	DescribeTable("Validate accepts value",
		func(f func(s *schema.Object), value interface{}) {
			s := defaultSchema()
			f(s)
			err := s.ValidateValue(value)
			Expect(err).ToNot(HaveOccurred())
		},
		Entry(
			"empty object",
			func(*schema.Object) {},
			map[string]interface{}{},
		),
		Entry(
			"non-empty object",
			func(*schema.Object) {},
			map[string]interface{}{
				"foo": int64(1),
			},
		),
		Entry(
			"complex object",
			func(*schema.Object) {},
			map[string]interface{}{
				"foo": int64(1),
				"bar": "value",
				"baz": map[string]interface{}{
					"qux": true,
				},
			},
		),
		Entry(
			"required value set",
			func(s *schema.Object) {
				s.Required = []string{"foo"}
			},
			map[string]interface{}{
				"foo": int64(1),
			},
		),
		Entry(
			"const value",
			func(s *schema.Object) {
				s.Const = &map[string]interface{}{
					"foo": int64(1),
				}
			},
			map[string]interface{}{
				"foo": int64(1),
			},
		),
		Entry(
			"const value - empty object",
			func(s *schema.Object) {
				s.Const = &map[string]interface{}{}
			},
			map[string]interface{}{},
		),
		Entry(
			"enum value - empty object",
			func(s *schema.Object) {
				s.Enum = []map[string]interface{}{}
			},
			map[string]interface{}{},
		),
		Entry(
			"enum value - single",
			func(s *schema.Object) {
				s.Enum = []map[string]interface{}{
					{"foo": int64(1)},
				}
			},
			map[string]interface{}{
				"foo": int64(1),
			},
		),
		Entry(
			"enum value - multiple",
			func(s *schema.Object) {
				s.Enum = []map[string]interface{}{
					{"foo": int64(1)},
					{"foo": int64(2)},
				}
			},
			map[string]interface{}{
				"foo": int64(1),
			},
		),
	)

	DescribeTable("Validate errors",
		func(f func(s *schema.Object), value interface{}, expectedErr error) {
			s := defaultSchema()
			f(s)
			err := s.ValidateValue(value)
			Expect(err).To(MatchError(expectedErr))
		},
		Entry(
			"non-object value",
			func(s *schema.Object) {},
			"foo",
			errors.New("must be object"),
		),
		Entry(
			"invalid object property value",
			func(s *schema.Object) {},
			map[string]interface{}{
				"unknown": int64(1),
			},
			schema.NewFieldError("unknown", errors.New("property does not exist")),
		),
		Entry(
			"invalid object property value",
			func(s *schema.Object) {},
			map[string]interface{}{
				"foo": "not an integer",
			},
			schema.NewFieldError("foo", errors.New("must be integer")),
		),
		Entry(
			"invalid object property on child object",
			func(*schema.Object) {},
			map[string]interface{}{
				"foo": int64(1),
				"baz": map[string]interface{}{
					"qux": int64(1),
				},
			},
			schema.NewFieldError("baz", schema.NewFieldError("qux", errors.New("must be boolean"))),
		),
		Entry(
			"required value not set",
			func(s *schema.Object) {
				s.Required = []string{"foo"}
			},
			map[string]interface{}{},
			schema.NewFieldError("foo", errors.New("required")),
		),
		Entry(
			"const value",
			func(s *schema.Object) {
				s.Const = &map[string]interface{}{
					"foo": int64(1),
				}
			},
			map[string]interface{}{
				"foo": int64(2),
			},
			errors.New("must be {foo: 1}"),
		),
		Entry(
			"const value - empty object",
			func(s *schema.Object) {
				s.Const = &map[string]interface{}{}
			},
			map[string]interface{}{
				"foo": int64(1),
			},
			errors.New("must be {}"),
		),
		Entry(
			"enum value - empty object",
			func(s *schema.Object) {
				s.Enum = []map[string]interface{}{{}}
			},
			map[string]interface{}{
				"foo": int64(1),
			},
			errors.New("must be {}"),
		),
		Entry(
			"enum value - single",
			func(s *schema.Object) {
				s.Enum = []map[string]interface{}{
					{"foo": int64(1)},
				}
			},
			map[string]interface{}{
				"foo": int64(2),
			},
			errors.New("must be {foo: 1}"),
		),
		Entry(
			"enum value - multiple",
			func(s *schema.Object) {
				s.Enum = []map[string]interface{}{
					{"foo": int64(1)},
					{"foo": int64(2)},
				}
			},
			map[string]interface{}{
				"foo": int64(3),
			},
			errors.New("must be one of {foo: 1}, {foo: 2}"),
		),
	)

	DescribeTable("GoString prints a valid Go struct",
		func(schema *schema.Object, expected string) {
			str := schema.GoString()
			Expect(str).To(Equal(expected))
		},
		Entry(
			"empty",
			&schema.Object{},
			`&schema.Object{
}`,
		),
		Entry(
			"const",
			&schema.Object{
				Const: &map[string]interface{}{"foo": "bar"},
			},
			`&schema.Object{
	Const: &map[string]interface {}{"foo":"bar"},
}`,
		),
		Entry(
			"default",
			&schema.Object{
				Default: &map[string]interface{}{"foo": "bar"},
			},
			`&schema.Object{
	Default: &map[string]interface {}{"foo":"bar"},
}`,
		),
		Entry(
			"enum",
			&schema.Object{
				Enum: []map[string]interface{}{{"foo": "bar"}},
			},
			`&schema.Object{
	Enum: []map[string]interface {}{map[string]interface {}{"foo":"bar"}},
}`,
		),
		Entry(
			"properties",
			&schema.Object{
				Properties: map[string]schema.Schema{
					"bar": &schema.String{Format: "f1"},
					"foo": &schema.String{Format: "f2"},
				},
			},
			`&schema.Object{
	Properties: map[string]schema.Schema{
		"bar": &schema.String{
			Format: "f1",
		},
		"foo": &schema.String{
			Format: "f2",
		},
	},
}`,
		),
		Entry(
			"required",
			&schema.Object{
				Required: []string{"foo"},
			},
			`&schema.Object{
	Required: []string{"foo"},
}`,
		),
		Entry(
			"structProperties",
			&schema.Object{
				StructProperties: map[string]string{"foo": "Foo"},
			},
			`&schema.Object{
	StructProperties: map[string]string{"foo":"Foo"},
}`,
		),
	)

	It("should marshal/unmarshal", func() {
		s := &schema.Object{
			Metadata: schema.Metadata{
				Annotations: map[string]string{
					"foo": "bar",
				},
				Description: "foo",
			},
			Properties: map[string]schema.Schema{
				"baz": &schema.String{
					Metadata: schema.Metadata{
						Description: "qux",
					},
				},
			},
			StructProperties: map[string]string{
				"baz": "Baz",
			},
		}
		j, err := json.Marshal(s)
		Expect(err).ToNot(HaveOccurred())

		s2 := &schema.Object{}
		err = json.Unmarshal(j, s2)
		Expect(err).ToNot(HaveOccurred())
		Expect(s2).To(Equal(s))
	})
})
