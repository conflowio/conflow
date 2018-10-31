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
		"block_simple":                 ocl.BlockFactoryCreatorFunc(fixtures.NewBlockSimpleFactory),
		"block_with_block":             ocl.BlockFactoryCreatorFunc(fixtures.NewBlockWithBlockFactory),
		"block_with_block_interface":   ocl.BlockFactoryCreatorFunc(fixtures.NewBlockWithBlockInterfaceFactory),
		"block_with_factory":           ocl.BlockFactoryCreatorFunc(fixtures.NewBlockWithFactoryFactory),
		"block_with_factory_interface": ocl.BlockFactoryCreatorFunc(fixtures.NewBlockWithFactoryInterfaceFactory),
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
})
