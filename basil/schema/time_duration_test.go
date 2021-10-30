// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package schema_test

import (
	"encoding/json"
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/opsidian/conflow/basil/schema"
)

var _ schema.Schema = &schema.TimeDuration{}

var _ = Describe("TimeDuration", func() {
	DescribeTable("Validate accepts value",
		func(schema *schema.TimeDuration, value interface{}) {
			err := schema.ValidateValue(value)
			Expect(err).ToNot(HaveOccurred())
		},
		Entry("zero value", &schema.TimeDuration{}, time.Duration(0)),
		Entry("non-zero value", &schema.TimeDuration{}, time.Second),
		Entry("const value", &schema.TimeDuration{Const: schema.TimeDurationPtr(time.Second)}, time.Second),
		Entry("enum value - single", &schema.TimeDuration{Enum: []time.Duration{time.Second}}, time.Second),
		Entry("enum value - multiple", &schema.TimeDuration{Enum: []time.Duration{time.Second, 2 * time.Second}}, time.Second),
	)

	DescribeTable("Validate errors",
		func(schema *schema.TimeDuration, value interface{}, expectedErr error) {
			err := schema.ValidateValue(value)
			Expect(err).To(MatchError(expectedErr))
		},
		Entry(
			"not a time duration value",
			&schema.TimeDuration{},
			"foo",
			errors.New("must be time duration"),
		),
		Entry(
			"const value",
			&schema.TimeDuration{Const: schema.TimeDurationPtr(time.Second)},
			2*time.Second,
			errors.New(`must be 1s`),
		),
		Entry(
			"enum value - single",
			&schema.TimeDuration{Enum: []time.Duration{time.Second}},
			2*time.Second,
			errors.New(`must be 1s`),
		),
		Entry(
			"enum value - multiple",
			&schema.TimeDuration{Enum: []time.Duration{time.Second, 2 * time.Second}},
			3*time.Second,
			errors.New(`must be one of 1s, 2s`),
		),
	)

	DescribeTable("GoString prints a valid Go struct",
		func(schema *schema.TimeDuration, expected string) {
			str := schema.GoString()
			Expect(str).To(Equal(expected))
		},
		Entry(
			"empty",
			&schema.TimeDuration{},
			`&schema.TimeDuration{
}`,
		),
		Entry(
			"const",
			&schema.TimeDuration{Const: schema.TimeDurationPtr(time.Second)},
			`&schema.TimeDuration{
	Const: schema.TimeDurationPtr(1000000000),
}`,
		),
		Entry(
			"default",
			&schema.TimeDuration{Default: schema.TimeDurationPtr(time.Second)},
			`&schema.TimeDuration{
	Default: schema.TimeDurationPtr(1000000000),
}`,
		),
		Entry(
			"enum",
			&schema.TimeDuration{Enum: []time.Duration{time.Second}},
			`&schema.TimeDuration{
	Enum: []time.Duration{1000000000},
}`,
		),
	)

	It("should marshal/unmarshal", func() {
		s := &schema.TimeDuration{
			Const:   schema.TimeDurationPtr(time.Second),
			Default: schema.TimeDurationPtr(2 * time.Second),
			Enum:    []time.Duration{3 * time.Second, 4 * time.Second},
		}
		j, err := json.Marshal(s)
		Expect(err).ToNot(HaveOccurred())

		s2 := &schema.TimeDuration{}
		err = json.Unmarshal(j, &s2)
		Expect(err).ToNot(HaveOccurred())
		Expect(s2).To(Equal(s))
	})
})
