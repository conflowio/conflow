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

var _ schema.Schema = &schema.Number{}

var _ = Describe("Number", func() {
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
})
