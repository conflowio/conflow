// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package values_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/pkg/values"
)

var _ = Describe("List", func() {
	It("copies input slice so mutating original does not affect list", func() {
		original := []string{"a", "b"}
		list := values.NewListFromSlice(original)
		original[0] = "mutated"

		Expect(list.At(0)).To(Equal("a"))
		Expect(list.At(1)).To(Equal("b"))
	})

	It("returns a copy from Elems that does not affect the list when mutated", func() {
		list := values.NewList("a", "b")
		elems := list.Elems()
		elems[0] = "mutated"

		Expect(list.At(0)).To(Equal("a"))
	})

	It("creates independent lists from NewList variadic args", func() {
		elems := []string{"a", "b"}
		list := values.NewList(elems...)
		elems[0] = "mutated"

		Expect(list.At(0)).To(Equal("a"))
	})
})

var _ = Describe("ListBuilder", func() {
	It("returns an independent immutable list from Freeze", func() {
		builder := values.NewListBuilder[string]()
		builder.Append("a")
		frozen := builder.Freeze()

		builder.Append("b")

		Expect(frozen.Len()).To(Equal(1))
		Expect(frozen.At(0)).To(Equal("a"))
	})

	It("does not share backing storage when builder slice grows after Freeze", func() {
		builder := values.NewListBuilder[string]()
		builder.Append("shared")
		frozen := builder.Freeze()

		for i := 0; i < 10; i++ {
			builder.Append("extra")
		}

		Expect(frozen.Len()).To(Equal(1))
		Expect(frozen.At(0)).To(Equal("shared"))
	})

	It("allows sharing the frozen pointer without aliasing builder state", func() {
		builder := values.NewListBuilder[string]()
		builder.Append("a")
		frozen1 := builder.Freeze()
		frozen2 := frozen1

		builder.Append("b")

		Expect(frozen1).To(BeIdenticalTo(frozen2))
		Expect(frozen1.Len()).To(Equal(1))
		Expect(frozen1.At(0)).To(Equal("a"))
	})
})
