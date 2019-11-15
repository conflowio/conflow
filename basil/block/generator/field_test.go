// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/basil/basil/block/generator"
	"github.com/opsidian/basil/basil/variable"
)

var _ = Describe("Field", func() {
	var f *generator.Field

	BeforeEach(func() {
		f = &generator.Field{
			Name:        "foo",
			ParamName:   "param_foo",
			Type:        "string",
			IsRequired:  false,
			Stage:       "main",
			IsID:        false,
			IsValue:     false,
			IsReference: false,
			IsBlock:     false,
		}
	})

	Describe("Validate", func() {

		It("returns no error for a valid field", func() {
			Expect(f.Validate()).To(BeNil())
		})

		It("allows the reference tag on an id field", func() {
			f.IsID = true
			f.IsReference = true
			f.Type = variable.TypeIdentifier
			Expect(f.Validate()).To(BeNil())
		})

		It("allows an array type with block", func() {
			f.IsBlock = true
			f.Type = "[]SomeBlock"
			Expect(f.Validate()).To(BeNil())
		})

		It("returns error if id and value are both set", func() {
			f.IsID = true
			f.IsValue = true
			Expect(f.Validate()).To(MatchError("field \"foo\" must have exactly one of: id, value, block or generated"))
		})

		It("returns error if reference is on a non-id field", func() {
			f.IsReference = true
			Expect(f.Validate()).To(MatchError("the \"reference\" tag can only be set on the id field"))
		})

		It("returns error if param name is invalid", func() {
			f.ParamName = "invalid name"
			Expect(f.Validate()).To(MatchError("\"name\" tag is invalid on field \"foo\", it must be a valid identifier"))
		})

		It("returns an error for an invalid type", func() {
			f.Type = "invalidtype"
			Expect(f.Validate()).To(MatchError("invalid field type \"invalidtype\" on field \"foo\", use a valid type or use ignore tag"))
		})

		It("returns an error for an empty stage", func() {
			f.Stage = ""
			Expect(f.Validate()).To(MatchError("\"stage\" can not be empty on field \"foo\""))
		})

		It("returns an error for a non-string id field", func() {
			f.IsID = true
			f.Type = "int64"
			Expect(f.Validate()).To(MatchError("field \"foo\" must be defined as basil.ID"))
		})

	})
})
