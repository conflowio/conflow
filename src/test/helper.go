// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package test

import (
	"context"
	"fmt"

	"github.com/conflowio/parsley/combinator"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/text"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"

	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/conflow/block"
	"github.com/conflowio/conflow/src/conflow/conflowfakes"
	"github.com/conflowio/conflow/src/conflow/function"
	"github.com/conflowio/conflow/src/directives"
	"github.com/conflowio/conflow/src/loggers/zerolog"
)

func ParseCtx(
	input string,
	blockRegistry parsley.NodeTransformerRegistry,
	functionRegistry parsley.NodeTransformerRegistry,
) *parsley.Context {
	if functionRegistry == nil {
		functionRegistry = function.InterpreterRegistry{
			"test.func0": TestFunc0Interpreter{},
			"test.func1": TestFunc1Interpreter{},
			"test.func2": TestFunc2Interpreter{},
		}
	}

	testBlockNode := block.NewNode(
		conflow.NewIDNode("test", conflow.ClassifierNone, parsley.NilPos, parsley.NilPos),
		conflow.NewNameNode(nil, nil, conflow.NewIDNode("testblock", conflow.ClassifierNone, parsley.NilPos, parsley.NilPos)),
		nil,
		nil,
		"TESTBLOCK",
		nil,
		parsley.NilPos,
		BlockInterpreter{},
		nil,
	)

	f := text.NewFile("testfile", []byte(input))
	fs := parsley.NewFileSet(f)
	r := text.NewReader(f)

	directiveTransformerRegistry := directives.DefaultRegistry().
		Register("testdirective", DirectiveInterpreter{}).
		Register("testdirective2", DirectiveInterpreter{})

	parseCtx := conflow.NewParseContext(fs, newIDRegistry(), directiveTransformerRegistry).New(conflow.ParseContextOverride{
		BlockTransformerRegistry:    blockRegistry,
		FunctionTransformerRegistry: functionRegistry,
	})
	_ = parseCtx.AddBlockNode(testBlockNode)

	ctx := parsley.NewContext(fs, r)
	ctx.EnableStaticCheck()
	ctx.EnableTransformation()
	ctx.RegisterKeywords("map", "testkeyword")
	ctx.SetUserContext(parseCtx)

	return ctx
}

func EvalUserCtx() *conflow.EvalContext {
	testBlock :=
		&Block{
			FieldString: "bar",
			FieldMap: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
			FieldArray: []interface{}{
				"value1",
				"value2",
			},
			FieldInt: int64(1),
		}

	testBlockContainer := &conflowfakes.FakeBlockContainer{}
	testBlockContainer.ParamCalls(func(name conflow.ID) interface{} {
		return BlockInterpreter{}.Param(testBlock, name)
	})

	containers := map[conflow.ID]conflow.BlockContainer{
		"test": testBlockContainer,
	}

	evalCtx := conflow.NewEvalContext(
		context.Background(),
		"userCtx",
		zerolog.NewDisabledLogger(),
		&Scheduler{},
		containers,
	)

	return evalCtx
}

func ExpectParserToEvaluate(p parsley.Parser) func(string, interface{}) {
	return func(input string, expected interface{}) {
		node, parseErr := parsley.Parse(ParseCtx(input, nil, nil), combinator.Sentence(p))
		Expect(parseErr).ToNot(HaveOccurred(), "input: %s", input)

		value, evalErr := parsley.EvaluateNode(EvalUserCtx(), node)
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
		res, err := parsley.Parse(ParseCtx(input, nil, nil), combinator.Sentence(p))

		Expect(err).To(HaveOccurred(), "input: %s", input)
		Expect(err).To(MatchError(fmt.Errorf("failed to parse the input: %w", expectedErr)), "input: %s", input)
		Expect(res).To(BeNil(), "input: %s", input)
	}
}

func ExpectParserToHaveStaticCheckError(p parsley.Parser) func(string, error) {
	return func(input string, expectedErr error) {
		res, err := parsley.Parse(ParseCtx(input, nil, nil), combinator.Sentence(p))

		Expect(err).To(HaveOccurred(), "input: %s", input)
		Expect(err).To(MatchError(expectedErr), "input: %s", input)
		Expect(res).To(BeNil(), "input: %s", input)
	}
}

func ExpectParserToHaveEvalError(p parsley.Parser) func(string, error) {
	return func(input string, expectedErr error) {
		parseCtx := ParseCtx(input, nil, nil)
		node, parseErr := parsley.Parse(parseCtx, combinator.Sentence(p))
		Expect(parseErr).ToNot(HaveOccurred(), "input: %s", input)

		value, evalErr := parsley.EvaluateNode(EvalUserCtx(), node)
		Expect(evalErr).To(HaveOccurred())
		Expect(parseCtx.FileSet().ErrorWithPosition(evalErr)).To(MatchError(expectedErr), "input: %s", input)
		Expect(value).To(BeNil(), "input: %s", input)
	}
}

func ExpectParserToReturn(p parsley.Parser, input string, expected parsley.Node) {
	res, err := parsley.Parse(ParseCtx(input, nil, nil), combinator.Sentence(p))
	Expect(err).ToNot(HaveOccurred())

	node, ok := res.(parsley.NonTerminalNode)
	Expect(ok).To(BeTrue())
	actual := node.Children()[0]

	Expect(actual).To(BeEquivalentTo(expected))
}

func ExpectBlockToEvaluate(p parsley.Parser, registry parsley.NodeTransformerRegistry) func(string, interface{}, func(interface{}, interface{}, string)) {
	return func(input string, expected interface{}, compare func(interface{}, interface{}, string)) {
		node, parseErr := parsley.Parse(ParseCtx(input, registry, nil), combinator.Sentence(p))
		Expect(parseErr).ToNot(HaveOccurred(), "input: %s", input)

		value, evalErr := parsley.EvaluateNode(EvalUserCtx(), node)
		Expect(evalErr).ToNot(HaveOccurred(), "eval failed, input: %s", input)

		compare(value, expected, input)
	}
}

func ExpectBlockToHaveParseError(p parsley.Parser, registry parsley.NodeTransformerRegistry) func(string, types.GomegaMatcher) {
	return func(input string, errMatcher types.GomegaMatcher) {
		node, parseErr := parsley.Parse(ParseCtx(input, registry, nil), combinator.Sentence(p))
		Expect(parseErr).To(HaveOccurred(), "input: %s", input)
		Expect(parseErr).To(errMatcher, "input: %s", input)
		Expect(node).To(BeNil(), "input: %s", input)
	}
}

func ExpectBlockToHaveStaticCheckError(p parsley.Parser, registry parsley.NodeTransformerRegistry) func(string, types.GomegaMatcher) {
	return func(input string, errMatcher types.GomegaMatcher) {
		node, parseErr := parsley.Parse(ParseCtx(input, registry, nil), combinator.Sentence(p))
		Expect(parseErr).To(HaveOccurred(), "input: %s", input)
		Expect(parseErr).To(errMatcher, "input: %s", input)
		Expect(node).To(BeNil(), "input: %s", input)
	}
}

func ExpectBlockToHaveEvalError(p parsley.Parser, registry parsley.NodeTransformerRegistry) func(string, types.GomegaMatcher) {
	return func(input string, errMatcher types.GomegaMatcher) {
		parseCtx := ParseCtx(input, registry, nil)
		node, parseErr := parsley.Parse(parseCtx, combinator.Sentence(p))
		Expect(parseErr).ToNot(HaveOccurred(), "input: %s", input)

		_, evalErr := parsley.EvaluateNode(EvalUserCtx(), node)
		Expect(evalErr).To(HaveOccurred())
		Expect(parseCtx.FileSet().ErrorWithPosition(evalErr)).To(errMatcher, "input: %s", input)
	}
}

func ExpectFunctionToEvaluate(p parsley.Parser, registry parsley.NodeTransformerRegistry) func(string, interface{}) {
	return func(input string, expected interface{}) {
		node, parseErr := parsley.Parse(ParseCtx(input, nil, registry), combinator.Sentence(p))
		Expect(parseErr).ToNot(HaveOccurred(), "input: %s", input)

		value, evalErr := parsley.EvaluateNode(EvalUserCtx(), node)
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
		res, err := parsley.Parse(ParseCtx(input, nil, registry), combinator.Sentence(p))
		Expect(err).To(HaveOccurred(), "input: %s", input)
		Expect(err).To(MatchError(expectedErr), "input: %s", input)
		Expect(res).To(BeNil(), "input: %s", input)
	}
}

func ExpectFunctionToHaveEvalError(p parsley.Parser, registry parsley.NodeTransformerRegistry) func(string, error) {
	return func(input string, expectedErr error) {
		parseCtx := ParseCtx(input, nil, registry)
		node, parseErr := parsley.Parse(parseCtx, combinator.Sentence(p))
		Expect(parseErr).ToNot(HaveOccurred(), "input: %s", input)

		value, evalErr := parsley.EvaluateNode(EvalUserCtx(), node)
		Expect(evalErr).To(HaveOccurred(), "input: %s", input)
		Expect(parseCtx.FileSet().ErrorWithPosition(evalErr)).To(MatchError(expectedErr), "input: %s", input)
		Expect(value).To(BeNil(), "input: %s", input)
	}
}

func ExpectFunctionNode(p parsley.Parser, registry parsley.NodeTransformerRegistry) func(string, func(interface{}, parsley.Node)) {
	return func(input string, test func(interface{}, parsley.Node)) {
		ctx := ParseCtx(input, nil, registry)
		node, err := parsley.Parse(ctx, combinator.Sentence(p))
		Expect(err).ToNot(HaveOccurred(), "input: %s", input)

		test(ctx.UserContext(), node)
	}
}
