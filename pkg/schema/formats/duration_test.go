// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package formats_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/pkg/conflow/types"
	"github.com/conflowio/conflow/pkg/schema"
	"github.com/conflowio/conflow/pkg/schema/formats"
)

var _ = Describe("Duration", func() {

	format := formats.Duration{}

	DescribeTable("Valid values",
		expectFormatToParse[types.Duration](format),

		Entry(
			"some duration",
			"1h2m3s",
			types.Duration(time.Hour+2*time.Minute+3*time.Second),
			"1h2m3s",
		),
	)

	DescribeTable("Invalid values",
		func(input string) {
			_, err := format.ValidateValue(input)
			Expect(err).To(HaveOccurred())
		},
		Entry("empty", ""),
		Entry("random string", "foo"),
	)

	When("a field type is types.Duration", func() {
		It("should be parsed as string schema with duration format", func() {
			source := `
				import "github.com/conflowio/conflow/pkg/conflow/types"
				// @block "configuration"
				type Foo struct {
					v types.Duration
				}
			`
			expectGoStructToHaveStringSchema(source, schema.FormatDuration, false)
		})
	})

	When("a field type is *types.Duration", func() {
		It("should be parsed as string schema with duration format", func() {
			source := `
				import "github.com/conflowio/conflow/pkg/conflow/types"
				// @block "configuration"
				type Foo struct {
					v *types.Duration
				}
			`
			expectGoStructToHaveStringSchema(source, schema.FormatDuration, true)
		})
	})

	It("should have a consistent JSON marshalling", func() {
		expectConsistentJSONMarshalling[*types.Duration]([]byte("null"))
	})

})
