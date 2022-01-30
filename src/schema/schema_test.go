// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema_test

import (
	"encoding/json"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/conflowio/conflow/src/conflow/types"

	"github.com/google/uuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/src/schema"
)

var _ = Describe("Metadata", func() {
	DescribeTable("GoString prints a valid Go struct",
		func(b schema.Metadata, expected string) {
			str := b.GoString()
			Expect(str).To(Equal(expected))
		},
		Entry(
			"empty",
			schema.Metadata{},
			`schema.Metadata{
}`,
		),
		Entry(
			"annotations",
			schema.Metadata{Annotations: map[string]string{"foo": "bar"}},
			`schema.Metadata{
	Annotations: map[string]string{"foo":"bar"},
}`,
		),
		Entry(
			"description",
			schema.Metadata{Description: "foo"},
			`schema.Metadata{
	Description: "foo",
}`,
		),
		Entry(
			"examples",
			schema.Metadata{Examples: []interface{}{"foo"}},
			`schema.Metadata{
	Examples: []interface {}{"foo"},
}`,
		),
		Entry(
			"readonly",
			schema.Metadata{ReadOnly: true},
			`schema.Metadata{
	ReadOnly: true,
}`,
		),
		Entry(
			"title",
			schema.Metadata{Title: "foo"},
			`schema.Metadata{
	Title: "foo",
}`,
		),
		Entry(
			"writeonly",
			schema.Metadata{WriteOnly: true},
			`schema.Metadata{
	WriteOnly: true,
}`,
		),
	)
})

var _ = Describe("Schema", func() {
	It("should unmarshal an empty schema into an untyped schema", func() {
		s, err := schema.UnmarshalJSON([]byte("{}"))
		Expect(err).ToNot(HaveOccurred())
		Expect(s).To(Equal(&schema.Untyped{}))
	})

	It("should unmarshal an empty schema into an untyped schema - whitespaces", func() {
		s, err := schema.UnmarshalJSON([]byte("{\n}"))
		Expect(err).ToNot(HaveOccurred())
		Expect(s).To(Equal(&schema.Untyped{}))
	})

	DescribeTable(
		"should marshal/unmarshal",
		func(s schema.Schema) {
			j, err := json.Marshal(s)
			Expect(err).ToNot(HaveOccurred())

			s2, err := schema.UnmarshalJSON(j)
			Expect(err).ToNot(HaveOccurred())
			Expect(s2).To(Equal(s))
		},
		Entry(
			"boolean",
			&schema.Boolean{
				Metadata: schema.Metadata{Description: "foo"},
			},
		),
		Entry(
			"integer",
			&schema.Integer{
				Metadata: schema.Metadata{Description: "foo"},
			},
		),
		Entry(
			"number",
			&schema.Number{
				Metadata: schema.Metadata{Description: "foo"},
			},
		),
		Entry(
			"string",
			&schema.String{
				Metadata: schema.Metadata{Description: "foo"},
			},
		),
		Entry(
			"untyped",
			&schema.Untyped{
				Metadata: schema.Metadata{Description: "foo"},
			},
		),
		Entry(
			"untyped with multiple types",
			&schema.Untyped{
				Metadata: schema.Metadata{Description: "foo"},
				Types:    []string{string(schema.TypeString), string(schema.TypeBoolean)},
			},
		),
		Entry(
			"null",
			&schema.Null{
				Metadata: schema.Metadata{Description: "foo"},
			},
		),
		Entry(
			"false",
			schema.False(),
		),
		Entry(
			"array",
			&schema.Array{
				Metadata: schema.Metadata{Description: "foo"},
				Items:    &schema.String{},
			},
		),
		Entry(
			"map",
			&schema.Map{
				Metadata:             schema.Metadata{Description: "foo"},
				AdditionalProperties: &schema.String{},
			},
		),
		Entry(
			"object",
			&schema.Object{
				Metadata: schema.Metadata{Description: "foo"},
				Parameters: map[string]schema.Schema{
					"foo": &schema.Integer{},
				},
			},
		),
	)

	DescribeTable("GetSchemaForValue",

		func(v interface{}, expected schema.Schema, expectedErr error) {
			s, err := schema.GetSchemaForValue(v)
			if expectedErr == nil {
				Expect(err).ToNot(HaveOccurred())
				Expect(s).To(Equal(expected))
			} else {
				Expect(err).To(MatchError(expectedErr))
			}
		},
		Entry("null", nil, schema.NullValue(), nil),
		Entry("boolean", false, schema.BooleanValue(), nil),
		Entry("integer", int64(1), schema.IntegerValue(), nil),
		Entry("number", float64(1), schema.NumberValue(), nil),
		Entry("string", "foo", schema.StringValue(), nil),
		Entry("byte stream", os.Stdin, schema.ByteStreamValue(), nil),
		Entry("unknown value", int8(1), nil, fmt.Errorf("value type int8 is not allowed")),

		Entry("byte", []byte("aGVsbG8="), &schema.String{Format: schema.FormatBinary}, nil),
		Entry("date-time", time.Now(), &schema.String{Format: schema.FormatDateTime}, nil),
		Entry("duration", types.RFC3339Duration{}, &schema.String{Format: schema.FormatDurationRFC3339}, nil),
		Entry("duration-go", 1*time.Second, &schema.String{Format: schema.FormatDurationGo}, nil),
		Entry("email", mail.Address{}, &schema.String{Format: schema.FormatEmail}, nil),
		Entry("ip", net.IP{}, &schema.String{Format: schema.FormatIP}, nil),
		Entry("ipc-cidr", types.CIDR{}, &schema.String{Format: schema.FormatIPCIDR}, nil),
		Entry("regexp", regexp.Regexp{}, &schema.String{Format: schema.FormatRegex}, nil),
		Entry("time", types.Time{}, &schema.String{Format: schema.FormatTime}, nil),
		Entry("uri", url.URL{}, &schema.String{Format: schema.FormatURI}, nil),
		Entry("uuid", uuid.UUID{}, &schema.String{Format: schema.FormatUUID}, nil),

		Entry(
			"bool array",
			[]interface{}{false, true},
			&schema.Array{Items: schema.BooleanValue()}, nil,
		),
		Entry(
			"integer array",
			[]interface{}{int64(1), int64(2)},
			&schema.Array{Items: schema.IntegerValue()}, nil,
		),
		Entry(
			"number array",
			[]interface{}{float64(1), float64(2)},
			&schema.Array{Items: schema.NumberValue()}, nil,
		),
		Entry(
			"mixed numeric array",
			[]interface{}{int64(1), float64(2)},
			&schema.Array{Items: schema.NumberValue()}, nil,
		),
		Entry(
			"mixed numeric array",
			[]interface{}{float64(1), int64(2)},
			&schema.Array{Items: schema.NumberValue()}, nil,
		),
		Entry(
			"string array",
			[]interface{}{"a", "b"},
			&schema.Array{Items: schema.StringValue()}, nil,
		),
		Entry(
			"byte stream array",
			[]interface{}{os.Stdin, os.Stderr},
			&schema.Array{Items: schema.ByteStreamValue()}, nil,
		),
		Entry(
			"time array",
			[]interface{}{time.Now(), time.Now()},
			&schema.Array{Items: &schema.String{Format: schema.FormatDateTime}}, nil,
		),
		Entry(
			"time duration array",
			[]interface{}{1 * time.Second, 2 * time.Second},
			&schema.Array{Items: &schema.String{Format: schema.FormatDurationGo}}, nil,
		),
		Entry(
			"mixed types in array",
			[]interface{}{int64(1), "foo"},
			nil, fmt.Errorf("items must have the same type, but found integer and string"),
		),
		Entry(
			"unknown value in array",
			[]interface{}{int64(1), int8(2)},
			nil, fmt.Errorf("value type int8 is not allowed"),
		),

		Entry(
			"bool map",
			map[string]interface{}{"a": false, "b": true},
			&schema.Map{AdditionalProperties: schema.BooleanValue()}, nil,
		),
		Entry(
			"integer map",
			map[string]interface{}{"a": int64(1), "b": int64(2)},
			&schema.Map{AdditionalProperties: schema.IntegerValue()}, nil,
		),
		Entry(
			"number map",
			map[string]interface{}{"a": float64(1), "b": float64(2)},
			&schema.Map{AdditionalProperties: schema.NumberValue()}, nil,
		),
		Entry(
			"mixed numeric map",
			map[string]interface{}{"a": int64(1), "b": float64(2)},
			&schema.Map{AdditionalProperties: schema.NumberValue()}, nil,
		),
		Entry(
			"mixed numeric map",
			map[string]interface{}{"a": float64(1), "b": int64(2)},
			&schema.Map{AdditionalProperties: schema.NumberValue()}, nil,
		),
		Entry(
			"string map",
			map[string]interface{}{"a": "a", "b": "b"},
			&schema.Map{AdditionalProperties: schema.StringValue()}, nil,
		),
		Entry(
			"byte stream map",
			map[string]interface{}{"a": os.Stdin, "b": os.Stderr},
			&schema.Map{AdditionalProperties: schema.ByteStreamValue()}, nil,
		),
		Entry(
			"time map",
			map[string]interface{}{"a": time.Now(), "b": time.Now()},
			&schema.Map{AdditionalProperties: &schema.String{Format: schema.FormatDateTime}}, nil,
		),
		Entry(
			"time duration map",
			map[string]interface{}{"a": 1 * time.Second, "b": 2 * time.Second},
			&schema.Map{AdditionalProperties: &schema.String{Format: schema.FormatDurationGo}}, nil,
		),
		Entry(
			"mixed types in map",
			map[string]interface{}{"a": int64(1), "b": "foo"},
			nil, fmt.Errorf("items must have the same type, but found integer and string"),
		),
		Entry(
			"unknown value in map",
			map[string]interface{}{"a": int64(1), "b": int8(2)},
			nil, fmt.Errorf("value type int8 is not allowed"),
		),

		Entry(
			"array of arrays",
			[]interface{}{
				[]interface{}{"a", "b"},
				[]interface{}{"c", "d"},
			},
			&schema.Array{Items: &schema.Array{Items: schema.StringValue()}}, nil,
		),

		Entry(
			"array of arrays with mixed numeric values",
			[]interface{}{
				[]interface{}{int64(1), int64(2)},
				[]interface{}{float64(3), float64(4)},
			},
			&schema.Array{Items: &schema.Array{Items: schema.NumberValue()}}, nil,
		),

		Entry(
			"maps of maps",
			map[string]interface{}{
				"a": map[string]interface{}{"a": "a", "b": "b"},
				"b": map[string]interface{}{"c": "c", "d": "d"},
			},
			&schema.Map{AdditionalProperties: &schema.Map{AdditionalProperties: schema.StringValue()}}, nil,
		),

		Entry(
			"maps of maps with mixed numeric values",
			map[string]interface{}{
				"a": map[string]interface{}{"a": int64(1), "b": int64(2)},
				"b": map[string]interface{}{"c": float64(3), "d": float64(4)},
			},
			&schema.Map{AdditionalProperties: &schema.Map{AdditionalProperties: schema.NumberValue()}}, nil,
		),
	)
})
