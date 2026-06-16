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

var _ = Describe("Map", func() {
	It("clones input Go map so mutating original does not affect immutable map", func() {
		original := map[string]int64{"k": 1}
		immutable := values.NewMapFromGoMap(original)
		original["k"] = 99
		original["new"] = 2

		v, ok := immutable.Get("k")
		Expect(ok).To(BeTrue())
		Expect(v).To(Equal(int64(1)))
		Expect(immutable.Len()).To(Equal(1))
	})

	It("returns ok false for missing keys", func() {
		immutable := values.NewMapFromGoMap(map[string]int64{"k": 1})

		_, ok := immutable.Get("missing")
		Expect(ok).To(BeFalse())
	})

	It("returns a copy from GoMap that does not affect the map when mutated", func() {
		immutable := values.NewMapFromGoMap(map[string]int64{"k": 1})
		cp := immutable.GoMap()
		cp["k"] = 99

		v, ok := immutable.Get("k")
		Expect(ok).To(BeTrue())
		Expect(v).To(Equal(int64(1)))
	})
})

var _ = Describe("MapBuilder", func() {
	It("returns an independent immutable map from Freeze", func() {
		builder := values.NewMapBuilder[string, int64]()
		builder.Set("k", 1)
		frozen := builder.Freeze()

		builder.Set("k", 99)
		builder.Set("other", 2)

		v, ok := frozen.Get("k")
		Expect(ok).To(BeTrue())
		Expect(v).To(Equal(int64(1)))
		Expect(frozen.Len()).To(Equal(1))
	})
})
