// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package formats_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/src/schema"
	"github.com/conflowio/conflow/src/schema/formats"
)

var _ = Describe("Hostname", func() {

	format := formats.Hostname{}

	DescribeTable("Valid values",
		expectFormatToParse(format),
		Entry("simple host", "host", "host", "host"),
		Entry("subdomain", "sub.host", "sub.host", "sub.host"),
		Entry("starts with a digit", "123domain.com", "123domain.com", "123domain.com"),
		Entry("IP address", "1.2.3.4", "1.2.3.4", "1.2.3.4"),
	)

	DescribeTable("Invalid values",
		func(input string) {
			_, err := format.ValidateValue(input)
			Expect(err).To(HaveOccurred())
		},
		Entry("empty", ""),
		Entry("invalid character", "invalid!"),
		Entry("starts with dot", ".domain"),
		Entry("ends with dot", "domain."),
		Entry("has underscore", "domain_name"),
		Entry("part starts with hyphen", "foo.-bar"),
		Entry("part ends with hyphen", "foo-.bar"),
	)

	When("a field type is string and has 'hostname' format", func() {
		It("should be parsed as string schema with hostname format", func() {
			source := `
				// @block "configuration"
				type Foo struct {
					// @format "hostname"
					v string
				}
			`
			expectGoStructToHaveStringSchema(source, schema.FormatHostname, false)
		})
	})

	When("a field type is *string and has hostname format", func() {
		It("should be parsed as string schema with hostname format", func() {
			source := `
				// @block "configuration"
				type Foo struct {
					// @format "hostname"
					v *string
				}
			`
			expectGoStructToHaveStringSchema(source, schema.FormatHostname, true)
		})
	})

})
