// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema_test

import (
	"encoding/json"
	"errors"

	"github.com/conflowio/conflow/conflow/schema"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ schema.Schema = &schema.Boolean{}

var _ = Describe("Boolean", func() {
	DescribeTable("Validate accepts value",
		func(schema *schema.Boolean, value interface{}) {
			err := schema.ValidateValue(value)
			Expect(err).ToNot(HaveOccurred())
		},
		Entry("true", &schema.Boolean{}, true),
		Entry("false", &schema.Boolean{}, false),
		Entry("const value", &schema.Boolean{Const: schema.BooleanPtr(true)}, true),
		Entry("enum value - single", &schema.Boolean{Enum: []bool{true}}, true),
	)

	DescribeTable("Validate errors",
		func(schema *schema.Boolean, value interface{}, expectedErr error) {
			err := schema.ValidateValue(value)
			Expect(err).To(MatchError(expectedErr))
		},
		Entry(
			"non-boolean value",
			&schema.Boolean{},
			"foo",
			errors.New("must be boolean"),
		),
		Entry(
			"const value",
			&schema.Boolean{Const: schema.BooleanPtr(true)},
			false,
			errors.New("must be true"),
		),
		Entry(
			"enum value - single",
			&schema.Boolean{Enum: []bool{true}},
			false,
			errors.New("must be true"),
		),
	)

	DescribeTable("GoString prints a valid Go struct",
		func(schema *schema.Boolean, expected string) {
			str := schema.GoString()
			Expect(str).To(Equal(expected))
		},
		Entry(
			"empty",
			&schema.Boolean{},
			`&schema.Boolean{
}`,
		),
		Entry(
			"const",
			&schema.Boolean{Const: schema.BooleanPtr(true)},
			`&schema.Boolean{
	Const: schema.BooleanPtr(true),
}`,
		),
		Entry(
			"default",
			&schema.Boolean{Default: schema.BooleanPtr(true)},
			`&schema.Boolean{
	Default: schema.BooleanPtr(true),
}`,
		),
		Entry(
			"enum",
			&schema.Boolean{Enum: []bool{true}},
			`&schema.Boolean{
	Enum: []bool{true},
}`,
		),
	)

	It("should marshal/unmarshal", func() {
		b := &schema.Boolean{
			Const:   schema.BooleanPtr(true),
			Default: schema.BooleanPtr(true),
			Enum:    []bool{true},
		}
		j, err := json.Marshal(b)
		Expect(err).ToNot(HaveOccurred())

		b2 := &schema.Boolean{}
		err = json.Unmarshal(j, &b2)
		Expect(err).ToNot(HaveOccurred())
		Expect(b2).To(Equal(b))
	})
})
