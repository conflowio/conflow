// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package formats_test

import (
	"github.com/conflowio/conflow/src/schema"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/src/schema/formats"
)

var _ = Describe("Byte", func() {

	format := formats.Binary{}

	DescribeTable("Valid values",
		expectFormatToParse(format),
		Entry(
			"empty string",
			"",
			[]byte(""),
			"",
		),
		Entry(
			"'hello' base64 encoded",
			"aGVsbG8=",
			[]byte("hello"),
			"aGVsbG8=",
		),
	)

	DescribeTable("Invalid values",
		func(input string) {
			_, err := format.ValidateValue(input)
			Expect(err).To(HaveOccurred())
		},
		Entry("not base64", "foo"),
	)

	When("a field type is []byte", func() {
		It("should be parsed as string schema with byte format", func() {
			source := `
				// @block "configuration"
				type Foo struct {
					v []byte
				}
			`
			expectGoStructToHaveStringSchema(source, schema.FormatBinary, false)
		})
	})

	When("a field type is *[]byte", func() {
		It("should be parsed as string schema with byte format", func() {
			source := `
				// @block "configuration"
				type Foo struct {
					v *[]byte
				}
			`
			expectGoStructToHaveStringSchema(source, schema.FormatBinary, true)
		})
	})

})
