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

	"github.com/conflowio/conflow/pkg/internal/testhelper"
	"github.com/conflowio/conflow/pkg/schema"
	"github.com/conflowio/conflow/pkg/util/validation"
)

var _ schema.Schema = &schema.Array{}

var _ = Describe("Array", func() {

	defaultSchema := func() *schema.Array {
		return &schema.Array{
			Items: &schema.Integer{},
		}
	}

	DescribeTable("CompareValues",
		func(f func(s *schema.Array), v1, v2 interface{}, expected int) {
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
		Entry("[] == nil", nil, []interface{}{}, nil, 0),
		Entry("nil == []", nil, nil, []interface{}{}, 0),
		Entry("[] == []", nil, []interface{}{}, []interface{}{}, 0),
		Entry(
			`[1] == [1]`,
			nil,
			[]interface{}{int64(1)},
			[]interface{}{int64(1)},
			0,
		),
		Entry(
			`[1] < [2]`,
			nil,
			[]interface{}{int64(1)},
			[]interface{}{int64(2)},
			-1,
		),
		Entry(
			`[1] < [1, 2]`,
			nil,
			[]interface{}{int64(1)},
			[]interface{}{int64(1), int64(2)},
			-1,
		),
		Entry(
			`[[1], [2]] == [[1], [2]]`,
			func(s *schema.Array) {
				s.Items = &schema.Array{Items: &schema.Integer{}}
			},
			[]interface{}{[]interface{}{int64(1)}, []interface{}{int64(2)}},
			[]interface{}{[]interface{}{int64(1)}, []interface{}{int64(2)}},
			0,
		),
	)

	DescribeTable("Validate accepts value",
		func(schema *schema.Array, value interface{}) {
			_, err := schema.ValidateValue(value)
			Expect(err).ToNot(HaveOccurred())
		},
		Entry(
			"empty array",
			&schema.Array{},
			[]interface{}{},
		),
		Entry(
			"non-empty array",
			&schema.Array{
				Items: schema.IntegerValue(),
			},
			[]interface{}{int64(1)},
		),
		Entry(
			"const value",
			&schema.Array{
				Items: schema.IntegerValue(),
				Const: []interface{}{int64(1)},
			},
			[]interface{}{int64(1)},
		),
		Entry(
			"const value - empty array",
			&schema.Array{
				Items: schema.IntegerValue(),
				Const: []interface{}{},
			},
			[]interface{}{},
		),
		Entry(
			"enum value - empty array",
			&schema.Array{
				Items: schema.IntegerValue(),
				Enum:  [][]interface{}{{}},
			},
			[]interface{}{},
		),
		Entry(
			"enum value - single",
			&schema.Array{
				Items: schema.IntegerValue(),
				Enum:  [][]interface{}{{int64(1)}},
			},
			[]interface{}{int64(1)},
		),
		Entry(
			"enum value - multiple",
			&schema.Array{
				Items: schema.IntegerValue(),
				Enum:  [][]interface{}{{int64(1)}, {int64(2)}},
			},
			[]interface{}{int64(1)},
		),
		Entry(
			"enum value - min length",
			&schema.Array{
				Items:    schema.IntegerValue(),
				MinItems: 1,
			},
			[]interface{}{int64(1)},
		),
		Entry(
			"enum value - min length zero",
			&schema.Array{
				Items: schema.IntegerValue(),
			},
			[]interface{}{int64(1)},
		),
		Entry(
			"enum value - max length",
			&schema.Array{
				Items:    schema.IntegerValue(),
				MaxItems: schema.Pointer(int64(1)),
			},
			[]interface{}{int64(1)},
		),
		Entry(
			"enum value - max length zero",
			&schema.Array{
				Items:    schema.IntegerValue(),
				MaxItems: schema.Pointer(int64(0)),
			},
			[]interface{}{},
		),
		Entry(
			"enum value - min and max length",
			&schema.Array{
				Items:    schema.IntegerValue(),
				MinItems: 1,
				MaxItems: schema.Pointer(int64(2)),
			},
			[]interface{}{int64(1)},
		),
		Entry(
			"unique items - zero item",
			&schema.Array{
				Items:       schema.IntegerValue(),
				UniqueItems: true,
			},
			[]interface{}{},
		),
		Entry(
			"unique items - one item",
			&schema.Array{
				Items:       schema.IntegerValue(),
				UniqueItems: true,
			},
			[]interface{}{int64(1)},
		),
		Entry(
			"unique items - two items",
			&schema.Array{
				Items:       schema.IntegerValue(),
				UniqueItems: true,
			},
			[]interface{}{int64(1), int64(2)},
		),
		Entry(
			"unique items - complex items",
			&schema.Array{
				Items:       &schema.Array{Items: schema.IntegerValue()},
				UniqueItems: true,
			},
			[]interface{}{
				[]interface{}{int64(1), int64(2)},
				[]interface{}{int64(1), int64(3)},
			},
		),
	)

	DescribeTable("Validate errors",
		func(schema *schema.Array, value interface{}, expectedErr error) {
			_, err := schema.ValidateValue(value)
			Expect(err).To(MatchError(expectedErr))
		},
		Entry(
			"non-array value",
			&schema.Array{},
			"foo",
			errors.New("must be array"),
		),
		Entry(
			"invalid array item",
			&schema.Array{
				Items: schema.IntegerValue(),
			},
			[]interface{}{"foo"},
			validation.NewFieldError("0", errors.New("must be integer")),
		),
		Entry(
			"invalid array value",
			&schema.Array{
				Items: &schema.Integer{
					Enum: []int64{1, 2},
				},
			},
			[]interface{}{int64(3)},
			validation.NewFieldError("0", errors.New("must be one of 1, 2")),
		),
		Entry(
			"multiple invalid array items",
			&schema.Array{
				Items: schema.IntegerValue(),
			},
			[]interface{}{int64(1), "foo", "bar"},
			validation.NewError(
				validation.NewFieldError("1", errors.New("must be integer")),
				validation.NewFieldError("2", errors.New("must be integer")),
			),
		),
		Entry(
			"recursive validation",
			&schema.Array{
				Items: &schema.Array{
					Items: schema.IntegerValue(),
				},
			},
			[]interface{}{[]interface{}{"foo"}},
			validation.NewFieldError("0", validation.NewFieldError("0", errors.New("must be integer"))),
		),
		Entry(
			"const value",
			&schema.Array{
				Items: schema.IntegerValue(),
				Const: []interface{}{int64(1)},
			},
			[]interface{}{int64(2)},
			errors.New("must be [1]"),
		),
		Entry(
			"enum value - single",
			&schema.Array{
				Items: schema.IntegerValue(),
				Enum:  [][]interface{}{{int64(1)}},
			},
			[]interface{}{int64(2)},
			errors.New("must be [1]"),
		),
		Entry(
			"enum value - multiple",
			&schema.Array{
				Items: schema.IntegerValue(),
				Enum:  [][]interface{}{{int64(1)}, {int64(2)}},
			},
			[]interface{}{int64(3)},
			errors.New("must be one of [1], [2]"),
		),
		Entry(
			"enum value - min length 1, zero items",
			&schema.Array{
				Items:    schema.IntegerValue(),
				MinItems: 1,
			},
			[]interface{}{},
			errors.New("must have at least one element"),
		),
		Entry(
			"min length 2, 1 item",
			&schema.Array{
				Items:    schema.IntegerValue(),
				MinItems: 2,
			},
			[]interface{}{int64(1)},
			errors.New("must have at least 2 elements"),
		),
		Entry(
			"max length zero, 1 item",
			&schema.Array{
				Items:    schema.IntegerValue(),
				MaxItems: schema.Pointer(int64(0)),
			},
			[]interface{}{int64(1)},
			errors.New("must be empty"),
		),
		Entry(
			"max length 1, 2 items",
			&schema.Array{
				Items:    schema.IntegerValue(),
				MaxItems: schema.Pointer(int64(1)),
			},
			[]interface{}{int64(1), int64(2)},
			errors.New("must not contain more than one element"),
		),
		Entry(
			"max length 2, 3 items",
			&schema.Array{
				Items:    schema.IntegerValue(),
				MaxItems: schema.Pointer(int64(2)),
			},
			[]interface{}{int64(1), int64(2), int64(3)},
			errors.New("must not contain more than 2 elements"),
		),
		Entry(
			"same min and max length",
			&schema.Array{
				Items:    schema.IntegerValue(),
				MinItems: 2,
				MaxItems: schema.Pointer(int64(2)),
			},
			[]interface{}{int64(1), int64(2), int64(3)},
			errors.New("must have exactly 2 elements"),
		),
		Entry(
			"unique items - same values",
			&schema.Array{
				Items:       schema.IntegerValue(),
				UniqueItems: true,
			},
			[]interface{}{int64(1), int64(2), int64(1)},
			errors.New("array must contain unique items"),
		),
		Entry(
			"unique items - same values, complex items",
			&schema.Array{
				Items:       &schema.Array{Items: schema.IntegerValue()},
				UniqueItems: true,
			},
			[]interface{}{
				[]interface{}{int64(1), int64(2)},
				[]interface{}{int64(1), int64(2)},
			},
			errors.New("array must contain unique items"),
		),
	)

	DescribeTable("GoString prints a valid Go struct",
		func(schema *schema.Array, expected string) {
			str := schema.GoString(map[string]string{})
			Expect(str).To(Equal(expected))
		},
		Entry(
			"empty",
			&schema.Array{},
			`&schema.Array{
}`,
		),
		Entry(
			"const",
			&schema.Array{
				Const: []interface{}{1},
			},
			`&schema.Array{
	Const: []interface {}{1},
}`,
		),
		Entry(
			"default",
			&schema.Array{
				Default: []interface{}{1},
			},
			`&schema.Array{
	Default: []interface {}{1},
}`,
		),
		Entry(
			"enum",
			&schema.Array{
				Enum: [][]interface{}{{1}},
			},
			`&schema.Array{
	Enum: [][]interface {}{[]interface {}{1}},
}`,
		),
		Entry(
			"items",
			&schema.Array{
				Items: &schema.String{Format: "foo"},
			},
			`&schema.Array{
	Items: &schema.String{
		Format: "foo",
	},
}`,
		),
		Entry(
			"MinItems",
			&schema.Array{
				MinItems: 1,
			},
			`&schema.Array{
	MinItems: 1,
}`,
		),
		Entry(
			"MaxItems",
			&schema.Array{
				MaxItems: schema.Pointer(int64(1)),
			},
			`&schema.Array{
	MaxItems: schema.Pointer(int64(1)),
}`,
		),
		Entry(
			"UniqueItems",
			&schema.Array{
				UniqueItems: true,
			},
			`&schema.Array{
	UniqueItems: true,
}`,
		),
	)

	It("should unmarshal/marshal a json", func() {
		testhelper.ExpectConsistentJSONMarshalling(
			`{
				"const": [1],
				"default": [2],
				"enum": [[3, 4], [5, 6]],
				"items": {
					"type": "string"
				},
				"minItems": 7,
				"maxItems": 8,
				"type": "array",
				"uniqueItems": true
			}`,
			&schema.Array{},
		)
	})
})
