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

var _ schema.Schema = &schema.OneOf{}

var _ = Describe("OneOf", func() {
	DescribeTable("Validate accepts value",
		func(schema *schema.OneOf, value interface{}) {
			_, err := schema.ValidateValue(value)
			Expect(err).ToNot(HaveOccurred())
		},
		Entry("accepts boolean", &schema.OneOf{Schemas: []schema.Schema{&schema.Boolean{}, &schema.Integer{}}}, true),
		Entry("accepts integer", &schema.OneOf{Schemas: []schema.Schema{&schema.Boolean{}, &schema.Integer{}}}, int64(1)),
		Entry("accepts pointer", &schema.OneOf{Schemas: []schema.Schema{&schema.Boolean{}, &schema.Integer{}}}, schema.Pointer(int64(1))),
		Entry("const value - boolean", &schema.OneOf{Schemas: []schema.Schema{&schema.Boolean{}, &schema.Integer{}}, Const: true}, true),
		Entry("const value - integer", &schema.OneOf{Schemas: []schema.Schema{&schema.Boolean{}, &schema.Integer{}}, Const: int64(1)}, int64(1)),
		Entry("enum value - boolean", &schema.OneOf{Schemas: []schema.Schema{&schema.Boolean{}, &schema.Integer{}}, Enum: []interface{}{true, int64(1)}}, true),
		Entry("enum value - integer", &schema.OneOf{Schemas: []schema.Schema{&schema.Boolean{}, &schema.Integer{}}, Enum: []interface{}{true, int64(1)}}, int64(1)),
		Entry("nil is not validated", &schema.OneOf{Schemas: []schema.Schema{&schema.Boolean{}, &schema.Integer{}}, Const: true}, nil),
		Entry("nil pointer is not validated", &schema.OneOf{Schemas: []schema.Schema{&schema.Boolean{}, &schema.Integer{}}, Const: true}, (*string)(nil)),
	)

	DescribeTable("Validate errors",
		func(schema *schema.OneOf, value interface{}, expectedErr error) {
			_, err := schema.ValidateValue(value)
			Expect(err).To(MatchError(expectedErr))
		},
		Entry(
			"invalid value",
			&schema.OneOf{Schemas: []schema.Schema{&schema.Boolean{}, &schema.Integer{}}},
			"foo",
			errors.New("must be boolean or integer"),
		),
		Entry(
			"const value",
			&schema.OneOf{Schemas: []schema.Schema{&schema.Boolean{}, &schema.Integer{}}, Const: true},
			int64(1),
			errors.New("must be true"),
		),
		Entry(
			"enum value - single",
			&schema.OneOf{Schemas: []schema.Schema{&schema.Boolean{}, &schema.Integer{}}, Enum: []interface{}{true}},
			false,
			errors.New("must be true"),
		),
		Entry(
			"enum value - multiple",
			&schema.OneOf{Schemas: []schema.Schema{&schema.Boolean{}, &schema.Integer{}}, Enum: []interface{}{true, int64(1)}},
			"bar",
			errors.New("must be one of true, 1"),
		),
	)

	DescribeTable("GoString prints a valid Go struct",
		func(schema *schema.OneOf, expected string) {
			str := schema.GoString(map[string]string{})
			Expect(str).To(Equal(expected))
		},
		Entry(
			"empty",
			&schema.OneOf{},
			`&schema.OneOf{
}`,
		),
		Entry(
			"const",
			&schema.OneOf{Const: true},
			`&schema.OneOf{
	Const: true,
}`,
		),
		Entry(
			"default",
			&schema.OneOf{Default: true},
			`&schema.OneOf{
	Default: true,
}`,
		),
		Entry(
			"enum",
			&schema.OneOf{Enum: []interface{}{true}},
			`&schema.OneOf{
	Enum: []interface {}{true},
}`,
		),
		Entry(
			"nullable",
			&schema.OneOf{Nullable: true},
			`&schema.OneOf{
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
				"nullable": true,
				"oneOf": [
					{
						"type": "boolean",
						"title": "t1"
					},
					{
						"type": "integer",
						"title": "t1"
					}
				]
			}`,
			&schema.OneOf{},
		)
	})
})
