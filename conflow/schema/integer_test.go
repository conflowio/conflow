// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema_test

import (
	"errors"

	"github.com/conflowio/conflow/internal/testhelper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/conflow/schema"
)

var _ schema.Schema = &schema.Integer{}

var _ = Describe("Integer", func() {
	DescribeTable("Validate accepts value",
		func(schema *schema.Integer, value interface{}) {
			err := schema.ValidateValue(value)
			Expect(err).ToNot(HaveOccurred())
		},
		Entry("zero", &schema.Integer{}, int64(0)),
		Entry("non-zero value", &schema.Integer{}, int64(1)),
		Entry("const value", &schema.Integer{Const: schema.IntegerPtr(1)}, int64(1)),
		Entry("enum value - single", &schema.Integer{Enum: []int64{1}}, int64(1)),
		Entry("enum value - multiple", &schema.Integer{Enum: []int64{1, 2}}, int64(1)),
		Entry("enum value - minimum - equal", &schema.Integer{Minimum: schema.IntegerPtr(1)}, int64(1)),
		Entry("enum value - minimum - greater", &schema.Integer{Minimum: schema.IntegerPtr(1)}, int64(2)),
		Entry("enum value - maximum - equal", &schema.Integer{Maximum: schema.IntegerPtr(2)}, int64(2)),
		Entry("enum value - maximum - less", &schema.Integer{Maximum: schema.IntegerPtr(2)}, int64(1)),
		Entry("enum value - exclusive minimum", &schema.Integer{ExclusiveMinimum: schema.IntegerPtr(1)}, int64(2)),
		Entry("enum value - exclusive maximum", &schema.Integer{ExclusiveMaximum: schema.IntegerPtr(2)}, int64(1)),
		Entry("enum value - multiple of", &schema.Integer{MultipleOf: schema.IntegerPtr(2)}, int64(4)),
	)

	DescribeTable("Validate errors",
		func(schema *schema.Integer, value interface{}, expectedErr error) {
			err := schema.ValidateValue(value)
			Expect(err).To(MatchError(expectedErr))
		},
		Entry(
			"non-integer value",
			&schema.Integer{},
			"foo",
			errors.New("must be integer"),
		),
		Entry(
			"const value",
			&schema.Integer{Const: schema.IntegerPtr(1)},
			int64(2),
			errors.New("must be 1"),
		),
		Entry(
			"enum value - single",
			&schema.Integer{Enum: []int64{1}},
			int64(2),
			errors.New("must be 1"),
		),
		Entry(
			"enum value - multiple",
			&schema.Integer{Enum: []int64{1, 2}},
			int64(3),
			errors.New("must be one of 1, 2"),
		),
		Entry(
			"minimum",
			&schema.Integer{Minimum: schema.IntegerPtr(2)},
			int64(1),
			errors.New("must be greater than or equal to 2"),
		),
		Entry(
			"maximum",
			&schema.Integer{Maximum: schema.IntegerPtr(2)},
			int64(3),
			errors.New("must be less than or equal to 2"),
		),
		Entry(
			"exclusive minimum - equal",
			&schema.Integer{ExclusiveMinimum: schema.IntegerPtr(2)},
			int64(2),
			errors.New("must be greater than 2"),
		),
		Entry(
			"exclusive minimum - less",
			&schema.Integer{ExclusiveMinimum: schema.IntegerPtr(2)},
			int64(1),
			errors.New("must be greater than 2"),
		),
		Entry(
			"exclusive maximum - equal",
			&schema.Integer{ExclusiveMaximum: schema.IntegerPtr(2)},
			int64(2),
			errors.New("must be less than 2"),
		),
		Entry(
			"exclusive maximum - greater",
			&schema.Integer{ExclusiveMaximum: schema.IntegerPtr(2)},
			int64(3),
			errors.New("must be less than 2"),
		),
		Entry(
			"multiple of",
			&schema.Integer{MultipleOf: schema.IntegerPtr(2)},
			int64(3),
			errors.New("must be multiple of 2"),
		),
	)

	DescribeTable("GoString prints a valid Go struct",
		func(schema *schema.Integer, expected string) {
			str := schema.GoString()
			Expect(str).To(Equal(expected))
		},
		Entry(
			"empty",
			&schema.Integer{},
			`&schema.Integer{
}`,
		),
		Entry(
			"const",
			&schema.Integer{Const: schema.IntegerPtr(1)},
			`&schema.Integer{
	Const: schema.IntegerPtr(1),
}`,
		),
		Entry(
			"default",
			&schema.Integer{Default: schema.IntegerPtr(1)},
			`&schema.Integer{
	Default: schema.IntegerPtr(1),
}`,
		),
		Entry(
			"enum",
			&schema.Integer{Enum: []int64{1}},
			`&schema.Integer{
	Enum: []int64{1},
}`,
		),
		Entry(
			"minimum",
			&schema.Integer{Minimum: schema.IntegerPtr(1)},
			`&schema.Integer{
	Minimum: schema.IntegerPtr(1),
}`,
		),
		Entry(
			"maximum",
			&schema.Integer{Maximum: schema.IntegerPtr(1)},
			`&schema.Integer{
	Maximum: schema.IntegerPtr(1),
}`,
		),
		Entry(
			"exclusive minimum",
			&schema.Integer{ExclusiveMinimum: schema.IntegerPtr(1)},
			`&schema.Integer{
	ExclusiveMinimum: schema.IntegerPtr(1),
}`,
		),
		Entry(
			"exclusive maximum",
			&schema.Integer{ExclusiveMaximum: schema.IntegerPtr(1)},
			`&schema.Integer{
	ExclusiveMaximum: schema.IntegerPtr(1),
}`,
		),
		Entry(
			"multiple of",
			&schema.Integer{MultipleOf: schema.IntegerPtr(1)},
			`&schema.Integer{
	MultipleOf: schema.IntegerPtr(1),
}`,
		),
	)

	It("should unmarshal/marshal a json", func() {
		testhelper.ExpectConsistentJSONMarshalling(
			`{
				"const": 1,
				"default": 2,
				"enum": [3, 4],
				"exclusiveMinimum": 5,
				"exclusiveMaximum": 6,
				"maximum": 7,
				"minimum": 8,
				"multipleOf": 9,
				"type": "integer"
			}`,
			&schema.Integer{},
		)
	})

})
