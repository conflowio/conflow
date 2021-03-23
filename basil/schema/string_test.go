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
	)

	It("should marshal/unmarshal", func() {
		s := &schema.String{
			Const:   schema.StringPtr("constval"),
			Default: schema.StringPtr("defaultval"),
			Enum:    []string{"enum1", "enum2"},
			Format:  "formatval",
		}
		j, err := json.Marshal(s)
		Expect(err).ToNot(HaveOccurred())

		s2 := &schema.String{}
		err = json.Unmarshal(j, &s2)
		Expect(err).ToNot(HaveOccurred())
		Expect(s2).To(Equal(s))
	})
})
