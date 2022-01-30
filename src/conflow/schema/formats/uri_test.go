// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package formats_test

import (
	"net/url"

	"github.com/conflowio/conflow/src/conflow/schema"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/src/conflow/schema/formats"
)

var _ = Describe("URI", func() {

	format := formats.URI{RequireScheme: true}

	mustParse := func(u string) url.URL {
		res, err := url.Parse(u)
		if err != nil {
			panic(err)
		}
		return *res
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
			_, err := format.ValidateValue(input)
			Expect(err).To(HaveOccurred())
		},
		Entry("empty", ""),
		Entry("missing scheme", "domain/path?a=b#foo"),
	)

	When("a field type is url.URL", func() {
		It("should be parsed as string schema with uri format", func() {
			source := `
				import "net/url"
				// @block "configuration"
				type Foo struct {
					v url.URL
				}
			`
			expectGoStructToHaveStringSchema(source, schema.FormatURI, false)
		})
	})

	When("a field type is *url.URL", func() {
		It("should be parsed as string schema with uri format", func() {
			source := `
				import "net/url"
				// @block "configuration"
				type Foo struct {
					v *url.URL
				}
			`
			expectGoStructToHaveStringSchema(source, schema.FormatURI, true)
		})
	})

})

var _ = Describe("URIReference", func() {

	format := formats.URI{}

	mustParse := func(u string) url.URL {
		res, err := url.Parse(u)
		if err != nil {
			panic(err)
		}
		return *res
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
			"no schema",
			"//domain/path?a=b#foo",
			mustParse("//domain/path?a=b#foo"),
			"//domain/path?a=b#foo",
		),
		Entry(
			"absolute path",
			"/path?a=b#foo",
			mustParse("/path?a=b#foo"),
			"/path?a=b#foo",
		),
		Entry(
			"relative path",
			"path?a=b#foo",
			mustParse("path?a=b#foo"),
			"path?a=b#foo",
		),
		Entry(
			"empty",
			"",
			mustParse(""),
			"",
		),
	)

	DescribeTable("Invalid values",
		func(input string) {
			_, err := format.ValidateValue(input)
			Expect(err).To(HaveOccurred())
		},
		Entry("not a url", `://@fooo-/lala`),
	)

	When("a field type is url.URL with format uri-reference", func() {
		It("should be parsed as string schema with uri-reference format", func() {
			source := `
				import "net/url"
				// @block "configuration"
				type Foo struct {
					// @format "uri-reference"
					v url.URL
				}
			`
			expectGoStructToHaveStringSchema(source, schema.FormatURIReference, false)
		})
	})

})
