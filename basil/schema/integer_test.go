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
	)

	It("should marshal/unmarshal", func() {
		i := &schema.Integer{
			Const:   schema.IntegerPtr(1),
			Default: schema.IntegerPtr(2),
			Enum:    []int64{3},
		}
		j, err := json.Marshal(i)
		Expect(err).ToNot(HaveOccurred())

		i2 := &schema.Integer{}
		err = json.Unmarshal(j, &i2)
		Expect(err).ToNot(HaveOccurred())
		Expect(i2).To(Equal(i))
	})
})
