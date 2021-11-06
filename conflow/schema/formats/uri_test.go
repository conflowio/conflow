// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package formats_test

import (
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/conflow/schema/formats"
)

var _ = Describe("URI", func() {

	format := formats.URI{}

	mustParse := func(u string) *url.URL {
		res, err := url.Parse(u)
		if err != nil {
			panic(err)
		}
		return res
	}

	DescribeTable("Valid values",
		expectFormatToParse(format),
		Entry(
			"valid URI",
			"http://domain/path?a=b#foo",
			mustParse("http://domain/path?a=b#foo"),
			"http://domain/path?a=b#foo",
		),
		Entry(
			"URI containing unsafe characters",
			"http://domain/my \\path",
			mustParse("http://domain/my \\path"),
			"http://domain/my%20%5Cpath",
		),
	)

	DescribeTable("Invalid values",
		func(input string) {
			_, err := format.Parse(input)
			Expect(err).To(HaveOccurred())
		},
		Entry("empty", ""),
		Entry("missing scheme", "domain/path?a=b#foo"),
	)

})
