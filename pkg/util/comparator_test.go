// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package util_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/pkg/util"
	"github.com/conflowio/conflow/pkg/util/ptr"
)

type testComparableObject struct {
	A string
}

func compareTestComparableObject(v1, v2 testComparableObject) bool {
	return util.ValueEquals(v1.A, v2.A)
}

func testValueEquals[T any](c util.Comparator[T], v1, v2 T, expected bool) func() {
	return func() {
		Expect(c(v1, v2)).To(Equal(expected))
		Expect(c(v2, v1)).To(Equal(expected))
	}
}

var _ = Describe("Equality", func() {

	DescribeTable("ValueEquals",
		func(f func()) {
			f()
		},
		Entry("equal strings", testValueEquals(util.ValueEquals[string], "foo", "foo", true)),
		Entry("different strings", testValueEquals(util.ValueEquals[string], "foo", "bar", false)),
		Entry("equal bools", testValueEquals(util.ValueEquals[bool], true, true, true)),
		Entry("different bools", testValueEquals(util.ValueEquals[bool], true, false, false)),
		Entry("equal ints", testValueEquals(util.ValueEquals[int64], int64(1), int64(1), true)),
		Entry("different ints", testValueEquals(util.ValueEquals[int64], int64(1), int64(2), false)),
		Entry("equal objects", testValueEquals(compareTestComparableObject, testComparableObject{A: "foo"}, testComparableObject{A: "foo"}, true)),
		Entry("different objects", testValueEquals(compareTestComparableObject, testComparableObject{A: "foo"}, testComparableObject{A: "bar"}, false)),
		Entry("equal floats", testValueEquals(util.FloatEquals[float64], 1.0, 1.0, true)),
		Entry("equal floats - tolerance 1", testValueEquals(util.FloatEquals[float64], 1.0, 1.0+util.Epsilon*0.9, true)),
		Entry("equal floats - tolerance 2", testValueEquals(util.FloatEquals[float64], 1.0+util.Epsilon*0.9, 1.0, true)),
		Entry("different floats", testValueEquals(util.FloatEquals[float64], 1.0, 2.0, false)),
		Entry("different floats - tolerance", testValueEquals(util.FloatEquals[float64], 1.0, 1.0+util.Epsilon, false)),
		Entry("different floats - tolerance", testValueEquals(util.FloatEquals[float64], 1.0+util.Epsilon, 1.0, false)),
		Entry("self comparator equals", testValueEquals(util.SelfComparator[time.Time](), time.Date(1, 2, 3, 4, 5, 6, 7, time.UTC), time.Date(1, 2, 3, 4, 5, 6, 7, time.UTC), true)),
		Entry("self comparator different", testValueEquals(util.SelfComparator[time.Time](), time.Date(1, 2, 3, 4, 5, 6, 7, time.UTC), time.Date(1, 2, 3, 4, 5, 6, 8, time.UTC), false)),
		Entry(
			"equal arrays",
			testValueEquals(
				util.ArrayEquals[int64](util.ValueEquals[int64]),
				[]int64{1},
				[]int64{1},
				true,
			),
		),
		Entry(
			"equal arrays - nil",
			testValueEquals(
				util.ArrayEquals[int64](util.ValueEquals[int64]),
				nil,
				nil,
				true,
			),
		),
		Entry(
			"different arrays - different values",
			testValueEquals(
				util.ArrayEquals[int64](util.ValueEquals[int64]),
				[]int64{1, 2},
				[]int64{1, 3},
				false,
			),
		),
		Entry(
			"different arrays - different length",
			testValueEquals(
				util.ArrayEquals[int64](util.ValueEquals[int64]),
				[]int64{1},
				[]int64{1, 2},
				false,
			),
		),
		Entry(
			"equal object arrays",
			testValueEquals(
				util.ArrayEquals[testComparableObject](compareTestComparableObject),
				[]testComparableObject{{A: "foo"}},
				[]testComparableObject{{A: "foo"}},
				true,
			),
		),
		Entry(
			"different object arrays",
			testValueEquals(
				util.ArrayEquals[testComparableObject](compareTestComparableObject),
				[]testComparableObject{{A: "foo"}},
				[]testComparableObject{{A: "bar"}},
				false,
			),
		),
		Entry(
			"equal maps",
			testValueEquals(
				util.MapEquals[int64](util.ValueEquals[int64]),
				map[string]int64{"foo": 1},
				map[string]int64{"foo": 1},
				true,
			),
		),
		Entry(
			"equal maps - nil",
			testValueEquals(
				util.MapEquals[int64](util.ValueEquals[int64]),
				nil,
				nil,
				true,
			),
		),
		Entry(
			"different maps - different values",
			testValueEquals(
				util.MapEquals[int64](util.ValueEquals[int64]),
				map[string]int64{"foo": 1},
				map[string]int64{"foo": 2},
				false,
			),
		),
		Entry(
			"different maps - different length",
			testValueEquals(
				util.MapEquals[int64](util.ValueEquals[int64]),
				map[string]int64{"foo": 1},
				map[string]int64{"foo": 1, "bar": 2},
				false,
			),
		),
		Entry(
			"different maps - different keys",
			testValueEquals(
				util.MapEquals[int64](util.ValueEquals[int64]),
				map[string]int64{"foo": 1, "bar": 2},
				map[string]int64{"foo": 1, "baz": 2},
				false,
			),
		),
		Entry(
			"equal object maps",
			testValueEquals(
				util.MapEquals[testComparableObject](compareTestComparableObject),
				map[string]testComparableObject{"foo": {A: "a1"}},
				map[string]testComparableObject{"foo": {A: "a1"}},
				true,
			),
		),
		Entry(
			"different object maps",
			testValueEquals(
				util.MapEquals[testComparableObject](compareTestComparableObject),
				map[string]testComparableObject{"foo": {A: "a1"}},
				map[string]testComparableObject{"foo": {A: "a2"}},
				false,
			),
		),
		Entry(
			"equal pointers",
			testValueEquals(
				util.PointerEquals[int64](util.ValueEquals[int64]),
				ptr.To(int64(1)),
				ptr.To(int64(1)),
				true,
			),
		),
		Entry(
			"equal pointers - nil",
			testValueEquals(
				util.PointerEquals[int64](util.ValueEquals[int64]),
				nil,
				nil,
				true,
			),
		),
		Entry(
			"different pointers - different values",
			testValueEquals(
				util.PointerEquals[int64](util.ValueEquals[int64]),
				ptr.To(int64(1)),
				ptr.To(int64(2)),
				false,
			),
		),
		Entry(
			"different pointers - one nil",
			testValueEquals(
				util.PointerEquals[int64](util.ValueEquals[int64]),
				nil,
				ptr.To(int64(2)),
				false,
			),
		),
		Entry(
			"different pointers - other nil",
			testValueEquals(
				util.PointerEquals[int64](util.ValueEquals[int64]),
				ptr.To(int64(1)),
				nil,
				false,
			),
		),
		Entry(
			"equal complex type",
			testValueEquals(
				util.MapEquals[[]int64](util.ArrayEquals[int64](util.ValueEquals[int64])),
				map[string][]int64{"foo": {1}},
				map[string][]int64{"foo": {1}},
				true,
			),
		),
		Entry(
			"different complex type",
			testValueEquals(
				util.MapEquals[[]int64](util.ArrayEquals[int64](util.ValueEquals[int64])),
				map[string][]int64{"foo": {1}},
				map[string][]int64{"foo": {2}},
				false,
			),
		),
	)
})
