package test

import (
	"context"
	"fmt"
	"time"

	"github.com/opsidian/basil/block"
	"github.com/opsidian/basil/function"

	. "github.com/onsi/gomega"
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
)

func parseCtx(
	input string,
	blockRegistry parsley.NodeTransformerRegistry,
	functionRegistry parsley.NodeTransformerRegistry,
) *parsley.Context {
	f := text.NewFile("testfile", []byte(input))
	fs := parsley.NewFileSet(f)
	r := text.NewReader(f)
	ctx := parsley.NewContext(fs, r)
	ctx.EnableStaticCheck()
	ctx.EnableTransformation()
	ctx.RegisterKeywords("true", "false", "nil", "map", "testkeyword")
	ctx.SetUserContext(parseUserCtx(blockRegistry, functionRegistry))
	return ctx
}

func parseUserCtx(
	blockRegistry parsley.NodeTransformerRegistry,
	functionRegistry parsley.NodeTransformerRegistry,
) *basil.ParseContext {
	if functionRegistry == nil {
		functionRegistry = function.InterpreterRegistry{
			"test_func0": TestFunc0Interpreter{},
			"test_func1": TestFunc1Interpreter{},
			"test_func2": TestFunc2Interpreter{},
		}
	}
	return basil.NewParseContext(
		blockRegistry,
		functionRegistry,
		newIDRegistry(),
		block.NewNodeRegistry(),
	)
}

func evalUserCtx() *basil.EvalContext {
	blockContainerRegistry := block.NewContainerRegistry()

	bc1 := &TestBlock{
		IDField:           "foo",
		FieldInt:          1,
		FieldFloat:        1.2,
		FieldString:       "string",
		FieldBool:         true,
		FieldTimeDuration: 90 * time.Minute,
		FieldArray:        []interface{}{2, 2.3, "str"},
		FieldMap: map[string]interface{}{
			"key1": "value1",
		},
		FieldCustomName: "customvalue",
	}

	bc2 := &TestBlock{
		IDField: "bar",
	}

	root := &TestBlock{
		IDField: basil.ID("root"),
		Blocks: []*TestBlock{
			bc1,
			bc2,
		},
	}

	blockContainerRegistry.AddBlockContainer(block.NewContainer(root, TestBlockInterpreter{}))
	blockContainerRegistry.AddBlockContainer(block.NewContainer(bc1, TestBlockInterpreter{}))
	blockContainerRegistry.AddBlockContainer(block.NewContainer(bc2, TestBlockInterpreter{}))

	return basil.NewEvalContext(context.Background(), "userctx", blockContainerRegistry)
}

func ExpectParserToEvaluate(p parsley.Parser) func(string, interface{}) {
	return func(input string, expected interface{}) {
		node, parseErr := parsley.Parse(parseCtx(input, nil, nil), combinator.Sentence(p))
		Expect(parseErr).ToNot(HaveOccurred(), "input: %s", input)

		value, evalErr := node.Value(evalUserCtx())
		Expect(evalErr).ToNot(HaveOccurred(), "input: %s", input)

		if value != nil {
			Expect(value).To(Equal(expected), "input: %s", input)
		} else {
			Expect(value).To(BeNil(), "input: %s", input)
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
		parseCtx := parseCtx(input, nil, nil)
		node, parseErr := parsley.Parse(parseCtx, combinator.Sentence(p))
		Expect(parseErr).ToNot(HaveOccurred(), "input: %s", input)

		value, evalErr := node.Value(evalUserCtx())
		Expect(evalErr).To(HaveOccurred())
		Expect(parseCtx.FileSet().ErrorWithPosition(evalErr)).To(HaveOccurred(), "input: %s", input)
		Expect(value).To(MatchError(expectedErr), "input: %s", input)
		Expect(value).To(BeNil(), "input: %s", input)
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

func ExpectBlockToEvaluate(p parsley.Parser, registry parsley.NodeTransformerRegistry) func(string, interface{}, func(interface{}, interface{}, string)) {
	return func(input string, expected interface{}, compare func(interface{}, interface{}, string)) {
		node, parseErr := parsley.Parse(parseCtx(input, registry, nil), combinator.Sentence(p))
		Expect(parseErr).ToNot(HaveOccurred(), "input: %s", input)

		value, evalErr := node.Value(evalUserCtx())
		Expect(evalErr).ToNot(HaveOccurred(), "eval failed, input: %s", input)

		compare(value, expected, input)
	}
}

func ExpectBlockToHaveParseError(p parsley.Parser, registry parsley.NodeTransformerRegistry) func(string, error) {
	return func(input string, expectedErr error) {
		node, parseErr := parsley.Parse(parseCtx(input, registry, nil), combinator.Sentence(p))
		Expect(parseErr).To(HaveOccurred(), "input: %s", input)
		Expect(parseErr).To(MatchError(expectedErr), "input: %s", input)
		Expect(node).To(BeNil(), "input: %s", input)
	}
}

func ExpectBlockToHaveEvalError(p parsley.Parser, registry parsley.NodeTransformerRegistry) func(string, error) {
	return func(input string, expectedErr error) {
		parseCtx := parseCtx(input, registry, nil)
		node, parseErr := parsley.Parse(parseCtx, combinator.Sentence(p))
		Expect(parseErr).ToNot(HaveOccurred(), "input: %s", input)

		_, evalErr := node.Value(evalUserCtx())
		Expect(evalErr).To(HaveOccurred())
		Expect((parseCtx.FileSet().ErrorWithPosition(evalErr))).To(MatchError(expectedErr), "input: %s", input)
	}
}

func ExpectBlockNodeToEvaluate(p parsley.Parser, registry parsley.NodeTransformerRegistry, block basil.Block, node basil.BlockNode) func(string, interface{}, func(interface{}, interface{}, string)) {
	return func(input string, expected interface{}, compare func(interface{}, interface{}, string)) {
		// TODO: registry is not used
		block, err := node.Value(evalUserCtx())
		Expect(err).ToNot(HaveOccurred(), "create block failed, input: %s", input)

		compare(block, expected, input)
	}
}

func ExpectFunctionToEvaluate(p parsley.Parser, registry parsley.NodeTransformerRegistry) func(string, interface{}) {
	return func(input string, expected interface{}) {
		node, parseErr := parsley.Parse(parseCtx(input, nil, registry), combinator.Sentence(p))
		Expect(parseErr).ToNot(HaveOccurred(), "input: %s", input)

		value, evalErr := node.Value(evalUserCtx())
		Expect(evalErr).ToNot(HaveOccurred(), "eval failed, input: %s", input)
		switch expected.(type) {
		case int64, float64:
			Expect(value).To(BeNumerically("~", expected))
		case nil:
			Expect(value).To(BeNil())
		default:
			Expect(value).To(Equal(expected))
		}
	}
}

func ExpectFunctionToHaveParseError(p parsley.Parser, registry parsley.NodeTransformerRegistry) func(string, error) {
	return func(input string, expectedErr error) {
		res, err := parsley.Parse(parseCtx(input, nil, registry), combinator.Sentence(p))
		Expect(err).To(HaveOccurred(), "input: %s", input)
		Expect(err).To(MatchError(expectedErr), "input: %s", input)
		Expect(res).To(BeNil(), "input: %s", input)
	}
}

func ExpectFunctionToHaveEvalError(p parsley.Parser, registry parsley.NodeTransformerRegistry) func(string, error) {
	return func(input string, expectedErr error) {
		parseCtx := parseCtx(input, nil, registry)
		node, parseErr := parsley.Parse(parseCtx, combinator.Sentence(p))
		Expect(parseErr).ToNot(HaveOccurred(), "input: %s", input)

		value, evalErr := node.Value(evalUserCtx())
		Expect(evalErr).To(HaveOccurred(), "input: %s", input)
		Expect(parseCtx.FileSet().ErrorWithPosition(evalErr)).To(MatchError(expectedErr), "input: %s", input)
		Expect(value).To(BeNil(), "input: %s", input)
	}
}

func ExpectFunctionNode(p parsley.Parser, registry parsley.NodeTransformerRegistry) func(string, func(interface{}, parsley.Node)) {
	return func(input string, test func(interface{}, parsley.Node)) {
		ctx := parseCtx(input, nil, registry)
		node, err := parsley.Parse(ctx, combinator.Sentence(p))
		Expect(err).ToNot(HaveOccurred(), "input: %s", input)

		test(ctx.UserContext(), node)
	}
}
