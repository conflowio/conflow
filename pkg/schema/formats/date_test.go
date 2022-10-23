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

var _ = Describe("Date", func() {

	format := formats.Date{}

	DescribeTable("Valid values",
		expectFormatToParse[types.Date](format),
		Entry(
			"any date",
			"2021-01-02",
			types.NewDate(2021, 1, 2),
			"2021-01-02",
		),
	)

	DescribeTable("Invalid values",
		func(input string) {
			_, err := format.ValidateValue(input)
			Expect(err).To(HaveOccurred())
		},
		Entry("empty", ""),
		Entry("random string", "foo"),
		Entry("incomplete - no day", "2021-01"),
		Entry("incomplete - short day", "2021-01-1"),
		Entry("non-existing day", "2021-02-31"),
	)

	When("a field type is types.Date and format is set as date", func() {
		It("should be parsed as string schema with date format", func() {
			source := `
				import "github.com/conflowio/conflow/pkg/conflow/types" 
				// @block "configuration"
				type Foo struct {
					v types.Date
				}
			`
			expectGoStructToHaveStringSchema(source, schema.FormatDate, false)
		})
	})

	When("a field type is *types.Date and format is set as date", func() {
		It("should be parsed as string schema with date format", func() {
			source := `
				import "github.com/conflowio/conflow/pkg/conflow/types"
				// @block "configuration"
				type Foo struct {
					v *types.Date
				}
			`
			expectGoStructToHaveStringSchema(source, schema.FormatDate, true)
		})
	})

	It("should have a consistent JSON marshalling", func() {
		expectConsistentJSONMarshalling[*types.Date]([]byte("null"))
	})

})
