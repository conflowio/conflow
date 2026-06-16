// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package bind_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/pkg/conflow/bind"
	"github.com/conflowio/conflow/pkg/schema"
	"github.com/conflowio/conflow/pkg/values"
)

var _ = Describe("Isolation", func() {
	arraySchema := &schema.Array{Items: &schema.String{}}

	It("does not alias mutable slice input across multiple binds", func() {
		upstream := []interface{}{"shared"}
		bound1, err := bind.BindValue(arraySchema, upstream)
		Expect(err).NotTo(HaveOccurred())
		bound2, err := bind.BindValue(arraySchema, upstream)
		Expect(err).NotTo(HaveOccurred())

		Expect(bound1).NotTo(BeIdenticalTo(&upstream))
		Expect(bound2).NotTo(BeIdenticalTo(&upstream))
		Expect(bound1).NotTo(BeIdenticalTo(bound2))

		upstream[0] = "mutated"

		list1 := bound1.(*values.List[interface{}])
		list2 := bound2.(*values.List[interface{}])
		Expect(list1.At(0)).To(Equal("shared"))
		Expect(list2.At(0)).To(Equal("shared"))
	})

	It("shares immutable list pointer across binds", func() {
		upstream := values.NewList("shared")
		bound1, err := bind.BindValue(arraySchema, upstream)
		Expect(err).NotTo(HaveOccurred())
		bound2, err := bind.BindValue(arraySchema, upstream)
		Expect(err).NotTo(HaveOccurred())

		Expect(bound1).To(BeIdenticalTo(upstream))
		Expect(bound2).To(BeIdenticalTo(upstream))
		Expect(bound1).To(BeIdenticalTo(bound2))
	})
})
