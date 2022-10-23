// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package formats_test

import (
	"regexp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/pkg/conflow/types"
	"github.com/conflowio/conflow/pkg/schema"
	"github.com/conflowio/conflow/pkg/schema/formats"
)

var _ = Describe("Regex", func() {

	format := formats.Regex{}

	re := regexp.MustCompile("^[a-z]+$")

	DescribeTable("Valid values",
		expectFormatToParse[types.Regexp](format),
		Entry("valid regexp", "^[a-z]+$", (types.Regexp)(*re), "^[a-z]+$"),
	)

	DescribeTable("Invalid values",
		func(input string) {
			_, err := format.ValidateValue(input)
			Expect(err).To(HaveOccurred())
		},
		Entry("empty", ""),
		Entry("missing parentheses", "(a-z"),
	)

	When("a field type is types.Regexp", func() {
		It("should be parsed as string schema with regex format", func() {
			source := `
				import "github.com/conflowio/conflow/pkg/conflow/types"
				// @block "configuration"
				type Foo struct {
					v types.Regexp
				}
			`
			expectGoStructToHaveStringSchema(source, schema.FormatRegex, false)
		})
	})

	When("a field type is *types.Regexp", func() {
		It("should be parsed as string schema with regex format", func() {
			source := `
				import "github.com/conflowio/conflow/pkg/conflow/types"
				// @block "configuration"
				type Foo struct {
					v *types.Regexp
				}
			`
			expectGoStructToHaveStringSchema(source, schema.FormatRegex, true)
		})
	})

	It("should have a consistent JSON marshalling", func() {
		expectConsistentJSONMarshalling[*types.Regexp]([]byte("null"))
	})

})
