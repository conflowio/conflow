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

var _ = Describe("Convert helpers", func() {
	It("converts slices to immutable lists", func() {
		original := []string{"a", "b"}
		list := values.FromSlice(original)
		original[0] = "mutated"

		Expect(list.At(0)).To(Equal("a"))
	})

	It("converts []interface{} to immutable lists", func() {
		original := []interface{}{"a", int64(1)}
		list, err := values.FromInterfaceSlice(original)
		Expect(err).NotTo(HaveOccurred())
		original[0] = "mutated"

		Expect(list.At(0)).To(Equal("a"))
		Expect(list.At(1)).To(Equal(int64(1)))
	})

	It("converts Go maps to immutable maps", func() {
		original := map[string]int64{"k": 1}
		immutable := values.FromGoMap(original)
		original["k"] = 99

		v, ok := immutable.Get("k")
		Expect(ok).To(BeTrue())
		Expect(v).To(Equal(int64(1)))
	})

	It("converts map[string]interface{} to immutable maps", func() {
		original := map[string]interface{}{"k": int64(1)}
		immutable, err := values.FromStringInterfaceMap(original)
		Expect(err).NotTo(HaveOccurred())
		original["k"] = int64(99)

		v, ok := immutable.Get("k")
		Expect(ok).To(BeTrue())
		Expect(v).To(Equal(int64(1)))
	})
})
