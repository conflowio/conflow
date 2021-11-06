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

	"github.com/conflowio/conflow/conflow/schema/formats"
)

var _ = Describe("IPv4", func() {

	format := formats.IPv4{}

	ipv4 := net.ParseIP("1.2.3.4")

	DescribeTable("Valid values",
		expectFormatToParse(format),
		Entry("valid IPv4", "1.2.3.4", &ipv4, "1.2.3.4"),
	)

	DescribeTable("Invalid values",
		func(input string) {
			_, err := format.Parse(input)
			Expect(err).To(HaveOccurred())
		},
		Entry("empty", ""),
		Entry("random", "foo"),
		Entry("incomplete IPv4", "1.2.3"),
		Entry("IPv6", "2001:0db8:85a3:0000:0000:8a2e:0370:7334"),
	)

})
