package test

import (
	"fmt"

	"github.com/opsidian/basil/block"

	. "github.com/onsi/gomega"
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
)

func parseCtx(input string, registry block.Registry) *parsley.Context {
	f := text.NewFile("testfile", []byte(input))
	fs := parsley.NewFileSet(f)
	r := text.NewReader(f)
	ctx := parsley.NewContext(fs, r)
	if registry != nil {
		ctx.SetNodeTransformer(block.TransformNode(registry))
	}
	ctx.RegisterKeywords("true", "false", "nil", "map", "testkeyword")
	return ctx
}

func evalCtx() interface{} {
	return basil.NewContext(
		nil,
		basil.ContextConfig{
			VariableProvider: testVariableProvider,
			FunctionRegistry: &functionRegistry{},
			IDRegistry:       newIDRegistry(),
		},
	)
}

func ExpectParserToEvaluate(p parsley.Parser) func(string, interface{}) {
	return func(input string, expected interface{}) {
		val, err := parsley.Evaluate(parseCtx(input, nil), combinator.Sentence(p), evalCtx())

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
		res, err := parsley.Parse(parseCtx(input, nil), combinator.Sentence(p))

		Expect(err).To(HaveOccurred(), "input: %s", input)
		Expect(err).To(MatchError(fmt.Errorf("failed to parse the input: %s", expectedErr)), "input: %s", input)
		Expect(res).To(BeNil(), "input: %s", input)
	}
}

func ExpectParserToHaveEvalError(p parsley.Parser) func(string, error) {
	return func(input string, expectedErr error) {
		val, err := parsley.Evaluate(parseCtx(input, nil), combinator.Sentence(p), evalCtx())

		Expect(err).To(HaveOccurred(), "input: %s", input)
		Expect(err).To(MatchError(expectedErr), "input: %s", input)
		Expect(val).To(BeNil(), "input: %s", input)
	}
}

func ExpectParserToReturn(p parsley.Parser, input string, expected parsley.Node) {
	res, err := parsley.Parse(parseCtx(input, nil), combinator.Sentence(p))

	Expect(err).ToNot(HaveOccurred())

	node, ok := res.(parsley.NonTerminalNode)
	Expect(ok).To(BeTrue())
	actual := node.Children()[0]

	Expect(actual).To(BeEquivalentTo(expected))
}

func ExpectBlockToEvaluate(p parsley.Parser, registry block.Registry) func(string, interface{}, func(interface{}, interface{}, string)) {
	return func(input string, expected interface{}, compare func(interface{}, interface{}, string)) {
		block, err := parsley.Evaluate(parseCtx(input, registry), combinator.Sentence(p), evalCtx())
		Expect(err).ToNot(HaveOccurred(), "eval failed, input: %s", input)

		compare(block, expected, input)
	}
}

func ExpectBlockToHaveParseError(p parsley.Parser, registry block.Registry) func(string, error) {
	return func(input string, expectedErr error) {
		res, err := parsley.Parse(parseCtx(input, registry), combinator.Sentence(p))
		Expect(err).To(HaveOccurred(), "input: %s", input)
		Expect(err).To(MatchError(fmt.Errorf("failed to process the input: %s", expectedErr)), "input: %s", input)
		Expect(res).To(BeNil(), "input: %s", input)
	}
}

func ExpectBlockToHaveCheckError(p parsley.Parser, registry block.Registry) func(string, error) {
	return func(input string, expectedErr error) {
		parseCtx := parseCtx(input, registry)
		res, err := parsley.Parse(parseCtx, combinator.Sentence(p))
		Expect(err).ToNot(HaveOccurred(), "input: %s", input)

		err = parsley.StaticCheck(parseCtx, res, evalCtx())
		Expect(err).To(HaveOccurred(), "input: %s", input)
		Expect(err).To(MatchError(expectedErr), "input: %s", input)
	}
}

func ExpectBlockToHaveEvalError(p parsley.Parser, registry block.Registry) func(string, error) {
	return func(input string, expectedErr error) {
		_, err := parsley.Evaluate(parseCtx(input, registry), combinator.Sentence(p), evalCtx())
		Expect(err).To(MatchError(expectedErr), "input: %s", input)
	}
}

func ExpectBlockNodeToEvaluate(p parsley.Parser, registry block.Registry, block basil.Block, node basil.BlockNode) func(string, interface{}, func(interface{}, interface{}, string)) {
	return func(input string, expected interface{}, compare func(interface{}, interface{}, string)) {
		evalCtx := block.Context(evalCtx())

		block, err := node.Value(evalCtx)
		Expect(err).ToNot(HaveOccurred(), "create block failed, input: %s", input)

		compare(block, expected, input)
	}
}
