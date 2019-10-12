// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package test

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"

	"github.com/opsidian/basil/logger"

	"github.com/opsidian/basil/basil/basilfakes"

	"github.com/onsi/gomega/types"

	"github.com/opsidian/basil/basil/function"

	. "github.com/onsi/gomega"
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
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

	testBlockNode := &basilfakes.FakeBlockNode{}
	testBlockNode.IDReturns(basil.ID("test"))
	testBlockNode.ParamTypeCalls(func(id basil.ID) (s string, b bool) {
		interpreter := TestBlockInterpreter{}
		if param, ok := interpreter.Params()[id]; ok {
			return param.Type, true
		}
		return "", false
	})

	f := text.NewFile("testfile", []byte(input))
	fs := parsley.NewFileSet(f)
	r := text.NewReader(f)

	parseCtx := basil.NewParseContext(newIDRegistry()).New(basil.ParseContextOverride{
		BlockTransformerRegistry:    blockRegistry,
		FunctionTransformerRegistry: functionRegistry,
	})
	_ = parseCtx.AddBlockNode(testBlockNode)

	ctx := parsley.NewContext(fs, r)
	ctx.EnableStaticCheck()
	ctx.EnableTransformation()
	ctx.RegisterKeywords("true", "false", "nil", "map", "testkeyword")
	ctx.SetUserContext(parseCtx)

	return ctx
}

func EvalUserCtx() *basil.EvalContext {
	testBlock :=
		&TestBlock{
			FieldString: "bar",
			FieldMap: map[string]interface{}{
				"key1": "value1",
				"key2": map[string]interface{}{
					"key3": "value3",
				},
			},
			FieldArray: []interface{}{
				"value1",
				"value2",
				[]interface{}{"value2"},
			},
			FieldInt: int64(1),
		}

	testBlockContainer := &basilfakes.FakeBlockContainer{}
	testBlockContainer.ParamCalls(func(name basil.ID) interface{} {
		return TestBlockInterpreter{}.Param(testBlock, name)
	})

	containers := map[basil.ID]basil.BlockContainer{
		"test": testBlockContainer,
	}

	evalCtx := basil.NewEvalContext(
		context.Background(),
		"userCtx",
		logger.NewZeroLogLogger(zerolog.New(os.Stderr).Level(zerolog.Disabled)),
		Scheduler{},
	).New(containers)

	return evalCtx
}

func ExpectParserToEvaluate(p parsley.Parser) func(string, interface{}) {
	return func(input string, expected interface{}) {
		node, parseErr := parsley.Parse(ParseCtx(input, nil, nil), combinator.Sentence(p))
		Expect(parseErr).ToNot(HaveOccurred(), "input: %s", input)

		value, evalErr := node.Value(EvalUserCtx())
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
		Expect(err).To(MatchError(fmt.Errorf("failed to parse the input: %s", expectedErr)), "input: %s", input)
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

		value, evalErr := node.Value(EvalUserCtx())
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

		value, evalErr := node.Value(EvalUserCtx())
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

func ExpectBlockToHaveEvalError(p parsley.Parser, registry parsley.NodeTransformerRegistry) func(string, types.GomegaMatcher) {
	return func(input string, errMatcher types.GomegaMatcher) {
		parseCtx := ParseCtx(input, registry, nil)
		node, parseErr := parsley.Parse(parseCtx, combinator.Sentence(p))
		Expect(parseErr).ToNot(HaveOccurred(), "input: %s", input)

		_, evalErr := node.Value(EvalUserCtx())
		Expect(evalErr).To(HaveOccurred())
		Expect(parseCtx.FileSet().ErrorWithPosition(evalErr)).To(errMatcher, "input: %s", input)
	}
}

func ExpectFunctionToEvaluate(p parsley.Parser, registry parsley.NodeTransformerRegistry) func(string, interface{}) {
	return func(input string, expected interface{}) {
		node, parseErr := parsley.Parse(ParseCtx(input, nil, registry), combinator.Sentence(p))
		Expect(parseErr).ToNot(HaveOccurred(), "input: %s", input)

		value, evalErr := node.Value(EvalUserCtx())
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

		value, evalErr := node.Value(EvalUserCtx())
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
