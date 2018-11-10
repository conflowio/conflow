package test

import (
	"fmt"

	"github.com/opsidian/basil/block"
	"github.com/opsidian/basil/function"

	. "github.com/onsi/gomega"
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
)

func parseCtx(input string, blockRegistry block.Registry, functionRegistry function.Registry) *parsley.Context {
	f := text.NewFile("testfile", []byte(input))
	fs := parsley.NewFileSet(f)
	r := text.NewReader(f)
	ctx := parsley.NewContext(fs, r)
	ctx.EnableStaticCheck()
	ctx.EnableTransformation()
	ctx.RegisterKeywords("true", "false", "nil", "map", "testkeyword")
	ctx.SetUserContext(userCtx(blockRegistry, functionRegistry))
	return ctx
}

func userCtx(blockRegistry block.Registry, functionRegistry function.Registry) interface{} {
	if functionRegistry == nil {
		functionRegistry = function.Registry{
			"test_func0": TestFunc0Interpreter{},
			"test_func1": TestFunc1Interpreter{},
			"test_func2": TestFunc2Interpreter{},
		}
	}
	return basil.NewContext(
		nil,
		basil.ContextConfig{
			VariableProvider: testVariableProvider,
			IDRegistry:       newIDRegistry(),
			BlockRegistry:    blockRegistry,
			FunctionRegistry: functionRegistry,
		},
	)
}

func ExpectParserToEvaluate(p parsley.Parser) func(string, interface{}) {
	return func(input string, expected interface{}) {
		val, err := parsley.Evaluate(parseCtx(input, nil, nil), combinator.Sentence(p))

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
		res, err := parsley.Parse(parseCtx(input, nil, nil), combinator.Sentence(p))

		Expect(err).To(HaveOccurred(), "input: %s", input)
		Expect(err).To(MatchError(fmt.Errorf("failed to parse the input: %s", expectedErr)), "input: %s", input)
		Expect(res).To(BeNil(), "input: %s", input)
	}
}

func ExpectParserToHaveStaticCheckError(p parsley.Parser) func(string, error) {
	return func(input string, expectedErr error) {
		res, err := parsley.Parse(parseCtx(input, nil, nil), combinator.Sentence(p))

		Expect(err).To(HaveOccurred(), "input: %s", input)
		Expect(err).To(MatchError(expectedErr), "input: %s", input)
		Expect(res).To(BeNil(), "input: %s", input)
	}
}

func ExpectParserToHaveEvalError(p parsley.Parser) func(string, error) {
	return func(input string, expectedErr error) {
		val, err := parsley.Evaluate(parseCtx(input, nil, nil), combinator.Sentence(p))

		Expect(err).To(HaveOccurred(), "input: %s", input)
		Expect(err).To(MatchError(expectedErr), "input: %s", input)
		Expect(val).To(BeNil(), "input: %s", input)
	}
}

func ExpectParserToReturn(p parsley.Parser, input string, expected parsley.Node) {
	res, err := parsley.Parse(parseCtx(input, nil, nil), combinator.Sentence(p))

	Expect(err).ToNot(HaveOccurred())

	node, ok := res.(parsley.NonTerminalNode)
	Expect(ok).To(BeTrue())
	actual := node.Children()[0]

	Expect(actual).To(BeEquivalentTo(expected))
}

func ExpectBlockToEvaluate(p parsley.Parser, registry block.Registry) func(string, interface{}, func(interface{}, interface{}, string)) {
	return func(input string, expected interface{}, compare func(interface{}, interface{}, string)) {
		block, err := parsley.Evaluate(parseCtx(input, registry, nil), combinator.Sentence(p))
		Expect(err).ToNot(HaveOccurred(), "eval failed, input: %s", input)

		compare(block, expected, input)
	}
}

func ExpectBlockToHaveParseError(p parsley.Parser, registry block.Registry) func(string, error) {
	return func(input string, expectedErr error) {
		res, err := parsley.Parse(parseCtx(input, registry, nil), combinator.Sentence(p))
		Expect(err).To(HaveOccurred(), "input: %s", input)
		Expect(err).To(MatchError(expectedErr), "input: %s", input)
		Expect(res).To(BeNil(), "input: %s", input)
	}
}

func ExpectBlockToHaveEvalError(p parsley.Parser, registry block.Registry) func(string, error) {
	return func(input string, expectedErr error) {
		_, err := parsley.Evaluate(parseCtx(input, registry, nil), combinator.Sentence(p))
		Expect(err).To(MatchError(expectedErr), "input: %s", input)
	}
}

func ExpectBlockNodeToEvaluate(p parsley.Parser, registry block.Registry, block basil.Block, node basil.BlockNode) func(string, interface{}, func(interface{}, interface{}, string)) {
	return func(input string, expected interface{}, compare func(interface{}, interface{}, string)) {
		userCtx := block.Context(userCtx(registry, nil))

		block, err := node.Value(userCtx)
		Expect(err).ToNot(HaveOccurred(), "create block failed, input: %s", input)

		compare(block, expected, input)
	}
}

func ExpectFunctionToEvaluate(p parsley.Parser, registry function.Registry) func(string, interface{}) {
	return func(input string, expected interface{}) {
		res, err := parsley.Evaluate(parseCtx(input, nil, registry), combinator.Sentence(p))
		Expect(err).ToNot(HaveOccurred(), "eval failed, input: %s", input)
		switch expected.(type) {
		case int64, float64:
			Expect(res).To(BeNumerically("~", expected))
		case nil:
			Expect(res).To(BeNil())
		default:
			Expect(res).To(Equal(expected))
		}
	}
}

func ExpectFunctionToHaveParseError(p parsley.Parser, registry function.Registry) func(string, error) {
	return func(input string, expectedErr error) {
		res, err := parsley.Parse(parseCtx(input, nil, registry), combinator.Sentence(p))
		Expect(err).To(HaveOccurred(), "input: %s", input)
		Expect(err).To(MatchError(expectedErr), "input: %s", input)
		Expect(res).To(BeNil(), "input: %s", input)
	}
}

func ExpectFunctionToHaveEvalError(p parsley.Parser, registry function.Registry) func(string, error) {
	return func(input string, expectedErr error) {
		val, err := parsley.Evaluate(parseCtx(input, nil, registry), combinator.Sentence(p))
		Expect(err).To(HaveOccurred(), "input: %s", input)
		Expect(err).To(MatchError(expectedErr), "input: %s", input)
		Expect(val).To(BeNil(), "input: %s", input)
	}
}

func ExpectFunctionNode(p parsley.Parser, registry function.Registry) func(string, func(interface{}, parsley.Node)) {
	return func(input string, test func(interface{}, parsley.Node)) {
		ctx := parseCtx(input, nil, registry)
		node, err := parsley.Parse(ctx, combinator.Sentence(p))
		Expect(err).ToNot(HaveOccurred(), "input: %s", input)

		test(ctx.UserContext(), node)
	}
}
