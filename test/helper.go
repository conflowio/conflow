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

func evalCtx(blockRegistry block.Registry) interface{} {
	if blockRegistry == nil {
		blockRegistry = block.Registry{
			"testblock": ocl.BlockFactoryCreatorFunc(NewTestBlockFactory),
		}
	}
	return ocl.NewContext(
		nil,
		ocl.ContextConfig{
			VariableProvider: testVariableProvider,
			FunctionRegistry: &functionRegistry{},
			BlockRegistry:    blockRegistry,
			IDRegistry:       newIDRegistry(),
		},
	)
}

func ExpectParserToEvaluate(p parsley.Parser) func(string, interface{}) {
	return func(input string, expected interface{}) {
		val, err := parsley.Evaluate(parseCtx(input), combinator.Sentence(p), evalCtx(nil))

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
		val, err := parsley.Evaluate(parseCtx(input), combinator.Sentence(p), evalCtx(nil))

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

func ExpectBlockToEvaluate(p parsley.Parser, blockRegistry block.Registry) func(string, interface{}, func(interface{}, interface{}, string)) {
	return func(input string, expected interface{}, compare func(interface{}, interface{}, string)) {
		evalCtx := evalCtx(blockRegistry)
		val, err := parsley.Evaluate(parseCtx(input), combinator.Sentence(p), evalCtx)
		Expect(err).ToNot(HaveOccurred(), "eval failed, input: %s", input)

		block, blockCtx, err := val.(ocl.BlockFactory).CreateBlock(evalCtx)
		Expect(err).ToNot(HaveOccurred(), "create block failed, input: %s", input)

		err = val.(ocl.BlockFactory).EvalBlock(blockCtx, "default", block)
		Expect(err).ToNot(HaveOccurred(), "eval block failed, input: %s", input)

		compare(block, expected, input)
	}
}

func ExpectBlockToHaveEvalError(p parsley.Parser, blockRegistry block.Registry) func(string, error) {
	return func(input string, expectedErr error) {
		ctx := parseCtx(input)
		evalCtx := evalCtx(blockRegistry)
		val, evalErr := parsley.Evaluate(ctx, combinator.Sentence(p), evalCtx)
		Expect(evalErr).ToNot(HaveOccurred(), "eval failed, input: %s", input)

		block, blockCtx, err := val.(ocl.BlockFactory).CreateBlock(evalCtx)
		if err != nil {
			errWithPos := ctx.FileSet().ErrorWithPosition(err)
			Expect(errWithPos).To(MatchError(expectedErr), "input: %s", input)
			return
		}

		err = val.(ocl.BlockFactory).EvalBlock(blockCtx, "default", block)
		Expect(err).To(HaveOccurred(), "input: %s", input)
		errWithPos := ctx.FileSet().ErrorWithPosition(err)
		Expect(errWithPos).To(MatchError(expectedErr), "input: %s", input)
	}
}

func ExpectBlockFactoryToEvaluate(p parsley.Parser, blockRegistry block.Registry, block ocl.Block, blockFactory ocl.BlockFactory) func(string, interface{}, func(interface{}, interface{}, string)) {
	return func(input string, expected interface{}, compare func(interface{}, interface{}, string)) {
		evalCtx := block.Context(evalCtx(blockRegistry))

		block, blockCtx, err := blockFactory.CreateBlock(evalCtx)
		Expect(err).ToNot(HaveOccurred(), "create block failed, input: %s", input)

		err = blockFactory.EvalBlock(blockCtx, "default", block)
		Expect(err).ToNot(HaveOccurred(), "eval block failed, input: %s", input)

		compare(block, expected, input)
	}
}
