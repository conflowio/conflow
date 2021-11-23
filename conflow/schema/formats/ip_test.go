// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package formats_test

import (
	"net"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/conflow/schema"

	"github.com/conflowio/conflow/conflow/schema/formats"
)

var _ = Describe("IP", func() {

	format := formats.IP{AllowIPv4: true, AllowIPv6: true}

	ipv4 := net.ParseIP("1.2.3.4")
	ipv6 := net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334")

	DescribeTable("Valid values",
		expectFormatToParse(format),
		Entry("valid IPv4", "1.2.3.4", ipv4, "1.2.3.4"),
		Entry("valid IPv6", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", ipv6, "2001:db8:85a3::8a2e:370:7334"),
		Entry("IPv6 shortening", "0:0:0:0:0:0:0:1", net.IPv6loopback, "::1"),
	)

	DescribeTable("Invalid values",
		func(input string) {
			_, err := format.ValidateValue(input)
			Expect(err).To(HaveOccurred())
		},
		Entry("empty", ""),
		Entry("random", "foo"),
		Entry("incomplete IPv4", "1.2.3"),
		Entry("incomplete IPv6", "2001:0db8:85a3:0000:0000:8a2e:0370"),
	)

	When("a field type is net.IP", func() {
		It("should be parsed as string schema with ip format", func() {
			source := `
				import "net"
				// @block "configuration"
				type Foo struct {
					v net.IP
				}
			`
			expectGoStructToHaveStringSchema(source, schema.FormatIP, false)
		})
	})

	When("a field type is *net.IP", func() {
		It("should be parsed as string schema with ip format", func() {
			source := `
				import "net"
				// @block "configuration"
				type Foo struct {
					v *net.IP
				}
			`
			expectGoStructToHaveStringSchema(source, schema.FormatIP, true)
		})
	})

})

var _ = Describe("IPv4", func() {

	format := formats.IP{AllowIPv4: true}

	ipv4 := net.ParseIP("1.2.3.4")

	DescribeTable("Valid values",
		expectFormatToParse(format),
		Entry("valid IPv4", "1.2.3.4", ipv4, "1.2.3.4"),
	)

	DescribeTable("Invalid values",
		func(input string) {
			_, err := format.ValidateValue(input)
			Expect(err).To(HaveOccurred())
		},
		Entry("empty", ""),
		Entry("random", "foo"),
		Entry("incomplete IPv4", "1.2.3"),
		Entry("IPv6", "2001:0db8:85a3:0000:0000:8a2e:0370:7334"),
	)

	When("a field type is net.IP with format ipv4", func() {
		It("should be parsed as string schema with ipv4 format", func() {
			source := `
				import "net"
				// @block "configuration"
				type Foo struct {
					// @format "ipv4"
					v net.IP
				}
			`
			expectGoStructToHaveStringSchema(source, schema.FormatIPv4, false)
		})
	})

})

var _ = Describe("IPv6", func() {

	format := formats.IP{AllowIPv6: true}

	ipv6 := net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334")

	DescribeTable("Valid values",
		expectFormatToParse(format),
		Entry("valid IPv6", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", ipv6, "2001:db8:85a3::8a2e:370:7334"),
		Entry("IPv6 shortening", "0:0:0:0:0:0:0:1", net.IPv6loopback, "::1"),
	)

	DescribeTable("Invalid values",
		func(input string) {
			_, err := format.ValidateValue(input)
			Expect(err).To(HaveOccurred())
		},
		Entry("empty", ""),
		Entry("random", "foo"),
		Entry("incomplete IPv6", "2001:0db8:85a3:0000:0000:8a2e:0370"),
		Entry("IPv4", "1.2.3.4"),
	)

	When("a field type is net.IP with format ipv6", func() {
		It("should be parsed as string schema with ipv6 format", func() {
			source := `
				import "net"
				// @block "configuration"
				type Foo struct {
					// @format "ipv6"
					v net.IP
				}
			`
			expectGoStructToHaveStringSchema(source, schema.FormatIPv6, false)
		})
	})

})
