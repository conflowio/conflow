// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package util_test

import (
	"reflect"
	"time"

	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/pkg/util"
	"github.com/conflowio/conflow/pkg/util/ptr"
)

type testClonableObject struct {
	A string
}

func cloneTestClonableObject(v1 testClonableObject) testClonableObject {
	return testClonableObject{
		A: util.CloneValue(v1.A),
	}
}

type testClonerObject struct {
	A string
}

func (t testClonerObject) Clone() testClonerObject {
	return testClonerObject{A: t.A}
}

func testClone[T any](c util.Cloner[T], v1 T) func() {
	return func() {
		v2 := c(v1)
		Expect(v1).To(Equal(v2))
		expectDeepCopy(reflect.ValueOf(v1), reflect.ValueOf(v2))
	}
}

func expectDeepCopy(v1, v2 reflect.Value) {
	k := v1.Kind()
	if (k == reflect.Pointer || k == reflect.Map || k == reflect.Slice) && !v1.IsNil() {
		Expect(v1.Pointer()).ToNot(Equal(v2.Pointer()))
	}
	switch k {
	case reflect.Map:
		iter := v1.MapRange()
		for iter.Next() {
			k := iter.Key()
			expectDeepCopy(iter.Value(), v2.MapIndex(k))
		}
	case reflect.Slice:
		for i := 0; i < v1.Len(); i++ {
			expectDeepCopy(v1.Index(i), v2.Index(i))
		}
	}
}

var _ = DescribeTable("Clone",
	func(f func()) {
		f()
	},
	Entry("string", testClone(util.CloneValue[string], "foo")),
	Entry("bool", testClone(util.CloneValue[bool], true)),
	Entry("int", testClone(util.CloneValue[int64], int64(1))),
	Entry("float", testClone(util.CloneValue[float64], 1.0)),
	Entry("time duration", testClone(util.CloneValue[time.Duration], time.Hour+30*time.Minute)),
	Entry("time", testClone(util.CloneValue[time.Time], time.Now())),
	Entry("object", testClone(cloneTestClonableObject, testClonableObject{A: "foo"})),
	Entry("self cloner", testClone(util.SelfCloner[testClonerObject](), testClonerObject{A: "foo"})),
	Entry(
		"array",
		testClone(
			util.CloneArray[int64](util.CloneValue[int64]),
			[]int64{1},
		),
	),
	Entry(
		"nil array",
		testClone(
			util.CloneArray[int64](util.CloneValue[int64]),
			nil,
		),
	),
	Entry(
		"object array",
		testClone(
			util.CloneArray[testClonableObject](cloneTestClonableObject),
			[]testClonableObject{{A: "foo"}},
		),
	),
	Entry(
		"map",
		testClone(
			util.CloneMap[int64](util.CloneValue[int64]),
			map[string]int64{"foo": 1},
		),
	),
	Entry(
		"map - nil",
		testClone(
			util.CloneMap[int64](util.CloneValue[int64]),
			nil,
		),
	),
	Entry(
		"object map",
		testClone(
			util.CloneMap[testClonableObject](cloneTestClonableObject),
			map[string]testClonableObject{"foo": {A: "a1"}},
		),
	),
	Entry(
		"pointer",
		testClone(
			util.ClonePointer[int64](util.CloneValue[int64]),
			ptr.To(int64(1)),
		),
	),
	Entry(
		"pointers - nil",
		testClone(
			util.ClonePointer[int64](util.CloneValue[int64]),
			nil,
		),
	),
	Entry(
		"equal complex type",
		testClone(
			util.CloneMap[[]int64](util.CloneArray[int64](util.CloneValue[int64])),
			map[string][]int64{"foo": {1}},
		),
	),
)
