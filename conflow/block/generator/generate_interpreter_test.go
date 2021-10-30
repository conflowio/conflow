// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/conflow/basil/block"
	"github.com/opsidian/conflow/basil/block/fixtures"
	"github.com/opsidian/conflow/conflow"
	"github.com/opsidian/conflow/parsers"
	"github.com/opsidian/conflow/test"
)

var _ = Describe("GenerateInterpreter", func() {

	p := parsers.Block(parsers.Expression())

	var registry = block.InterpreterRegistry{
		"block_simple":          fixtures.BlockSimpleInterpreter{},
		"block_value_required":  fixtures.BlockValueRequiredInterpreter{},
		"block_with_one_block":  fixtures.BlockWithOneBlockInterpreter{},
		"block_with_many_block": fixtures.BlockWithManyBlockInterpreter{},
		"block_with_default":    fixtures.BlockWithDefaultInterpreter{},
		"block_with_interface":  fixtures.BlockWithInterfaceInterpreter{},
	}

	Context("fixtures/block_simple.go", func() {
		It("should parse the input", func() {
			test.ExpectBlockToEvaluate(p, registry)(
				`foo block_simple`,
				&fixtures.BlockSimple{IDField: "foo"},
				func(b1i interface{}, b2i interface{}, input string) {
					Expect(b1i).To(Equal(b2i), "input was %s", input)
				},
			)
		})

		It("should parse the input in short format", func() {
			test.ExpectBlockToEvaluate(p, registry)(
				`foo block_simple "bar"`,
				&fixtures.BlockSimple{IDField: "foo", Value: "bar"},
				func(b1i interface{}, b2i interface{}, input string) {
					Expect(b1i).To(Equal(b2i), "input was %s", input)
				},
			)
		})

		It("should not parse fields with nil values", func() {
			test.ExpectBlockToEvaluate(p, registry)(
				`foo block_simple {
					value = null
				}`,
				&fixtures.BlockSimple{IDField: "foo"},
				func(b1i interface{}, b2i interface{}, input string) {
					Expect(b1i).To(Equal(b2i), "input was %s", input)
				},
			)
		})
	})

	Context("fixtures/block_value_required.go", func() {
		It("should parse the input in short format", func() {
			test.ExpectBlockToEvaluate(p, registry)(
				`foo block_value_required "bar"`,
				&fixtures.BlockValueRequired{IDField: "foo", Value: "bar"},
				func(b1i interface{}, b2i interface{}, input string) {
					Expect(b1i).To(Equal(b2i), "input was %s", input)
				},
			)
		})

		It("should parse the input in short f", func() {
			test.ExpectBlockToEvaluate(p, registry)(
				`foo block_value_required {
					value = "bar"
				}`,
				&fixtures.BlockValueRequired{IDField: "foo", Value: "bar"},
				func(b1i interface{}, b2i interface{}, input string) {
					Expect(b1i).To(Equal(b2i), "input was %s", input)
				},
			)
		})
	})

	Context("fixtures/block_with_one_block.go", func() {
		It("should parse the input", func() {
			test.ExpectBlockToEvaluate(p, registry)(
				`foo block_with_one_block {
					bar block_simple
				}`,
				&fixtures.BlockWithOneBlock{IDField: "foo", BlockSimple: &fixtures.BlockSimple{IDField: "bar"}},
				func(b1i interface{}, b2i interface{}, input string) {
					b1 := b1i.(*fixtures.BlockWithOneBlock)
					b2 := b2i.(*fixtures.BlockWithOneBlock)
					Expect(b1.IDField).To(Equal(b2.IDField), "IDProperty does not match, input was %s", input)
					Expect(b1.BlockSimple).To(Equal(b2.BlockSimple), "Blocks does not match, input was %s", input)
				},
			)
		})
	})

	Context("fixtures/block_with_many_block.go", func() {
		It("should parse the input", func() {
			test.ExpectBlockToEvaluate(p, registry)(
				`foo block_with_many_block {
					bar block_simple
					baz block_simple
				}`,
				&fixtures.BlockWithManyBlock{IDField: "foo", BlockSimple: []*fixtures.BlockSimple{
					{IDField: "bar"},
					{IDField: "baz"},
				}},

				func(b1i interface{}, b2i interface{}, input string) {
					b1 := b1i.(*fixtures.BlockWithManyBlock)
					b2 := b2i.(*fixtures.BlockWithManyBlock)
					Expect(b1.IDField).To(Equal(b2.IDField), "IDProperty does not match, input was %s", input)
					Expect(b1.BlockSimple).To(ConsistOf(b2.BlockSimple), "Blocks does not match, input was %s", input)
				},
			)
		})
	})

	Context("fixtures/block_with_default.go", func() {
		It("should parse the input", func() {
			test.ExpectBlockToEvaluate(p, registry)(
				`block_with_default {}`,
				&fixtures.BlockWithDefault{IDField: "0", Value: "foo"},
				func(b1i interface{}, b2i interface{}, input string) {
					b1 := b1i.(*fixtures.BlockWithDefault)
					b2 := b2i.(*fixtures.BlockWithDefault)
					Expect(b1.IDField).To(Equal(b2.IDField), "IDProperty does not match, input was %s", input)
					Expect(b1.Value).To(Equal(b2.Value), "Value does not match, input was %s", input)
				},
			)
		})
	})

	Context("fixtures/block_with_interface.go", func() {
		It("should parse the block", func() {
			test.ExpectBlockToEvaluate(p, registry)(
				`foo block_with_interface {
					bar block:block_simple
				}`,
				&fixtures.BlockWithInterface{IDField: "foo", Block: &fixtures.BlockSimple{IDField: "bar"}},
				func(b1i interface{}, b2i interface{}, input string) {
					b1 := b1i.(*fixtures.BlockWithInterface)
					b2 := b2i.(*fixtures.BlockWithInterface)
					Expect(b1.IDField).To(Equal(b2.IDField), "IDProperty does not match, input was %s", input)
					Expect(b1.Block).To(Equal(b2.Block), "Blocks does not match, input was %s", input)
				},
			)
		})

		It("should parse the blocks array", func() {
			test.ExpectBlockToEvaluate(p, registry)(
				`foo block_with_interface {
					bar blocks:block_simple
				}`,
				&fixtures.BlockWithInterface{IDField: "foo", Blocks: []conflow.Block{&fixtures.BlockSimple{IDField: "bar"}}},
				func(b1i interface{}, b2i interface{}, input string) {
					b1 := b1i.(*fixtures.BlockWithInterface)
					b2 := b2i.(*fixtures.BlockWithInterface)
					Expect(b1.IDField).To(Equal(b2.IDField), "IDProperty does not match, input was %s", input)
					Expect(b1.Blocks).To(Equal(b2.Blocks), "Blocks does not match, input was %s", input)
				},
			)
		})
	})
})
