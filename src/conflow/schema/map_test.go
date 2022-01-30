// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema_test

import (
	"errors"

	"github.com/conflowio/conflow/src/internal/testhelper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/src/conflow/schema"
)

var _ schema.Schema = &schema.Map{}
var _ schema.MapKind = &schema.Map{}

var _ = Describe("Map", func() {
	defaultSchema := func() *schema.Map {
		return &schema.Map{
			AdditionalProperties: &schema.Integer{},
		}
	}

	DescribeTable("CompareValues",
		func(f func(s *schema.Map), v1, v2 interface{}, expected int) {
			s := defaultSchema()
			if f != nil {
				f(s)
			}
			res := s.CompareValues(v1, v2)
			Expect(res).To(Equal(expected))

			if expected != 0 {
				res2 := s.CompareValues(v2, v1)
				Expect(res2).To(Equal(0 - expected))
			}
		},
		Entry("nil == nil", nil, nil, nil, 0),
		Entry("map{} == nil", nil, map[string]interface{}{}, nil, 0),
		Entry("nil == map{}", nil, nil, map[string]interface{}{}, 0),
		Entry("map{} == map{}", nil, map[string]interface{}{}, map[string]interface{}{}, 0),
		Entry(
			`map{a=1} == map{a=1}`,
			nil,
			map[string]interface{}{"a": int64(1)},
			map[string]interface{}{"a": int64(1)},
			0,
		),
		Entry(
			`map{a=1} < map{a=2}`,
			nil,
			map[string]interface{}{"a": int64(1)},
			map[string]interface{}{"a": int64(2)},
			-1,
		),
		Entry(
			`map{a=1} < map{b=1}`,
			nil,
			map[string]interface{}{"a": int64(1)},
			map[string]interface{}{"b": int64(1)},
			-1,
		),
		Entry(
			`map{a=1} < map{a=1,b=2}`,
			nil,
			map[string]interface{}{"a": int64(1)},
			map[string]interface{}{"a": int64(1), "b": int64(2)},
			-1,
		),
		Entry(
			`map{a=map{b=1}} == map{a=map{b=1}}`,
			func(s *schema.Map) {
				s.AdditionalProperties = &schema.Map{AdditionalProperties: &schema.Integer{}}
			},
			map[string]interface{}{"a": map[string]interface{}{"b": int64(1)}},
			map[string]interface{}{"a": map[string]interface{}{"b": int64(1)}},
			0,
		),
	)

	DescribeTable("Validate accepts value",
		func(f func(s *schema.Map), value interface{}) {
			s := defaultSchema()
			f(s)
			_, err := s.ValidateValue(value)
			Expect(err).ToNot(HaveOccurred())
		},
		Entry(
			"empty map",
			func(*schema.Map) {},
			map[string]interface{}{},
		),
		Entry(
			"non-empty map",
			func(*schema.Map) {},
			map[string]interface{}{
				"foo": int64(1),
				"bar": int64(2),
			},
		),
		Entry(
			"const value",
			func(s *schema.Map) {
				s.Const = map[string]interface{}{
					"foo": int64(1),
				}
			},
			map[string]interface{}{
				"foo": int64(1),
			},
		),
		Entry(
			"const value - empty map",
			func(s *schema.Map) {
				s.Const = map[string]interface{}{}
			},
			map[string]interface{}{},
		),
		Entry(
			"enum value - empty map",
			func(s *schema.Map) {
				s.Enum = []map[string]interface{}{}
			},
			map[string]interface{}{},
		),
		Entry(
			"enum value - single",
			func(s *schema.Map) {
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
			func(s *schema.Map) {
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
			func(s *schema.Map) {
				s.MinProperties = 1
			},
			map[string]interface{}{
				"foo": int64(1),
			},
		),
		Entry(
			"minProperties=1, 2 elements",
			func(s *schema.Map) {
				s.MinProperties = 2
			},
			map[string]interface{}{
				"foo": int64(1),
				"bar": int64(2),
			},
		),
		Entry(
			"maxProperties=2, 2 elements",
			func(s *schema.Map) {
				s.MaxProperties = schema.IntegerPtr(2)
			},
			map[string]interface{}{
				"foo": int64(1),
				"bar": int64(2),
			},
		),
		Entry(
			"maxProperties=2, 1 element",
			func(s *schema.Map) {
				s.MaxProperties = schema.IntegerPtr(2)
			},
			map[string]interface{}{
				"foo": int64(1),
			},
		),
	)

	DescribeTable("Validate errors",
		func(f func(s *schema.Map), value interface{}, expectedErr error) {
			s := defaultSchema()
			f(s)
			_, err := s.ValidateValue(value)
			Expect(err).To(MatchError(expectedErr))
		},
		Entry(
			"non-map value",
			func(s *schema.Map) {},
			"foo",
			errors.New("must be map"),
		),
		Entry(
			"invalid map value",
			func(s *schema.Map) {},
			map[string]interface{}{
				"foo": "not an integer",
			},
			schema.NewFieldError("foo", errors.New("must be integer")),
		),
		Entry(
			"const value",
			func(s *schema.Map) {
				s.Const = map[string]interface{}{
					"foo": int64(1),
				}
			},
			map[string]interface{}{
				"foo": int64(2),
			},
			errors.New("must be map{foo: 1}"),
		),
		Entry(
			"const value - empty map",
			func(s *schema.Map) {
				s.Const = map[string]interface{}{}
			},
			map[string]interface{}{
				"foo": int64(1),
			},
			errors.New("must be map{}"),
		),
		Entry(
			"enum value - empty map",
			func(s *schema.Map) {
				s.Enum = []map[string]interface{}{{}}
			},
			map[string]interface{}{
				"foo": int64(1),
			},
			errors.New("must be map{}"),
		),
		Entry(
			"enum value - single",
			func(s *schema.Map) {
				s.Enum = []map[string]interface{}{
					{"foo": int64(1)},
				}
			},
			map[string]interface{}{
				"foo": int64(2),
			},
			errors.New("must be map{foo: 1}"),
		),
		Entry(
			"enum value - multiple",
			func(s *schema.Map) {
				s.Enum = []map[string]interface{}{
					{"foo": int64(1)},
					{"foo": int64(2)},
				}
			},
			map[string]interface{}{
				"foo": int64(3),
			},
			errors.New("must be one of map{foo: 1}, map{foo: 2}"),
		),
		Entry(
			"minProperties: 1, empty",
			func(s *schema.Map) {
				s.MinProperties = 1
			},
			map[string]interface{}{},
			errors.New("the map can not be empty"),
		),
		Entry(
			"minProperties: 2, 1 element",
			func(s *schema.Map) {
				s.MinProperties = 2
			},
			map[string]interface{}{
				"foo": int64(1),
			},
			errors.New("the map must contain at least 2 elements"),
		),
		Entry(
			"maxProperties: 0, 1 element",
			func(s *schema.Map) {
				s.MaxProperties = schema.IntegerPtr(0)
			},
			map[string]interface{}{
				"foo": int64(1),
			},
			errors.New("the map must be empty"),
		),
		Entry(
			"maxProperties: 1, 2 elements",
			func(s *schema.Map) {
				s.MaxProperties = schema.IntegerPtr(1)
			},
			map[string]interface{}{
				"foo": int64(1),
				"bar": int64(2),
			},
			errors.New("the map can only have a single element"),
		),
	)

	DescribeTable("GoString prints a valid Go struct",
		func(schema *schema.Map, expected string) {
			str := schema.GoString(map[string]string{})
			Expect(str).To(Equal(expected))
		},
		Entry(
			"empty",
			&schema.Map{},
			`&schema.Map{
}`,
		),
		Entry(
			"additionalProperties",
			&schema.Map{
				AdditionalProperties: &schema.String{Format: "foo"},
			},
			`&schema.Map{
	AdditionalProperties: &schema.String{
		Format: "foo",
	},
}`,
		),
		Entry(
			"const",
			&schema.Map{
				Const: map[string]interface{}{"foo": "bar"},
			},
			`&schema.Map{
	Const: map[string]interface {}{"foo":"bar"},
}`,
		),
		Entry(
			"default",
			&schema.Map{
				Default: map[string]interface{}{"foo": "bar"},
			},
			`&schema.Map{
	Default: map[string]interface {}{"foo":"bar"},
}`,
		),
		Entry(
			"enum",
			&schema.Map{
				Enum: []map[string]interface{}{{"foo": "bar"}},
			},
			`&schema.Map{
	Enum: []map[string]interface {}{map[string]interface {}{"foo":"bar"}},
}`,
		),
		Entry(
			"minProperties",
			&schema.Map{
				MinProperties: 1,
			},
			`&schema.Map{
	MinProperties: 1,
}`,
		),
		Entry(
			"maxProperties",
			&schema.Map{
				MaxProperties: schema.IntegerPtr(1),
			},
			`&schema.Map{
	MaxProperties: schema.IntegerPtr(1),
}`,
		),
	)

	It("should unmarshal/marshal a json", func() {
		testhelper.ExpectConsistentJSONMarshalling(
			`{
				"additionalProperties": {
					"type": "string"
				},
				"const": {
					"p1": "v1"
				},
				"default": {
					"p2": "v2"
				},
				"enum": [
					{
						"p3": "v3"
					},
					{
						"p4": "v4"
					}
				],
				"minProperties": 1,
				"maxProperties": 2,
				"type": "object"
			}`,
			&schema.Map{},
		)
	})
})
