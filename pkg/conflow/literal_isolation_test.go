// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conflow_test

import (
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/text/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/schema"
	"github.com/conflowio/conflow/pkg/values"
)

var _ = Describe("Literal isolation", func() {
	intNode := func(v int64) parsley.Node {
		return terminal.NewIntegerNode(schema.IntegerValue(), v, parsley.NilPos, parsley.NilPos)
	}
	strNode := func(v string) parsley.Node {
		return terminal.NewStringNode(schema.StringValue(), v, parsley.NilPos, parsley.NilPos)
	}

	It("evaluates array literals to immutable lists", func() {
		node := conflow.NewArrayNode(
			[]parsley.Node{intNode(1), intNode(2)},
			parsley.NilPos,
			parsley.NilPos,
			nil,
		)

		value, err := node.Value(nil)
		Expect(err).NotTo(HaveOccurred())

		list, ok := value.(*values.List[interface{}])
		Expect(ok).To(BeTrue())
		Expect(list.Len()).To(Equal(2))
		Expect(list.At(0)).To(Equal(int64(1)))
		Expect(list.At(1)).To(Equal(int64(2)))
	})

	It("evaluates empty array literals to empty immutable lists", func() {
		node := conflow.NewArrayNode(nil, parsley.NilPos, parsley.NilPos, nil)

		value, err := node.Value(nil)
		Expect(err).NotTo(HaveOccurred())

		list, ok := value.(*values.List[interface{}])
		Expect(ok).To(BeTrue())
		Expect(list.Len()).To(Equal(0))
	})

	It("evaluates nested array literals to nested immutable lists", func() {
		inner := conflow.NewArrayNode(
			[]parsley.Node{intNode(1), intNode(2)},
			parsley.NilPos,
			parsley.NilPos,
			nil,
		)
		outer := conflow.NewArrayNode(
			[]parsley.Node{inner, conflow.NewArrayNode([]parsley.Node{intNode(3), strNode("foo")}, parsley.NilPos, parsley.NilPos, nil)},
			parsley.NilPos,
			parsley.NilPos,
			nil,
		)

		value, err := outer.Value(nil)
		Expect(err).NotTo(HaveOccurred())

		list, ok := value.(*values.List[interface{}])
		Expect(ok).To(BeTrue())
		Expect(list.Len()).To(Equal(2))

		first, ok := list.At(0).(*values.List[interface{}])
		Expect(ok).To(BeTrue())
		Expect(first.At(1)).To(Equal(int64(2)))

		second, ok := list.At(1).(*values.List[interface{}])
		Expect(ok).To(BeTrue())
		Expect(second.At(1)).To(Equal("foo"))
	})

	It("does not allow mutation through Elems copy", func() {
		node := conflow.NewArrayNode([]parsley.Node{strNode("shared")}, parsley.NilPos, parsley.NilPos, nil)

		value, err := node.Value(nil)
		Expect(err).NotTo(HaveOccurred())

		list := value.(*values.List[interface{}])
		copy := list.Elems()
		copy[0] = "mutated"
		Expect(list.At(0)).To(Equal("shared"))
	})

	It("evaluates map literals to immutable maps", func() {
		node := conflow.NewMapNode(
			[]string{"a", "b"},
			[]parsley.Node{strNode("one"), intNode(2)},
			parsley.NilPos,
			parsley.NilPos,
			nil,
		)

		value, err := node.Value(nil)
		Expect(err).NotTo(HaveOccurred())

		m, ok := value.(*values.Map[string, interface{}])
		Expect(ok).To(BeTrue())
		Expect(m.Len()).To(Equal(2))

		v, found := m.Get("a")
		Expect(found).To(BeTrue())
		Expect(v).To(Equal("one"))

		v, found = m.Get("b")
		Expect(found).To(BeTrue())
		Expect(v).To(Equal(int64(2)))
	})

	It("evaluates empty map literals to empty immutable maps", func() {
		node := conflow.NewMapNode(nil, nil, parsley.NilPos, parsley.NilPos, nil)

		value, err := node.Value(nil)
		Expect(err).NotTo(HaveOccurred())

		m, ok := value.(*values.Map[string, interface{}])
		Expect(ok).To(BeTrue())
		Expect(m.Len()).To(Equal(0))
	})

	It("does not allow mutation through GoMap copy", func() {
		node := conflow.NewMapNode(
			[]string{"key"},
			[]parsley.Node{strNode("shared")},
			parsley.NilPos,
			parsley.NilPos,
			nil,
		)

		value, err := node.Value(nil)
		Expect(err).NotTo(HaveOccurred())

		m := value.(*values.Map[string, interface{}])
		copy := m.GoMap()
		copy["key"] = "mutated"
		v, _ := m.Get("key")
		Expect(v).To(Equal("shared"))
	})
})
