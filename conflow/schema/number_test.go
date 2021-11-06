// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema_test

import (
	"encoding/json"
	"errors"

	"github.com/conflowio/conflow/internal/testhelper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/conflow/schema"
)

var _ schema.Schema = &schema.Number{}

var _ = Describe("Number", func() {
	// Used in "enum value - multiple of - big numbers"
	a := 185736583293475.127285723743
	b := 8562629561893.0
	c := a * b

	DescribeTable("Validate accepts value",
		func(schema *schema.Number, value interface{}) {
			err := schema.ValidateValue(value)
			Expect(err).ToNot(HaveOccurred())
		},
		Entry("zero", &schema.Number{}, float64(0)),
		Entry("non-zero value", &schema.Number{}, float64(1)),
		Entry("integer", &schema.Number{}, int64(1)),
		Entry("const value", &schema.Number{Const: schema.NumberPtr(1)}, float64(1)),
		Entry("const value integer", &schema.Number{Const: schema.NumberPtr(1)}, int64(1)),
		Entry("enum value - single", &schema.Number{Enum: []float64{1}}, float64(1)),
		Entry("enum value - integer", &schema.Number{Enum: []float64{1}}, int64(1)),
		Entry("enum value - multiple", &schema.Number{Enum: []float64{1, 2}}, float64(1)),

		Entry("enum value - minimum - equal", &schema.Number{Minimum: schema.NumberPtr(1)}, 1.0),
		Entry("enum value - minimum - equal - eps", &schema.Number{Minimum: schema.NumberPtr(1)}, 1-schema.Epsilon*0.1),
		Entry("enum value - minimum - greater", &schema.Number{Minimum: schema.NumberPtr(1)}, 2.0),
		Entry("enum value - maximum - equal", &schema.Number{Maximum: schema.NumberPtr(2)}, 2.0),
		Entry("enum value - maximum - equal + eps", &schema.Number{Maximum: schema.NumberPtr(2)}, 2.0+schema.Epsilon*0.1),
		Entry("enum value - maximum - less", &schema.Number{Maximum: schema.NumberPtr(2)}, 1.0),
		Entry("enum value - exclusive minimum", &schema.Number{ExclusiveMinimum: schema.NumberPtr(1)}, 2.0),
		Entry("enum value - exclusive minimum - eps", &schema.Number{ExclusiveMinimum: schema.NumberPtr(1)}, 1.0+schema.Epsilon),
		Entry("enum value - exclusive maximum", &schema.Number{ExclusiveMaximum: schema.NumberPtr(2)}, 1.0),
		Entry("enum value - exclusive maximum", &schema.Number{ExclusiveMaximum: schema.NumberPtr(2)}, 2.0-schema.Epsilon),
		Entry("enum value - multiple of", &schema.Number{MultipleOf: schema.NumberPtr(2)}, 4.0),
		Entry("enum value - multiple of - eps 1", &schema.Number{MultipleOf: schema.NumberPtr(2)}, 4.0-schema.Epsilon*0.1),
		Entry("enum value - multiple of - eps 2", &schema.Number{MultipleOf: schema.NumberPtr(2)}, 4.0+schema.Epsilon*0.1),
		Entry("enum value - multiple of - big numbers", &schema.Number{MultipleOf: schema.NumberPtr(a)}, c),
	)

	DescribeTable("Validate errors",
		func(schema *schema.Number, value interface{}, expectedErr error) {
			err := schema.ValidateValue(value)
			Expect(err).To(MatchError(expectedErr))
		},
		Entry(
			"non-number value",
			&schema.Number{},
			"foo",
			errors.New("must be number"),
		),
		Entry(
			"const value",
			&schema.Number{Const: schema.NumberPtr(1)},
			float64(2),
			errors.New("must be 1"),
		),
		Entry(
			"enum value - single",
			&schema.Number{Enum: []float64{1}},
			float64(2),
			errors.New("must be 1"),
		),
		Entry(
			"enum value - multiple",
			&schema.Number{Enum: []float64{1, 2}},
			float64(3),
			errors.New("must be one of 1, 2"),
		),
		Entry(
			"minimum",
			&schema.Number{Minimum: schema.NumberPtr(2)},
			1.0,
			errors.New("must be greater than or equal to 2"),
		),
		Entry(
			"minimum - eps",
			&schema.Number{Minimum: schema.NumberPtr(2)},
			2.0-schema.Epsilon,
			errors.New("must be greater than or equal to 2"),
		),
		Entry(
			"exclusive minimum",
			&schema.Number{ExclusiveMinimum: schema.NumberPtr(2)},
			1.0,
			errors.New("must be greater than 2"),
		),
		Entry(
			"exclusive minimum - equals",
			&schema.Number{ExclusiveMinimum: schema.NumberPtr(2)},
			2.0,
			errors.New("must be greater than 2"),
		),
		Entry(
			"exclusive minimum + eps",
			&schema.Number{ExclusiveMinimum: schema.NumberPtr(2)},
			2.0+schema.Epsilon*0.1,
			errors.New("must be greater than 2"),
		),
		Entry(
			"maximum",
			&schema.Number{Maximum: schema.NumberPtr(1)},
			2.0,
			errors.New("must be less than or equal to 1"),
		),
		Entry(
			"maximum - eps",
			&schema.Number{Maximum: schema.NumberPtr(1)},
			1.0+schema.Epsilon,
			errors.New("must be less than or equal to 1"),
		),
		Entry(
			"exclusive maximum",
			&schema.Number{ExclusiveMaximum: schema.NumberPtr(1)},
			2.0,
			errors.New("must be less than 1"),
		),
		Entry(
			"exclusive maximum - equals",
			&schema.Number{ExclusiveMaximum: schema.NumberPtr(1)},
			1.0,
			errors.New("must be less than 1"),
		),
		Entry(
			"exclusive maximum - eps",
			&schema.Number{ExclusiveMaximum: schema.NumberPtr(1)},
			1.0-schema.Epsilon*0.1,
			errors.New("must be less than 1"),
		),
	)

	DescribeTable("GoString prints a valid Go struct",
		func(schema *schema.Number, expected string) {
			str := schema.GoString()
			Expect(str).To(Equal(expected))
		},
		Entry(
			"empty",
			&schema.Number{},
			`&schema.Number{
}`,
		),
		Entry(
			"const",
			&schema.Number{Const: schema.NumberPtr(1.2)},
			`&schema.Number{
	Const: schema.NumberPtr(1.2),
}`,
		),
		Entry(
			"default",
			&schema.Number{Default: schema.NumberPtr(1.2)},
			`&schema.Number{
	Default: schema.NumberPtr(1.2),
}`,
		),
		Entry(
			"enum",
			&schema.Number{Enum: []float64{1.2}},
			`&schema.Number{
	Enum: []float64{1.2},
}`,
		),
		Entry(
			"minimum",
			&schema.Number{Minimum: schema.NumberPtr(1)},
			`&schema.Number{
	Minimum: schema.NumberPtr(1),
}`,
		),
		Entry(
			"maximum",
			&schema.Number{Maximum: schema.NumberPtr(1)},
			`&schema.Number{
	Maximum: schema.NumberPtr(1),
}`,
		),
		Entry(
			"exclusive minimum",
			&schema.Number{ExclusiveMinimum: schema.NumberPtr(1)},
			`&schema.Number{
	ExclusiveMinimum: schema.NumberPtr(1),
}`,
		),
		Entry(
			"exclusive maximum",
			&schema.Number{ExclusiveMaximum: schema.NumberPtr(1)},
			`&schema.Number{
	ExclusiveMaximum: schema.NumberPtr(1),
}`,
		),
		Entry(
			"multiple of",
			&schema.Number{MultipleOf: schema.NumberPtr(1)},
			`&schema.Number{
	MultipleOf: schema.NumberPtr(1),
}`,
		),
	)

	It("should marshal/unmarshal", func() {
		n := &schema.Number{
			Const:   schema.NumberPtr(1),
			Default: schema.NumberPtr(2),
			Enum:    []float64{3},
		}
		j, err := json.Marshal(n)
		Expect(err).ToNot(HaveOccurred())

		n2 := &schema.Number{}
		err = json.Unmarshal(j, &n2)
		Expect(err).ToNot(HaveOccurred())
		Expect(n2).To(Equal(n))
	})

	It("should unmarshal/marshal a json", func() {
		testhelper.ExpectConsistentJSONMarshalling(
			`{
				"const": 1.1,
				"default": 2.1,
				"enum": [3.1, 4.1],
				"exclusiveMinimum": 5.1,
				"exclusiveMaximum": 6.1,
				"maximum": 7.1,
				"minimum": 8.1,
				"multipleOf": 9.1,
				"type": "number"
			}`,
			&schema.Number{},
		)
	})
})
