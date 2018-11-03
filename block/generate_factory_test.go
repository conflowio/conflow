package block_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/ocl/block"
	"github.com/opsidian/ocl/block/fixtures"
	"github.com/opsidian/ocl/ocl"
	"github.com/opsidian/ocl/parser"
	"github.com/opsidian/ocl/test"
)

var _ = Describe("GenerateFactory", func() {

	var blockRegistry = block.Registry{
		"block_simple":               ocl.BlockFactoryCreatorFunc(fixtures.NewBlockSimpleFactory),
		"block_value_required":       ocl.BlockFactoryCreatorFunc(fixtures.NewBlockValueRequiredFactory),
		"block_with_block":           ocl.BlockFactoryCreatorFunc(fixtures.NewBlockWithBlockFactory),
		"block_with_block_interface": ocl.BlockFactoryCreatorFunc(fixtures.NewBlockWithBlockInterfaceFactory),
		"block_with_factory":         ocl.BlockFactoryCreatorFunc(fixtures.NewBlockWithFactoryFactory),
	}

	Context("fixtures/block_simple.go", func() {
		It("should parse the input", func() {
			test.ExpectBlockToEvaluate(parser.Block(), blockRegistry)(
				`block_simple foo`,
				&fixtures.BlockSimple{IDField: "foo"},
				func(b1i interface{}, b2i interface{}, input string) {
					Expect(b1i).To(Equal(b2i), "input was %s", input)
				},
			)
		})

		It("should parse the input in short format", func() {
			test.ExpectBlockToEvaluate(parser.Block(), blockRegistry)(
				`block_simple foo "bar"`,
				&fixtures.BlockSimple{IDField: "foo", Value: "bar"},
				func(b1i interface{}, b2i interface{}, input string) {
					Expect(b1i).To(Equal(b2i), "input was %s", input)
				},
			)
		})

		It("should not parse fields with nil values", func() {
			test.ExpectBlockToEvaluate(parser.Block(), blockRegistry)(
				`block_simple foo {
					value = nil
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
			test.ExpectBlockToEvaluate(parser.Block(), blockRegistry)(
				`block_value_required foo "bar"`,
				&fixtures.BlockValueRequired{IDField: "foo", Value: "bar"},
				func(b1i interface{}, b2i interface{}, input string) {
					Expect(b1i).To(Equal(b2i), "input was %s", input)
				},
			)
		})

		It("should parse the input in short f", func() {
			test.ExpectBlockToEvaluate(parser.Block(), blockRegistry)(
				`block_value_required foo {
					value = "bar"
				}`,
				&fixtures.BlockValueRequired{IDField: "foo", Value: "bar"},
				func(b1i interface{}, b2i interface{}, input string) {
					Expect(b1i).To(Equal(b2i), "input was %s", input)
				},
			)
		})
	})

	Context("fixtures/block_with_block.go", func() {
		It("should parse the input", func() {
			test.ExpectBlockToEvaluate(parser.Block(), blockRegistry)(
				`block_with_block foo {
					block_simple bar
				}`,
				&fixtures.BlockWithBlock{IDField: "foo", Blocks: []*fixtures.BlockSimple{
					{IDField: "bar"},
				}},
				func(b1i interface{}, b2i interface{}, input string) {
					b1 := b1i.(*fixtures.BlockWithBlock)
					b2 := b2i.(*fixtures.BlockWithBlock)
					Expect(b1.IDField).To(Equal(b2.IDField), "IDField does not match, input was %s", input)
					Expect(b1.Blocks).To(Equal(b2.Blocks), "Blocks does not match, input was %s", input)
				},
			)
		})
	})

	Context("fixtures/block_with_block_interface.go", func() {
		It("should parse the input", func() {
			test.ExpectBlockToEvaluate(parser.Block(), blockRegistry)(
				`block_with_block_interface foo {
					block_simple bar
				}`,
				&fixtures.BlockWithBlockInterface{IDField: "foo", Blocks: []fixtures.BlockInterface{
					&fixtures.BlockSimple{IDField: "bar"},
				}},
				func(b1i interface{}, b2i interface{}, input string) {
					b1 := b1i.(*fixtures.BlockWithBlockInterface)
					b2 := b2i.(*fixtures.BlockWithBlockInterface)
					Expect(b1.IDField).To(Equal(b2.IDField), "IDField does not match, input was %s", input)
					Expect(b1.Blocks).To(Equal(b2.Blocks), "Blocks does not match, input was %s", input)
				},
			)
		})
	})

	Context("fixtures/block_with_factory.go", func() {
		It("should parse the input", func() {
			test.ExpectBlockToEvaluate(parser.Block(), blockRegistry)(
				`block_with_factory foo {
					block_simple bar
				}`,
				&fixtures.BlockWithFactory{IDField: "foo"},
				func(b1i interface{}, b2i interface{}, input string) {
					b1 := b1i.(*fixtures.BlockWithFactory)
					b2 := b2i.(*fixtures.BlockWithFactory)
					Expect(b1.IDField).To(Equal(b2.IDField), "IDField does not match, input was %s", input)
					test.ExpectBlockFactoryToEvaluate(parser.Block(), blockRegistry, b1, b1.BlockFactories[0])(
						input,
						&fixtures.BlockSimple{IDField: "bar"},
						func(b1i interface{}, b2i interface{}, input string) {
							b1 := b1i.(*fixtures.BlockSimple)
							b2 := b2i.(*fixtures.BlockSimple)
							Expect(b1.IDField).To(Equal(b2.IDField), "IDField does not match, input was %s", input)
						},
					)
				},
			)
		})
	})
})
