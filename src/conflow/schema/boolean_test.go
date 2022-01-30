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

var _ schema.Schema = &schema.Boolean{}

var _ = Describe("Boolean", func() {
	DescribeTable("Validate accepts value",
		func(schema *schema.Boolean, value interface{}) {
			_, err := schema.ValidateValue(value)
			Expect(err).ToNot(HaveOccurred())
		},
		Entry("true", &schema.Boolean{}, true),
		Entry("false", &schema.Boolean{}, false),
		Entry("const value", &schema.Boolean{Const: schema.BooleanPtr(true)}, true),
		Entry("enum value - single", &schema.Boolean{Enum: []bool{true}}, true),
	)

	DescribeTable("Validate errors",
		func(schema *schema.Boolean, value interface{}, expectedErr error) {
			_, err := schema.ValidateValue(value)
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
			str := schema.GoString(map[string]string{})
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
		Entry(
			"nullable",
			&schema.Boolean{Nullable: true},
			`&schema.Boolean{
	Nullable: true,
}`,
		),
	)

	It("should unmarshal/marshal a json", func() {
		testhelper.ExpectConsistentJSONMarshalling(
			`{
				"const": true,
				"default": false,
				"enum": [false, true],
				"nullable": true,
				"type": "boolean"
			}`,
			&schema.Boolean{},
		)
	})
})
