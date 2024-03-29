// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package formats_test

import (
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/pkg/schema"
	"github.com/conflowio/conflow/pkg/schema/formats"
)

var _ = Describe("UUID", func() {

	format := formats.UUID{}

	uuid1 := uuid.MustParse("c94b8114-3f28-11ec-9bbc-0242ac130002")
	uuid4 := uuid.MustParse("cb95ec16-54d2-46f0-a63f-e94c0803d0d1")

	DescribeTable("Valid values",
		expectFormatToParse[uuid.UUID](format),
		Entry("version 1", "c94b8114-3f28-11ec-9bbc-0242ac130002", uuid1, "c94b8114-3f28-11ec-9bbc-0242ac130002"),
		Entry("version 4", "cb95ec16-54d2-46f0-a63f-e94c0803d0d1", uuid4, "cb95ec16-54d2-46f0-a63f-e94c0803d0d1"),
	)

	DescribeTable("Invalid values",
		func(input string) {
			_, err := format.ValidateValue(input)
			Expect(err).To(HaveOccurred())
		},
		Entry("empty", ""),
		Entry("random string", "foo"),
		Entry("incomplete", "c94b8114-3f28-11ec-9bbc-0242ac13000"),
	)

	When("a field type is uuid.UUID", func() {
		It("should be parsed as string schema with uuid format", func() {
			source := `
				import "github.com/google/uuid"
				// @block "configuration"
				type Foo struct {
					v uuid.UUID
				}
			`
			expectGoStructToHaveStringSchema(source, schema.FormatUUID, false)
		})
	})

	When("a field type is *uuid.UUID", func() {
		It("should be parsed as string schema with uuid format", func() {
			source := `
				import "github.com/google/uuid"
				// @block "configuration"
				type Foo struct {
					v *uuid.UUID
				}
			`
			expectGoStructToHaveStringSchema(source, schema.FormatUUID, true)
		})
	})

	It("should have a consistent JSON marshalling", func() {
		expectConsistentJSONMarshalling[*uuid.UUID]([]byte("null"))
	})

})
