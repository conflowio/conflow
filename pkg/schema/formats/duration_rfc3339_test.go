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

	"github.com/conflowio/conflow/pkg/conflow/types"
	"github.com/conflowio/conflow/pkg/schema"
	"github.com/conflowio/conflow/pkg/schema/formats"
)

var _ = Describe("DurationRFC3339", func() {

	format := formats.DurationRFC3339{}

	DescribeTable("Valid values",
		expectFormatToParse[types.RFC3339Duration](format),

		Entry(
			"day",
			"P13D",
			types.RFC3339Duration{Day: 13},
			"P13D",
		),
		Entry(
			"month + day",
			"P12M13D",
			types.RFC3339Duration{Month: 12, Day: 13},
			"P12M13D",
		),
		Entry(
			"year + month + day",
			"P11Y12M13D",
			types.RFC3339Duration{Year: 11, Month: 12, Day: 13},
			"P11Y12M13D",
		),
		Entry(
			"year + month + day + second",
			"P11Y12M13DT16S",
			types.RFC3339Duration{Year: 11, Month: 12, Day: 13, Second: 16},
			"P11Y12M13DT16S",
		),
		Entry(
			"year + month + day + minute + second",
			"P11Y12M13DT15M16S",
			types.RFC3339Duration{Year: 11, Month: 12, Day: 13, Minute: 15, Second: 16},
			"P11Y12M13DT15M16S",
		),
		Entry(
			"year + month + day + hour + minute + second",
			"P11Y12M13DT14H15M16S",
			types.RFC3339Duration{Year: 11, Month: 12, Day: 13, Hour: 14, Minute: 15, Second: 16},
			"P11Y12M13DT14H15M16S",
		),
		Entry(
			"second",
			"PT16S",
			types.RFC3339Duration{Second: 16},
			"PT16S",
		),
		Entry(
			"minute + second",
			"PT15M16S",
			types.RFC3339Duration{Minute: 15, Second: 16},
			"PT15M16S",
		),
		Entry(
			"hour + minute + second",
			"PT14H15M16S",
			types.RFC3339Duration{Hour: 14, Minute: 15, Second: 16},
			"PT14H15M16S",
		),
		Entry(
			"week",
			"P17W",
			types.RFC3339Duration{Week: 17},
			"P17W",
		),
	)

	DescribeTable("Invalid values",
		func(input string) {
			_, err := format.ValidateValue(input)
			Expect(err).To(HaveOccurred())
		},
		Entry("empty", ""),
		Entry("random string", "foo"),
		Entry("P only", "P"),
		Entry("PT only", "PT"),
		Entry("valid prefix", "P1Sx"),
		Entry("valid suffix", "xP1S"),
	)

	When("a field type is types.RFC3339Duration", func() {
		It("should be parsed as string schema with duration RFC3339 format", func() {
			source := `
				import "github.com/conflowio/conflow/pkg/conflow/types"
				// @block "configuration"
				type Foo struct {
					v types.RFC3339Duration
				}
			`
			expectGoStructToHaveStringSchema(source, schema.FormatDurationRFC3339, false)
		})
	})

	When("a field type is *types.RFC3339Duration", func() {
		It("should be parsed as string schema with duration RFC3339 format", func() {
			source := `
				import "github.com/conflowio/conflow/pkg/conflow/types"
				// @block "configuration"
				type Foo struct {
					v *types.RFC3339Duration
				}
			`
			expectGoStructToHaveStringSchema(source, schema.FormatDurationRFC3339, true)
		})
	})

	It("should have a consistent JSON marshalling", func() {
		expectConsistentJSONMarshalling[*types.RFC3339Duration]([]byte("null"))
	})

})
