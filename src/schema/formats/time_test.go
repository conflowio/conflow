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

	"github.com/conflowio/conflow/src/conflow/types"
	"github.com/conflowio/conflow/src/schema"
	"github.com/conflowio/conflow/src/schema/formats"
)

var _ = Describe("Time", func() {

	format := formats.Time{}

	DescribeTable("Valid values",
		expectFormatToParse(format),
		Entry(
			"time with no timezone",
			"12:23:34",
			types.Time{Hour: 12, Minute: 23, Second: 34, Location: time.UTC},
			"12:23:34Z",
		),
		Entry(
			"time with zero timezone",
			"12:23:34Z",
			types.Time{Hour: 12, Minute: 23, Second: 34, Location: time.UTC},
			"12:23:34Z",
		),
		Entry(
			"time with positive timezone",
			"12:23:34+01:00",
			types.Time{Hour: 12, Minute: 23, Second: 34, Location: time.FixedZone("", 3600)},
			"12:23:34+01:00",
		),
		Entry(
			"time with negative timezone",
			"12:23:34-01:00",
			types.Time{Hour: 12, Minute: 23, Second: 34, Location: time.FixedZone("", -3600)},
			"12:23:34-01:00",
		),
		Entry(
			"time with fractional seconds and no timezone",
			"12:23:34.123",
			types.Time{Hour: 12, Minute: 23, Second: 34, NanoSecond: 123000000, Location: time.UTC},
			"12:23:34.123Z",
		),
		Entry(
			"time with fractional seconds and empty timezone",
			"12:23:34.123Z",
			types.Time{Hour: 12, Minute: 23, Second: 34, NanoSecond: 123000000, Location: time.UTC},
			"12:23:34.123Z",
		),
	)

	DescribeTable("Invalid values",
		func(input string) {
			_, err := format.ValidateValue(input)
			Expect(err).To(HaveOccurred())
		},
		Entry("empty", ""),
		Entry("random string", "foo"),
		Entry("incomplete - no seconds", "12:23"),
		Entry("timezone - Z with value", "12:23:34Z+01:00"),
	)

	When("a field type is types.Time", func() {
		It("should be parsed as string schema with time format", func() {
			source := `
				import "github.com/conflowio/conflow/src/conflow/types"
				// @block "configuration"
				type Foo struct {
					v types.Time
				}
			`
			expectGoStructToHaveStringSchema(source, schema.FormatTime, false)
		})
	})

	When("a field type is *types.Time", func() {
		It("should be parsed as string schema with time format", func() {
			source := `
				import "github.com/conflowio/conflow/src/conflow/types"
				// @block "configuration"
				type Foo struct {
					v *types.Time
				}
			`
			expectGoStructToHaveStringSchema(source, schema.FormatTime, true)
		})
	})

})
