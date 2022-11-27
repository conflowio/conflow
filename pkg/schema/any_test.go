// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/pkg/internal/testhelper"
	"github.com/conflowio/conflow/pkg/schema"
)

var _ schema.Schema = &schema.Any{}

var _ = Describe("Any", func() {
	DescribeTable("Validate accepts value",
		func(schema *schema.Any, value interface{}) {
			_, err := schema.ValidateValue(value)
			Expect(err).ToNot(HaveOccurred())
		},
		Entry("accepts anything", &schema.Any{}, "foo"),
		Entry("accepts pointer", &schema.Any{}, schema.Pointer("foo")),
		Entry("const value", &schema.Any{Const: "foo"}, "foo"),
		Entry("const value - pointer", &schema.Any{Const: "foo"}, schema.Pointer("foo")),
		Entry("enum value", &schema.Any{Enum: []interface{}{"foo"}}, "foo"),
		Entry("enum value - pointer", &schema.Any{Enum: []interface{}{"foo"}}, schema.Pointer("foo")),
		Entry("nil is not validated", &schema.Any{Const: "foo"}, nil),
		Entry("nil pointer is not validated", &schema.Any{Const: "foo"}, (*string)(nil)),
	)

	DescribeTable("Validate errors",
		func(schema *schema.Any, value interface{}, expectedErr error) {
			_, err := schema.ValidateValue(value)
			Expect(err).To(MatchError(expectedErr))
		},
		Entry(
			"const value",
			&schema.Any{Const: "foo"},
			"bar",
			errors.New("must be foo"),
		),
		Entry(
			"enum value - single",
			&schema.Any{Enum: []interface{}{"foo"}},
			"bar",
			errors.New("must be foo"),
		),
		Entry(
			"enum value - multiple",
			&schema.Any{Enum: []interface{}{"foo", int64(1)}},
			"bar",
			errors.New("must be one of foo, 1"),
		),
	)

	DescribeTable("GoString prints a valid Go struct",
		func(schema *schema.Any, expected string) {
			str := schema.GoString(map[string]string{})
			Expect(str).To(Equal(expected))
		},
		Entry(
			"empty",
			&schema.Any{},
			`&schema.Any{
}`,
		),
		Entry(
			"const",
			&schema.Any{Const: true},
			`&schema.Any{
	Const: true,
}`,
		),
		Entry(
			"default",
			&schema.Any{Default: true},
			`&schema.Any{
	Default: true,
}`,
		),
		Entry(
			"enum",
			&schema.Any{Enum: []interface{}{true}},
			`&schema.Any{
	Enum: []interface {}{true},
}`,
		),
		Entry(
			"nullable",
			&schema.Any{Nullable: true},
			`&schema.Any{
	Nullable: true,
}`,
		),
	)

	It("should unmarshal/marshal a json", func() {
		testhelper.ExpectConsistentJSONMarshalling(
			`{
				"const": true,
				"default": false,
				"enum": [true, 1],
				"nullable": true
			}`,
			&schema.Any{},
		)
	})
})
