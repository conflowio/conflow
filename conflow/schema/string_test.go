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

var _ schema.Schema = &schema.String{}

var _ = Describe("String", func() {
	DescribeTable("Validate accepts value",
		func(schema *schema.String, value interface{}) {
			err := schema.ValidateValue(value)
			Expect(err).ToNot(HaveOccurred())
		},
		Entry("empty string", &schema.String{}, ""),
		Entry("non-empty string", &schema.String{}, "foo"),
		Entry("const value", &schema.String{Const: schema.StringPtr("foo")}, "foo"),
		Entry("enum value - single", &schema.String{Enum: []string{"foo"}}, "foo"),
		Entry("enum value - multiple", &schema.String{Enum: []string{"foo", "bar"}}, "foo"),
		Entry("min length - equal", &schema.String{MinLength: 1}, "a"),
		Entry("min length - longer", &schema.String{MinLength: 1}, "ab"),
		Entry("min length - unicode", &schema.String{MinLength: 1}, "🍕"),
		Entry("max length - empty", &schema.String{MaxLength: schema.IntegerPtr(0)}, ""),
		Entry("max length - equal", &schema.String{MaxLength: schema.IntegerPtr(1)}, "a"),
		Entry("max length - shorter", &schema.String{MaxLength: schema.IntegerPtr(2)}, "a"),
		Entry("max length - unicode", &schema.String{MaxLength: schema.IntegerPtr(1)}, "🍕"),
	)

	DescribeTable("Validate errors",
		func(schema *schema.String, value interface{}, expectedErr error) {
			err := schema.ValidateValue(value)
			Expect(err).To(MatchError(expectedErr))
		},
		Entry(
			"non-string value",
			&schema.String{},
			1,
			errors.New("must be string"),
		),
		Entry(
			"const value",
			&schema.String{Const: schema.StringPtr("foo")},
			"bar",
			errors.New(`must be "foo"`),
		),
		Entry(
			"enum value - single",
			&schema.String{Enum: []string{"foo"}},
			"bar",
			errors.New(`must be "foo"`),
		),
		Entry(
			"enum value - multiple",
			&schema.String{Enum: []string{"foo", "bar"}},
			"baz",
			errors.New(`must be one of "foo", "bar"`),
		),
		Entry(
			"min length - empty",
			&schema.String{MinLength: 1},
			"",
			errors.New(`can not be empty string`),
		),
		Entry(
			"min length - shorter",
			&schema.String{MinLength: 2},
			"a",
			errors.New(`must be at least 2 characters long`),
		),
		Entry(
			"min length - unicode",
			&schema.String{MinLength: 2},
			"🍕",
			errors.New(`must be at least 2 characters long`),
		),
		Entry(
			"max length - empty",
			&schema.String{MaxLength: schema.IntegerPtr(0)},
			"a",
			errors.New(`must be empty string`),
		),
		Entry(
			"max length - 1",
			&schema.String{MaxLength: schema.IntegerPtr(1)},
			"ab",
			errors.New(`must be empty string or a single character`),
		),
		Entry(
			"max length - 2",
			&schema.String{MaxLength: schema.IntegerPtr(2)},
			"abc",
			errors.New(`must be no more than 2 characters long`),
		),
		Entry(
			"min length = max length - 1",
			&schema.String{MinLength: 1, MaxLength: schema.IntegerPtr(1)},
			"ab",
			errors.New(`must be a single character`),
		),
		Entry(
			"min length = max length - 2",
			&schema.String{MinLength: 2, MaxLength: schema.IntegerPtr(2)},
			"abc",
			errors.New(`must be exactly 2 characters long`),
		),
	)

	DescribeTable("GoString prints a valid Go struct",
		func(schema *schema.String, expected string) {
			str := schema.GoString()
			Expect(str).To(Equal(expected))
		},
		Entry(
			"empty",
			&schema.String{},
			`&schema.String{
}`,
		),
		Entry(
			"const",
			&schema.String{Const: schema.StringPtr("foo")},
			`&schema.String{
	Const: schema.StringPtr("foo"),
}`,
		),
		Entry(
			"default",
			&schema.String{Default: schema.StringPtr("foo")},
			`&schema.String{
	Default: schema.StringPtr("foo"),
}`,
		),
		Entry(
			"enum",
			&schema.String{Enum: []string{"foo"}},
			`&schema.String{
	Enum: []string{"foo"},
}`,
		),
		Entry(
			"format",
			&schema.String{Format: "foo"},
			`&schema.String{
	Format: "foo",
}`,
		),
		Entry(
			"min length",
			&schema.String{MinLength: 1},
			`&schema.String{
	MinLength: 1,
}`,
		),
		Entry(
			"max length",
			&schema.String{MaxLength: schema.IntegerPtr(1)},
			`&schema.String{
	MaxLength: schema.IntegerPtr(1),
}`,
		),
	)

	It("should unmarshal/marshal a json", func() {
		testhelper.ExpectConsistentJSONMarshalling(
			`{
				"const": "constval",
				"default": "defaultval",
				"enum": ["enum1", "enum2"],
				"format": "formatval",
				"minLength": 1,
				"maxLength": 2,
				"type": "string"
			}`,
			&schema.String{},
		)
	})
})
