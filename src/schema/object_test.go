// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/src/internal/testhelper"
	"github.com/conflowio/conflow/src/schema"
)

var _ schema.Schema = &schema.Object{}
var _ schema.ObjectKind = &schema.Object{}

var _ = Describe("Object", func() {
	defaultSchema := func() *schema.Object {
		return &schema.Object{
			Parameters: map[string]schema.Schema{
				"foo": &schema.Integer{},
				"bar": &schema.String{},
				"baz": &schema.Object{
					Parameters: map[string]schema.Schema{
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
			_, err := s.ValidateValue(value)
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
				s.Const = map[string]interface{}{
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
				s.Const = map[string]interface{}{}
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
		Entry(
			"minProperties=1, 1 element",
			func(s *schema.Object) {
				s.MinProperties = 1
			},
			map[string]interface{}{
				"foo": int64(1),
			},
		),
		Entry(
			"minProperties=1, 2 elements",
			func(s *schema.Object) {
				s.MinProperties = 2
			},
			map[string]interface{}{
				"foo": int64(1),
				"bar": "2",
			},
		),
		Entry(
			"maxProperties=2, 2 elements",
			func(s *schema.Object) {
				s.MaxProperties = schema.IntegerPtr(2)
			},
			map[string]interface{}{
				"foo": int64(1),
				"bar": "2",
			},
		),
		Entry(
			"maxProperties=2, 1 element",
			func(s *schema.Object) {
				s.MaxProperties = schema.IntegerPtr(2)
			},
			map[string]interface{}{
				"foo": int64(1),
			},
		),
		Entry(
			"dependentRequired",
			func(s *schema.Object) {
				s.DependentRequired = map[string][]string{
					"foo": {"bar"},
				}
			},
			map[string]interface{}{
				"foo": int64(1),
				"bar": "val",
			},
		),
	)

	DescribeTable("Validate errors",
		func(f func(s *schema.Object), value interface{}, expectedErr error) {
			s := defaultSchema()
			f(s)
			_, err := s.ValidateValue(value)
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
				s.Const = map[string]interface{}{
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
				s.Const = map[string]interface{}{}
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
		Entry(
			"minProperties: 1, empty",
			func(s *schema.Object) {
				s.MinProperties = 1
			},
			map[string]interface{}{},
			errors.New("the object can not be empty"),
		),
		Entry(
			"minProperties: 2, 1 element",
			func(s *schema.Object) {
				s.MinProperties = 2
			},
			map[string]interface{}{
				"foo": int64(1),
			},
			errors.New("the object must have at least 2 properties defined"),
		),
		Entry(
			"maxProperties: 0, 1 element",
			func(s *schema.Object) {
				s.MaxProperties = schema.IntegerPtr(0)
			},
			map[string]interface{}{
				"foo": int64(1),
			},
			errors.New("the object must be empty"),
		),
		Entry(
			"maxProperties: 1, 2 elements",
			func(s *schema.Object) {
				s.MaxProperties = schema.IntegerPtr(1)
			},
			map[string]interface{}{
				"foo": int64(1),
				"bar": int64(2),
			},
			errors.New("the object can only have a single property defined"),
		),
		Entry(
			"dependentRequired - one missing",
			func(s *schema.Object) {
				s.DependentRequired = map[string][]string{
					"foo": {"bar"},
				}
			},
			map[string]interface{}{
				"foo": int64(1),
			},
			schema.NewFieldError("bar", errors.New("required")),
		),
	)

	DescribeTable("GoString prints a valid Go struct",
		func(schema *schema.Object, expected string) {
			str := schema.GoString(map[string]string{})
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
				Const: map[string]interface{}{"foo": "bar"},
			},
			`&schema.Object{
	Const: map[string]interface {}{"foo":"bar"},
}`,
		),
		Entry(
			"default",
			&schema.Object{
				Default: map[string]interface{}{"foo": "bar"},
			},
			`&schema.Object{
	Default: map[string]interface {}{"foo":"bar"},
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
			"parameters",
			&schema.Object{
				Parameters: map[string]schema.Schema{
					"bar": &schema.String{Format: "f1"},
					"foo": &schema.String{Format: "f2"},
				},
			},
			`&schema.Object{
	Parameters: map[string]schema.Schema{
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
			"fieldNames",
			&schema.Object{
				FieldNames: map[string]string{"myField": "MyField"},
			},
			`&schema.Object{
	FieldNames: map[string]string{"myField":"MyField"},
}`,
		),
		Entry(
			"JSON property names",
			&schema.Object{
				JSONPropertyNames: map[string]string{"my_field": "myField"},
			},
			`&schema.Object{
	JSONPropertyNames: map[string]string{"my_field":"myField"},
}`,
		),
		Entry(
			"minProperties",
			&schema.Object{
				MinProperties: 1,
			},
			`&schema.Object{
	MinProperties: 1,
}`,
		),
		Entry(
			"maxProperties",
			&schema.Object{
				MaxProperties: schema.IntegerPtr(1),
			},
			`&schema.Object{
	MaxProperties: schema.IntegerPtr(1),
}`,
		),
		Entry(
			"dependentRequired",
			&schema.Object{
				DependentRequired: map[string][]string{"foo": {"bar"}},
			},
			`&schema.Object{
	DependentRequired: map[string][]string{"foo":[]string{"bar"}},
}`,
		),
	)

	It("should unmarshal/marshal a json", func() {
		testhelper.ExpectConsistentJSONMarshalling(
			`{
				"const": {
					"myField": "val1"
				},
				"default": {
					"myField": "val2"
				},
				"dependentRequired": {
					"foo": ["bar"]
				},
				"enum": [
					{
						"myField": "val3"
					},
					{
						"myField": "val4"
					}
				],
				"type": "object",
				"properties": {
					"myField": {
						"type": "string"
					}
				},
				"fieldNames": {
					"myField": "MyField"
				},
				"minProperties": 1,
				"maxProperties": 2,
				"parameterNames": {
					"myField": "my_field"
				},
				"required": ["myField"]
			}`,
			&schema.Object{},
		)
	})

	It("should ignore an invalid parameter name", func() {
		input := `{
			"type": "object",
			"properties": {
				"myField": {
					"type": "string"
				}
			},
			"parameterNames": {
				"myField": "MyField"
			}
		}`
		s, err := schema.UnmarshalJSON([]byte(input))
		Expect(err).ToNot(HaveOccurred())
		Expect(s).To(Equal(&schema.Object{
			JSONPropertyNames: map[string]string{
				"my_field": "myField",
			},
			Parameters: map[string]schema.Schema{
				"my_field": &schema.String{},
			},
		}))
	})

	It("should ignore an invalid field name", func() {
		input := `{
			"type": "object",
			"properties": {
				"my_field": {
					"type": "string"
				},
				"my_field_2": {
					"type": "string"
				}
			},
			"fieldNames": {
				"my_field": "field name",
				"my_field_2": "myField2"
			}
		}`
		s, err := schema.UnmarshalJSON([]byte(input))
		Expect(err).ToNot(HaveOccurred())
		Expect(s).To(Equal(&schema.Object{
			Parameters: map[string]schema.Schema{
				"my_field":   &schema.String{},
				"my_field_2": &schema.String{},
			},
			FieldNames: map[string]string{
				"my_field_2": "myField2",
			},
		}))
	})

})
