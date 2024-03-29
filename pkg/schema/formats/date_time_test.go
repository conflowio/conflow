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

	"github.com/conflowio/conflow/pkg/schema"
	"github.com/conflowio/conflow/pkg/schema/formats"
)

var _ = Describe("DateTime", func() {

	format := formats.DateTime{}

	DescribeTable("Valid values",
		expectFormatToParse[time.Time](format, func(t1 time.Time, t2 time.Time) bool {
			return t1.Equal(t2)
		}),
		Entry(
			"date and time with zero timezone",
			"2021-01-02T12:23:34Z",
			time.Date(2021, 01, 02, 12, 23, 34, 0, time.UTC),
			"2021-01-02T12:23:34Z",
		),
		Entry(
			"date and time with positive timezone",
			"2021-01-02T12:23:34+01:00",
			time.Date(2021, 01, 02, 12, 23, 34, 0, time.FixedZone("", 3600)),
			"2021-01-02T12:23:34+01:00",
		),
		Entry(
			"date and time with negative timezone",
			"2021-01-02T12:23:34-01:00",
			time.Date(2021, 01, 02, 12, 23, 34, 0, time.FixedZone("", -3600)),
			"2021-01-02T12:23:34-01:00",
		),
		Entry(
			"date and time with fractional seconds and empty timezone",
			"2021-01-02T12:23:34.123Z",
			time.Date(2021, 01, 02, 12, 23, 34, 123000000, time.UTC),
			"2021-01-02T12:23:34.123Z",
		),
	)

	DescribeTable("Invalid values",
		func(input string) {
			_, err := format.ValidateValue(input)
			Expect(err).To(HaveOccurred())
		},
		Entry("empty", ""),
		Entry("random string", "foo"),
		Entry("just date", "2021-01-02"),
		Entry("date and time with no timezone", "2021-01-02T12:23:34"),
		Entry("date and time with fractional seconds and no timezone", "2021-01-02T12:23:34.123"),
		Entry("incomplete - no seconds", "2021-01-02T12:23"),
		Entry("timezone - Z with value", "2021-01-02T12:23:34Z+01:00"),
	)

	When("a field type is time.Time", func() {
		It("should be parsed as string schema with datetime format", func() {
			source := `
				// @block "configuration"
				type Foo struct {
					v time.Time
				}
			`
			expectGoStructToHaveStringSchema(source, schema.FormatDateTime, false)
		})
	})

	When("a field type is *time.Time", func() {
		It("should be parsed as string schema with datetime format", func() {
			source := `
				// @block "configuration"
				type Foo struct {
					v *time.Time
				}
			`
			expectGoStructToHaveStringSchema(source, schema.FormatDateTime, true)
		})
	})

	It("should have a consistent JSON marshalling", func() {
		expectConsistentJSONMarshalling[*time.Time]([]byte("null"))
	})

})
