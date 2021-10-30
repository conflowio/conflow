// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/conflowio/conflow/conflow/schema"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ schema.Schema = &schema.Time{}

var _ = Describe("Time", func() {

	time1, _ := time.Parse(time.RFC3339, "2001-01-01T00:00:01Z")
	time2, _ := time.Parse(time.RFC3339, "2001-01-01T00:00:02Z")
	time3, _ := time.Parse(time.RFC3339, "2001-01-01T00:00:03Z")

	DescribeTable("Validate accepts value",
		func(schema *schema.Time, value interface{}) {
			err := schema.ValidateValue(value)
			Expect(err).ToNot(HaveOccurred())
		},
		Entry("zero value", &schema.Time{}, time.Time{}),
		Entry("non-zero value", &schema.Time{}, time1),
		Entry("const value", &schema.Time{Const: schema.TimePtr(time1)}, time1),
		Entry("enum value - single", &schema.Time{Enum: []time.Time{time1}}, time1),
		Entry("enum value - multiple", &schema.Time{Enum: []time.Time{time1, time2}}, time1),
	)

	DescribeTable("Validate errors",
		func(schema *schema.Time, value interface{}, expectedErr error) {
			err := schema.ValidateValue(value)
			Expect(err).To(MatchError(expectedErr))
		},
		Entry(
			"not a time duration value",
			&schema.Time{},
			"foo",
			errors.New("must be date-time"),
		),
		Entry(
			"const value",
			&schema.Time{Const: schema.TimePtr(time1)},
			time2,
			fmt.Errorf("must be %s", time1.String()),
		),
		Entry(
			"enum value - single",
			&schema.Time{Enum: []time.Time{time1}},
			time2,
			fmt.Errorf("must be %s", time1.String()),
		),
		Entry(
			"enum value - multiple",
			&schema.Time{Enum: []time.Time{time1, time2}},
			time3,
			fmt.Errorf("must be one of %s, %s", time1.String(), time2.String()),
		),
	)

	DescribeTable("GoString prints a valid Go struct",
		func(schema *schema.Time, expected string) {
			str := schema.GoString()
			Expect(str).To(Equal(expected))
		},
		Entry(
			"empty",
			&schema.Time{},
			`&schema.Time{
}`,
		),
		Entry(
			"const",
			&schema.Time{Const: schema.TimePtr(time1)},
			`&schema.Time{
	Const: schema.TimePtr(time.Date(2001, 1, 1, 0, 0, 1, 0, time.UTC)),
}`,
		),
		Entry(
			"default",
			&schema.Time{Default: schema.TimePtr(time1)},
			`&schema.Time{
	Default: schema.TimePtr(time.Date(2001, 1, 1, 0, 0, 1, 0, time.UTC)),
}`,
		),
		Entry(
			"enum",
			&schema.Time{Enum: []time.Time{time1}},
			`&schema.Time{
	Enum: []time.Time{time.Date(2001, 1, 1, 0, 0, 1, 0, time.UTC)},
}`,
		),
	)

	It("should marshal/unmarshal", func() {
		s := &schema.Time{
			Const:   schema.TimePtr(time1),
			Default: schema.TimePtr(time2),
			Enum:    []time.Time{time3},
		}
		j, err := json.Marshal(s)
		Expect(err).ToNot(HaveOccurred())

		s2 := &schema.Time{}
		err = json.Unmarshal(j, &s2)
		Expect(err).ToNot(HaveOccurred())
		Expect(s2).To(Equal(s))
	})
})
