package test

import (
	"fmt"

	"github.com/opsidian/ocl/block"

	. "github.com/onsi/gomega"
	"github.com/opsidian/ocl/ocl"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
)

func parseCtx(input string) *parsley.Context {
	f := text.NewFile("testfile", []byte(input))
	fs := parsley.NewFileSet(f)
	r := text.NewReader(f)
	ctx := parsley.NewContext(fs, r)
	ctx.RegisterKeywords("true", "false", "nil", "map", "testkeyword")
	return ctx
}

func evalCtx() interface{} {
	return ocl.NewContext(
		nil,
		ocl.ContextConfig{
			VariableProvider: testVariableProvider,
			FunctionRegistry: &functionRegistry{},
			BlockRegistry: block.Registry{
				"testblock": ocl.BlockFactoryCreatorFunc(NewTestBlockFactory),
			},
			IDRegistry: newIDRegistry(),
		},
	)
}

func ExpectParserToEvaluate(p parsley.Parser) func(string, interface{}) {
	return func(input string, expected interface{}) {
		val, err := parsley.Evaluate(parseCtx(input), combinator.Sentence(p), evalCtx())

		Expect(err).ToNot(HaveOccurred(), "input: %s", input)

		if val != nil {
			Expect(val).To(Equal(expected), "input: %s", input)
		} else {
			Expect(val).To(BeNil(), "input: %s", input)
		}
	}
}

func ExpectParserToHaveParseError(p parsley.Parser) func(string, error) {
	return func(input string, expectedErr error) {
		res, err := parsley.Parse(parseCtx(input), combinator.Sentence(p))

		Expect(err).To(HaveOccurred(), "input: %s", input)
		Expect(err).To(MatchError(fmt.Errorf("failed to parse the input: %s", expectedErr)), "input: %s", input)
		Expect(res).To(BeNil(), "input: %s", input)
	}
}

func ExpectParserToHaveEvalError(p parsley.Parser) func(string, error) {
	return func(input string, expectedErr error) {
		val, err := parsley.Evaluate(parseCtx(input), combinator.Sentence(p), evalCtx())

		Expect(err).To(HaveOccurred(), "input: %s", input)
		Expect(err).To(MatchError(expectedErr), "input: %s", input)
		Expect(val).To(BeNil(), "input: %s", input)
	}
}

func ExpectParserToReturn(p parsley.Parser, input string, expected parsley.Node) {
	res, err := parsley.Parse(parseCtx(input), combinator.Sentence(p))

	Expect(err).ToNot(HaveOccurred())

	node, ok := res.(*ast.NonTerminalNode)
	Expect(ok).To(BeTrue())
	actual := node.Children()[0]

	Expect(actual).To(BeEquivalentTo(expected))
}

func ExpectBlockToEvaluate(p parsley.Parser) func(string, *TestBlock) {
	return func(input string, expected *TestBlock) {
		evalCtx := evalCtx()
		val, err := parsley.Evaluate(parseCtx(input), combinator.Sentence(p), evalCtx)
		Expect(err).ToNot(HaveOccurred(), "eval failed, input: %s", input)

		block, blockCtx, err := val.(ocl.BlockFactory).CreateBlock(evalCtx)
		Expect(err).ToNot(HaveOccurred(), "create block failed, input: %s", input)

		err = val.(ocl.BlockFactory).EvalBlock(blockCtx, "default", block)
		Expect(err).ToNot(HaveOccurred(), "eval block failed, input: %s", input)

		testBlock := block.(*TestBlock)

		expectBlockToEqual(testBlock, expected, input)
	}
}

func expectBlockToEqual(b1 *TestBlock, b2 *TestBlock, input string) {
	Expect(b1.IDField).To(Equal(b2.IDField), "IDField does not match, input: %s", input)
	Expect(b1.Value).To(Equal(b2.Value), "Value does not match, input: %s", input)
	Expect(b1.FieldString).To(Equal(b2.FieldString), "FieldString does not match, input: %s", input)
	Expect(b1.FieldInt).To(Equal(b2.FieldInt), "FieldInt does not match, input: %s", input)
	Expect(b1.FieldFloat).To(Equal(b2.FieldFloat), "FieldFloat does not match, input: %s", input)
	Expect(b1.FieldBool).To(Equal(b2.FieldBool), "FieldBool does not match, input: %s", input)
	Expect(b1.FieldArray).To(Equal(b2.FieldArray), "FieldArray does not match, input: %s", input)
	Expect(b1.FieldMap).To(Equal(b2.FieldMap), "FieldMap does not match, input: %s", input)
	Expect(b1.FieldTimeDuration).To(Equal(b2.FieldTimeDuration), "FieldTimeDuration does not match, input: %s", input)
	Expect(b1.FieldCustomName).To(Equal(b2.FieldCustomName), "FieldCustomName does not match, input: %s", input)

	Expect(len(b1.Blocks)).To(Equal(len(b2.Blocks)), "child block count does not match, input: %s", input)

	for i, block := range b1.Blocks {
		expectBlockToEqual(block, b2.Blocks[i], input)
	}

}
