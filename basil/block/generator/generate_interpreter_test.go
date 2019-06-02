package generator_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/basil/basil/block"
	"github.com/opsidian/basil/basil/block/fixtures"
	"github.com/opsidian/basil/parser"
	"github.com/opsidian/basil/test"
)

var _ = Describe("GenerateInterpreter", func() {

	p := parser.Block(parser.Expression())

	var registry = block.InterpreterRegistry{
		"block_simple":               fixtures.BlockSimpleInterpreter{},
		"block_value_required":       fixtures.BlockValueRequiredInterpreter{},
		"block_with_one_block":       fixtures.BlockWithOneBlockInterpreter{},
		"block_with_many_block":      fixtures.BlockWithManyBlockInterpreter{},
		"block_with_block_interface": fixtures.BlockWithBlockInterfaceInterpreter{},
		"block_with_reference":       fixtures.BlockWithReferenceInterpreter{},
	}

	Context("fixtures/block_simple.go", func() {
		It("should parse the input", func() {
			test.ExpectBlockToEvaluate(p, registry)(
				`block_simple foo`,
				&fixtures.BlockSimple{IDField: "foo"},
				func(b1i interface{}, b2i interface{}, input string) {
					Expect(b1i).To(Equal(b2i), "input was %s", input)
				},
			)
		})

		It("should parse the input in short format", func() {
			test.ExpectBlockToEvaluate(p, registry)(
				`block_simple foo "bar"`,
				&fixtures.BlockSimple{IDField: "foo", Value: "bar"},
				func(b1i interface{}, b2i interface{}, input string) {
					Expect(b1i).To(Equal(b2i), "input was %s", input)
				},
			)
		})

		It("should not parse fields with nil values", func() {
			test.ExpectBlockToEvaluate(p, registry)(
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
			test.ExpectBlockToEvaluate(p, registry)(
				`block_value_required foo "bar"`,
				&fixtures.BlockValueRequired{IDField: "foo", Value: "bar"},
				func(b1i interface{}, b2i interface{}, input string) {
					Expect(b1i).To(Equal(b2i), "input was %s", input)
				},
			)
		})

		It("should parse the input in short f", func() {
			test.ExpectBlockToEvaluate(p, registry)(
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

	Context("fixtures/block_with_one_block.go", func() {
		It("should parse the input", func() {
			test.ExpectBlockToEvaluate(p, registry)(
				`block_with_one_block foo {
					block_simple bar
				}`,
				&fixtures.BlockWithOneBlock{IDField: "foo", BlockSimple: &fixtures.BlockSimple{IDField: "bar"}},
				func(b1i interface{}, b2i interface{}, input string) {
					b1 := b1i.(*fixtures.BlockWithOneBlock)
					b2 := b2i.(*fixtures.BlockWithOneBlock)
					Expect(b1.IDField).To(Equal(b2.IDField), "IDField does not match, input was %s", input)
					Expect(b1.BlockSimple).To(Equal(b2.BlockSimple), "Blocks does not match, input was %s", input)
				},
			)
		})
	})

	Context("fixtures/block_with_many_block.go", func() {
		It("should parse the input", func() {
			test.ExpectBlockToEvaluate(p, registry)(
				`block_with_many_block foo {
					block_simple bar
					block_simple baz
				}`,
				&fixtures.BlockWithManyBlock{IDField: "foo", BlockSimple: []*fixtures.BlockSimple{
					{IDField: "bar"},
					{IDField: "baz"},
				}},
				func(b1i interface{}, b2i interface{}, input string) {
					b1 := b1i.(*fixtures.BlockWithManyBlock)
					b2 := b2i.(*fixtures.BlockWithManyBlock)
					Expect(b1.IDField).To(Equal(b2.IDField), "IDField does not match, input was %s", input)
					Expect(b1.BlockSimple).To(Equal(b2.BlockSimple), "Blocks does not match, input was %s", input)
				},
			)
		})
	})

	Context("fixtures/block_with_block_interface.go", func() {
		It("should parse the input", func() {
			test.ExpectBlockToEvaluate(p, registry)(
				`block_with_block_interface foo {
					block_simple bar
				}`,
				&fixtures.BlockWithBlockInterface{IDField: "foo", BlockSimple: []fixtures.BlockInterface{
					&fixtures.BlockSimple{IDField: "bar"},
				}},
				func(b1i interface{}, b2i interface{}, input string) {
					b1 := b1i.(*fixtures.BlockWithBlockInterface)
					b2 := b2i.(*fixtures.BlockWithBlockInterface)
					Expect(b1.IDField).To(Equal(b2.IDField), "IDField does not match, input was %s", input)
					Expect(b1.BlockSimple).To(Equal(b2.BlockSimple), "Blocks does not match, input was %s", input)
				},
			)
		})
	})

	Context("fixtures/block_with_reference.go", func() {
		It("should parse the input", func() {
			test.ExpectBlockToHaveEvalError(p, registry)(
				`block_with_reference foo {}`,
				MatchError(errors.New("\"foo\" is referencing a non-existing block at testfile:1:22")),
			)
		})
	})
})
