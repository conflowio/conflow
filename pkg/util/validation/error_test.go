// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package validation_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/pkg/util/validation"
)

var _ = Describe("Error", func() {
	It("should add a field error", func() {
		e := &validation.Error{}
		e.AddFieldError("foo", errors.New("some error"))
		Expect(e).To(MatchError("foo: some error"))
	})

	It("should add a formatted field error", func() {
		e := &validation.Error{}
		e.AddFieldErrorf("foo", "some error: %w", errors.New("other error"))
		Expect(e).To(MatchError("foo: some error: other error"))
	})

	It("should merge field errors", func() {
		e := &validation.Error{}
		e.AddFieldError("foo", validation.NewFieldError("bar", errors.New("some error")))
		Expect(e).To(MatchError("foo.bar: some error"))
	})

	It("should merge field with numeric field index", func() {
		e := &validation.Error{}
		e.AddFieldError("foo", validation.NewFieldError("[1]", errors.New("some error")))
		Expect(e).To(MatchError("foo[1]: some error"))
	})

	It("should merge field with string field index", func() {
		e := &validation.Error{}
		e.AddFieldError("foo", validation.NewFieldError(`["bar"]`, errors.New("some error")))
		Expect(e).To(MatchError(`foo["bar"]: some error`))
	})

	It("should flatten a validation error in a field error", func() {
		e := &validation.Error{}
		e.AddFieldError("foo", validation.NewError(
			validation.NewFieldError("bar", errors.New("some error")),
			validation.NewFieldError("baz", errors.New("other error")),
		))
		Expect(e).To(MatchError("foo.bar: some error, foo.baz: other error"))
	})

	It("should merge validation errors", func() {
		e := validation.NewError(
			validation.NewFieldError("foo", errors.New("some error")),
		)
		e.AddError(validation.NewError(
			validation.NewFieldError("bar", errors.New("other error")),
		))
		Expect(e).To(MatchError("foo: some error, bar: other error"))
	})
})
