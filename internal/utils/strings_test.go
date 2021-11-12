// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package utils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/internal/utils"
)

var _ = Describe("ToSnakeCase", func() {
	DescribeTable(
		"It should convert string to snake_case",
		func(input string, expected string) {
			Expect(utils.ToSnakeCase(input)).To(Equal(expected))
		},
		Entry("lowercase", "foo", "foo"),
		Entry("starts with capital", "Foo", "foo"),
		Entry("capital in the middle", "fooBar", "foo_bar"),
		Entry("multiple capitals", "FooBarBaz", "foo_bar_baz"),
		Entry("capital at the end", "fooB", "foo_b"),
		Entry("already snake_case", "foo_bar", "foo_bar"),
		Entry("only capitals", "FOO", "foo"),
		Entry("multiple capitals at the beginning", "FOOBar", "foo_bar"),
		Entry("multiple capitals in the middle", "fooBARBaz", "foo_bar_baz"),
		Entry("multiple capitals at the end", "fooBarBAZ", "foo_bar_baz"),
		Entry("with numbers", "Foo9", "foo9"),
		Entry("numbers after single capital letter", "I18nDict", "i18n_dict"),
		Entry("single letter, single number", "F9", "f9"),
	)
})

var _ = Describe("ToCamelCase", func() {
	DescribeTable(
		"It should convert string to camelCase",
		func(input string, expected string) {
			Expect(utils.ToCamelCase(input)).To(Equal(expected))
		},
		Entry("lowercase", "foo", "Foo"),
		Entry("starts with capital", "Foo", "Foo"),
		Entry("two parts with underscore", "foo_bar", "FooBar"),
		Entry("two parts with multiple underscores", "foo__bar", "FooBar"),
		Entry("three parts capitals", "foo_bar_baz", "FooBarBaz"),
		Entry("single character in the second part", "foo_b", "FooB"),
		Entry("single character in the first part", "f_bar", "FBar"),
		Entry("already camelCase", "FooBar", "FooBar"),
		Entry("only capitals", "FOO", "FOO"),
		Entry("with numbers", "foo_9", "Foo9"),
		Entry("numbers after first letter", "i18n_dict", "I18nDict"),
		Entry("single letter, single number", "F9", "F9"),
		Entry("abbreviation should be left as is", "YAMLDocument", "YAMLDocument"),
	)
})
