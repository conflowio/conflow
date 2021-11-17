// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package utils_test

import (
	"reflect"
	"time"

	"github.com/conflowio/conflow/internal/testhelper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/internal/utils"
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
			"time": "other/time",
		}
		res := utils.GoType(imports, reflect.TypeOf(time.Time{}), false)
		Expect(res).To(Equal("time2.Time"))
		Expect(imports["time2"]).To(Equal("time"))
	})

	It("should handle a built in type pointer - another import with same name", func() {
		imports := map[string]string{
			"time": "other/time",
		}
		res := utils.GoType(imports, reflect.TypeOf(time.Time{}), true)
		Expect(res).To(Equal("*time2.Time"))
		Expect(imports["time2"]).To(Equal("time"))
	})

	It("should handle a user type", func() {
		imports := map[string]string{}
		res := utils.GoType(imports, reflect.TypeOf(testhelper.CustomStruct{}), false)
		Expect(res).To(Equal("testhelper.CustomStruct"))
		Expect(imports["testhelper"]).To(Equal("github.com/conflowio/conflow/internal/testhelper"))
	})

	It("should handle a user type pointer", func() {
		imports := map[string]string{}
		res := utils.GoType(imports, reflect.TypeOf(testhelper.CustomStruct{}), true)
		Expect(res).To(Equal("*testhelper.CustomStruct"))
		Expect(imports["testhelper"]).To(Equal("github.com/conflowio/conflow/internal/testhelper"))
	})

	It("should handle a user type - another import with same name", func() {
		imports := map[string]string{
			"testhelper": "other/testhelper",
		}
		res := utils.GoType(imports, reflect.TypeOf(testhelper.CustomStruct{}), false)
		Expect(res).To(Equal("testhelper2.CustomStruct"))
		Expect(imports["testhelper2"]).To(Equal("github.com/conflowio/conflow/internal/testhelper"))
	})

	It("should handle a user type pointer - another import with same name", func() {
		imports := map[string]string{
			"testhelper": "other/testhelper",
		}
		res := utils.GoType(imports, reflect.TypeOf(testhelper.CustomStruct{}), true)
		Expect(res).To(Equal("*testhelper2.CustomStruct"))
		Expect(imports["testhelper2"]).To(Equal("github.com/conflowio/conflow/internal/testhelper"))
	})
})
