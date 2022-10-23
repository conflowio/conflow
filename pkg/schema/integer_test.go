// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema_test

import (
	"errors"
	"fmt"
	"math"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/pkg/internal/testhelper"
	"github.com/conflowio/conflow/pkg/schema"
)

var _ schema.Schema = &schema.Integer{}

var _ = Describe("Integer", func() {
	DescribeTable("Validate accepts value",
		func(schema *schema.Integer, value interface{}) {
			_, err := schema.ValidateValue(value)
			Expect(err).ToNot(HaveOccurred())
		},
		Entry("zero", &schema.Integer{}, int64(0)),
		Entry("non-zero value", &schema.Integer{}, int64(1)),
		Entry("const value", &schema.Integer{Const: schema.Pointer(int64(1))}, int64(1)),
		Entry("enum value - single", &schema.Integer{Enum: []int64{1}}, int64(1)),
		Entry("enum value - multiple", &schema.Integer{Enum: []int64{1, 2}}, int64(1)),
		Entry("enum value - minimum - equal", &schema.Integer{Minimum: schema.Pointer(int64(1))}, int64(1)),
		Entry("enum value - minimum - greater", &schema.Integer{Minimum: schema.Pointer(int64(1))}, int64(2)),
		Entry("enum value - maximum - equal", &schema.Integer{Maximum: schema.Pointer(int64(2))}, int64(2)),
		Entry("enum value - maximum - less", &schema.Integer{Maximum: schema.Pointer(int64(2))}, int64(1)),
		Entry("enum value - exclusive minimum", &schema.Integer{ExclusiveMinimum: schema.Pointer(int64(1))}, int64(2)),
		Entry("enum value - exclusive maximum", &schema.Integer{ExclusiveMaximum: schema.Pointer(int64(2))}, int64(1)),
		Entry("enum value - multiple of", &schema.Integer{MultipleOf: schema.Pointer(int64(2))}, int64(4)),
		Entry("int32 with maximum valid value", &schema.Integer{Format: "int32"}, int64(math.MaxInt32)),
		Entry("int32 with minimum valid value", &schema.Integer{Format: "int32"}, int64(math.MinInt32)),
	)

	DescribeTable("Validate errors",
		func(schema *schema.Integer, value interface{}, expectedErr error) {
			_, err := schema.ValidateValue(value)
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
			&schema.Integer{Const: schema.Pointer(int64(1))},
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
			&schema.Integer{Minimum: schema.Pointer(int64(2))},
			int64(1),
			errors.New("must be greater than or equal to 2"),
		),
		Entry(
			"maximum",
			&schema.Integer{Maximum: schema.Pointer(int64(2))},
			int64(3),
			errors.New("must be less than or equal to 2"),
		),
		Entry(
			"exclusive minimum - equal",
			&schema.Integer{ExclusiveMinimum: schema.Pointer(int64(2))},
			int64(2),
			errors.New("must be greater than 2"),
		),
		Entry(
			"exclusive minimum - less",
			&schema.Integer{ExclusiveMinimum: schema.Pointer(int64(2))},
			int64(1),
			errors.New("must be greater than 2"),
		),
		Entry(
			"exclusive maximum - equal",
			&schema.Integer{ExclusiveMaximum: schema.Pointer(int64(2))},
			int64(2),
			errors.New("must be less than 2"),
		),
		Entry(
			"exclusive maximum - greater",
			&schema.Integer{ExclusiveMaximum: schema.Pointer(int64(2))},
			int64(3),
			errors.New("must be less than 2"),
		),
		Entry(
			"multiple of",
			&schema.Integer{MultipleOf: schema.Pointer(int64(2))},
			int64(3),
			errors.New("must be multiple of 2"),
		),
		Entry(
			"int32 minimum",
			&schema.Integer{Format: "int32"},
			int64(math.MinInt32-1),
			fmt.Errorf("must be greater than or equal to %d", math.MinInt32),
		),
		Entry(
			"int32 maximum",
			&schema.Integer{Format: "int32"},
			int64(math.MaxInt32+1),
			fmt.Errorf("must be less than or equal to %d", math.MaxInt32),
		),
	)

	DescribeTable("GoString prints a valid Go struct",
		func(schema *schema.Integer, expected string) {
			str := schema.GoString(map[string]string{})
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
			&schema.Integer{Const: schema.Pointer(int64(1))},
			`&schema.Integer{
	Const: schema.Pointer(int64(1)),
}`,
		),
		Entry(
			"default",
			&schema.Integer{Default: schema.Pointer(int64(1))},
			`&schema.Integer{
	Default: schema.Pointer(int64(1)),
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
			"format",
			&schema.Integer{Format: "int32"},
			`&schema.Integer{
	Format: "int32",
}`,
		),
		Entry(
			"minimum",
			&schema.Integer{Minimum: schema.Pointer(int64(1))},
			`&schema.Integer{
	Minimum: schema.Pointer(int64(1)),
}`,
		),
		Entry(
			"maximum",
			&schema.Integer{Maximum: schema.Pointer(int64(1))},
			`&schema.Integer{
	Maximum: schema.Pointer(int64(1)),
}`,
		),
		Entry(
			"exclusive minimum",
			&schema.Integer{ExclusiveMinimum: schema.Pointer(int64(1))},
			`&schema.Integer{
	ExclusiveMinimum: schema.Pointer(int64(1)),
}`,
		),
		Entry(
			"exclusive maximum",
			&schema.Integer{ExclusiveMaximum: schema.Pointer(int64(1))},
			`&schema.Integer{
	ExclusiveMaximum: schema.Pointer(int64(1)),
}`,
		),
		Entry(
			"multiple of",
			&schema.Integer{MultipleOf: schema.Pointer(int64(1))},
			`&schema.Integer{
	MultipleOf: schema.Pointer(int64(1)),
}`,
		),
		Entry(
			"nullable",
			&schema.Integer{Nullable: true},
			`&schema.Integer{
	Nullable: true,
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
				"format": "int32",
				"maximum": 7,
				"minimum": 8,
				"multipleOf": 9,
				"nullable": true,
				"type": "integer"
			}`,
			&schema.Integer{},
		)
	})

})
