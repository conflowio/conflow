// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package formats_test

import (
	"net/mail"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/pkg/conflow/types"
	"github.com/conflowio/conflow/pkg/schema"
	"github.com/conflowio/conflow/pkg/schema/formats"
)

var _ = Describe("Email", func() {

	format := formats.Email{}

	DescribeTable("Valid values",
		expectFormatToParse[types.Email](format),
		Entry(
			"short",
			"a@b.c",
			types.Email(mail.Address{Address: "a@b.c"}),
			"a@b.c",
		),
		Entry(
			"simple hostname",
			"a@localhost",
			types.Email(mail.Address{Address: "a@localhost"}),
			"a@localhost",
		),
		Entry(
			"complete",
			"my.name@email.company.com",
			types.Email(mail.Address{Address: "my.name@email.company.com"}),
			"my.name@email.company.com",
		),
		Entry(
			"+ in first part",
			"my.name+foo@email.company.com",
			types.Email(mail.Address{Address: "my.name+foo@email.company.com"}),
			"my.name+foo@email.company.com",
		),
		Entry(
			"with name",
			`"My Name" <my.name@email.company.com>`,
			types.Email(mail.Address{Name: "My Name", Address: "my.name@email.company.com"}),
			`"My Name" <my.name@email.company.com>`,
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

	When("a field type is mail.Address", func() {
		It("should be parsed as string schema with email format", func() {
			source := `
				import "github.com/conflowio/conflow/pkg/conflow/types"
				// @block "configuration"
				type Foo struct {
					v types.Email
				}
			`
			expectGoStructToHaveStringSchema(source, schema.FormatEmail, false)
		})
	})

	When("a field type is *mail.Address", func() {
		It("should be parsed as string schema with email format", func() {
			source := `
				import "github.com/conflowio/conflow/pkg/conflow/types"
				// @block "configuration"
				type Foo struct {
					v *types.Email
				}
			`
			expectGoStructToHaveStringSchema(source, schema.FormatEmail, true)
		})
	})

	It("should have a consistent JSON marshalling", func() {
		expectConsistentJSONMarshalling[*types.Email]([]byte("null"))
	})

})
