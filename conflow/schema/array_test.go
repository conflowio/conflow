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

	"github.com/conflowio/conflow/conflow/schema"
)

var _ schema.Schema = &schema.Array{}
var _ schema.ArrayKind = &schema.Array{}

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
			err := schema.ValidateValue(value)
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
				MinItems: schema.IntegerPtr(1),
			},
			[]interface{}{int64(1)},
		),
		Entry(
			"enum value - min length zero",
			&schema.Array{
				Items:    schema.IntegerValue(),
				MinItems: schema.IntegerPtr(0),
			},
			[]interface{}{int64(1)},
		),
		Entry(
			"enum value - max length",
			&schema.Array{
				Items:    schema.IntegerValue(),
				MaxItems: schema.IntegerPtr(1),
			},
			[]interface{}{int64(1)},
		),
		Entry(
			"enum value - max length zero",
			&schema.Array{
				Items:    schema.IntegerValue(),
				MaxItems: schema.IntegerPtr(0),
			},
			[]interface{}{},
		),
		Entry(
			"enum value - min and max length",
			&schema.Array{
				Items:    schema.IntegerValue(),
				MinItems: schema.IntegerPtr(1),
				MaxItems: schema.IntegerPtr(2),
			},
			[]interface{}{int64(1)},
		),
	)

	DescribeTable("Validate errors",
		func(schema *schema.Array, value interface{}, expectedErr error) {
			err := schema.ValidateValue(value)
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
			schema.NewFieldError("0", errors.New("must be integer")),
		),
		Entry(
			"invalid array value",
			&schema.Array{
				Items: &schema.Integer{
					Enum: []int64{1, 2},
				},
			},
			[]interface{}{int64(3)},
			schema.NewFieldError("0", errors.New("must be one of 1, 2")),
		),
		Entry(
			"multiple invalid array items",
			&schema.Array{
				Items: schema.IntegerValue(),
			},
			[]interface{}{int64(1), "foo", "bar"},
			schema.ValidationError{Errors: []error{
				schema.NewFieldError("1", errors.New("must be integer")),
				schema.NewFieldError("2", errors.New("must be integer")),
			}},
		),
		Entry(
			"recursive validation",
			&schema.Array{
				Items: &schema.Array{
					Items: schema.IntegerValue(),
				},
			},
			[]interface{}{[]interface{}{"foo"}},
			schema.NewFieldError("0", schema.NewFieldError("0", errors.New("must be integer"))),
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
				MinItems: schema.IntegerPtr(1),
			},
			[]interface{}{},
			errors.New("must have at least one element"),
		),
		Entry(
			"min length 2, 1 item",
			&schema.Array{
				Items:    schema.IntegerValue(),
				MinItems: schema.IntegerPtr(2),
			},
			[]interface{}{int64(1)},
			errors.New("must have at least 2 elements"),
		),
		Entry(
			"max length zero, 1 item",
			&schema.Array{
				Items:    schema.IntegerValue(),
				MaxItems: schema.IntegerPtr(0),
			},
			[]interface{}{int64(1)},
			errors.New("must be empty"),
		),
		Entry(
			"max length 1, 2 items",
			&schema.Array{
				Items:    schema.IntegerValue(),
				MaxItems: schema.IntegerPtr(1),
			},
			[]interface{}{int64(1), int64(2)},
			errors.New("must not contain more than one element"),
		),
		Entry(
			"max length 2, 3 items",
			&schema.Array{
				Items:    schema.IntegerValue(),
				MaxItems: schema.IntegerPtr(2),
			},
			[]interface{}{int64(1), int64(2), int64(3)},
			errors.New("must not contain more than 2 elements"),
		),
		Entry(
			"same min and max length",
			&schema.Array{
				Items:    schema.IntegerValue(),
				MinItems: schema.IntegerPtr(2),
				MaxItems: schema.IntegerPtr(2),
			},
			[]interface{}{int64(1), int64(2), int64(3)},
			errors.New("must have exactly 2 elements"),
		),
	)

	DescribeTable("GoString prints a valid Go struct",
		func(schema *schema.Array, expected string) {
			str := schema.GoString()
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
			"MaxItems",
			&schema.Array{
				MinItems: schema.IntegerPtr(1),
			},
			`&schema.Array{
	MinItems: schema.IntegerPtr(1),
}`,
		),
		Entry(
			"MaxItems",
			&schema.Array{
				MaxItems: schema.IntegerPtr(1),
			},
			`&schema.Array{
	MaxItems: schema.IntegerPtr(1),
}`,
		),
	)

	It("should marshal/unmarshal", func() {
		s := &schema.Array{
			Metadata: schema.Metadata{
				Description: "foo",
			},
			Const:   []interface{}{"constval"},
			Default: []interface{}{"defaultval"},
			Enum:    [][]interface{}{{"enumval"}},
			Items: &schema.String{
				Metadata: schema.Metadata{
					Description: "bar",
				},
			},
		}
		j, err := json.Marshal(s)
		Expect(err).ToNot(HaveOccurred())

		s2 := &schema.Array{}
		err = json.Unmarshal(j, s2)
		Expect(err).ToNot(HaveOccurred())
		Expect(s2).To(Equal(s))
	})
})
