// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package utils_test

import (
	"reflect"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/src/internal/testhelper"
	"github.com/conflowio/conflow/src/internal/utils"
)

var _ = DescribeTable(
	"EnsureUniqueGoPackageSelector",
	func(imports map[string]string, pkg string, expectedImports map[string]string, expectedSel string) {
		sel := utils.EnsureUniqueGoPackageSelector(imports, pkg)
		Expect(sel).To(Equal(expectedSel))
		Expect(imports).To(Equal(expectedImports))
	},
	Entry(
		"adds a package with no alias",
		map[string]string{},
		"github.com/foo/bar",
		map[string]string{"github.com/foo/bar": "bar"},
		"bar.",
	),
	Entry(
		"adds a standard lib",
		map[string]string{},
		"fmt",
		map[string]string{"fmt": "fmt"},
		"fmt.",
	),
	Entry(
		"adds an alias if the selector is already taken",
		map[string]string{"github.com/foo/bar": "bar"},
		"github.com/baz/bar",
		map[string]string{"github.com/foo/bar": "bar", "github.com/baz/bar": "bar2"},
		"bar2.",
	),
	Entry(
		"returns empty string for the current package",
		map[string]string{"github.com/foo/bar": ""},
		"github.com/foo/bar",
		map[string]string{"github.com/foo/bar": ""},
		"",
	),
	Entry(
		"returns the real package name for version packages",
		map[string]string{},
		"github.com/foo/bar/v1",
		map[string]string{"github.com/foo/bar/v1": "bar"},
		"bar.",
	),
)

var _ = Describe("GoType", func() {
	It("should return with a simple type", func() {
		imports := map[string]string{}
		res := utils.GoType(imports, reflect.TypeOf(""), false)
		Expect(res).To(Equal("string"))
		Expect(imports).To(BeEmpty())
	})

	It("should return with a simple type pointer", func() {
		imports := map[string]string{}
		res := utils.GoType(imports, reflect.TypeOf(""), true)
		Expect(res).To(Equal("*string"))
		Expect(imports).To(BeEmpty())
	})

	It("should handle a built in type", func() {
		imports := map[string]string{}
		res := utils.GoType(imports, reflect.TypeOf(time.Time{}), false)
		Expect(res).To(Equal("time.Time"))
		Expect(imports["time"]).To(Equal("time"))
	})

	It("should handle a built in type pointer", func() {
		imports := map[string]string{}
		res := utils.GoType(imports, reflect.TypeOf(time.Time{}), true)
		Expect(res).To(Equal("*time.Time"))
		Expect(imports["time"]).To(Equal("time"))
	})

	It("should handle a built in type - another import with same name", func() {
		imports := map[string]string{
			"other/time": "time",
		}
		res := utils.GoType(imports, reflect.TypeOf(time.Time{}), false)
		Expect(res).To(Equal("time2.Time"))
		Expect(imports["time"]).To(Equal("time2"))
	})

	It("should handle a built in type pointer - another import with same name", func() {
		imports := map[string]string{
			"other/time": "time",
		}
		res := utils.GoType(imports, reflect.TypeOf(time.Time{}), true)
		Expect(res).To(Equal("*time2.Time"))
		Expect(imports["time"]).To(Equal("time2"))
	})

	It("should handle a user type", func() {
		imports := map[string]string{}
		res := utils.GoType(imports, reflect.TypeOf(testhelper.CustomStruct{}), false)
		Expect(res).To(Equal("testhelper.CustomStruct"))
		Expect(imports["github.com/conflowio/conflow/src/internal/testhelper"]).To(Equal("testhelper"))
	})

	It("should handle a user type pointer", func() {
		imports := map[string]string{}
		res := utils.GoType(imports, reflect.TypeOf(testhelper.CustomStruct{}), true)
		Expect(res).To(Equal("*testhelper.CustomStruct"))
		Expect(imports["github.com/conflowio/conflow/src/internal/testhelper"]).To(Equal("testhelper"))
	})

	It("should handle a user type - another import with same name", func() {
		imports := map[string]string{
			"other/testhelper": "testhelper",
		}
		res := utils.GoType(imports, reflect.TypeOf(testhelper.CustomStruct{}), false)
		Expect(res).To(Equal("testhelper2.CustomStruct"))
		Expect(imports["github.com/conflowio/conflow/src/internal/testhelper"]).To(Equal("testhelper2"))
	})

	It("should handle a user type pointer - another import with same name", func() {
		imports := map[string]string{
			"other/testhelper": "testhelper",
		}
		res := utils.GoType(imports, reflect.TypeOf(testhelper.CustomStruct{}), true)
		Expect(res).To(Equal("*testhelper2.CustomStruct"))
		Expect(imports["github.com/conflowio/conflow/src/internal/testhelper"]).To(Equal("testhelper2"))
	})
})
