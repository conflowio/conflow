package test

import (
	"fmt"

	. "github.com/onsi/gomega"
	"github.com/opsidian/ocl/ocl"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
)

var vp = VariableProvider{map[string]interface{}{
	"foo": "bar",
	"map": map[string]interface{}{
		"key1": "value1",
		"key2": map[string]interface{}{
			"key3": "value3",
		},
		"key4": []interface{}{
			"value4",
		},
	},
	"arr": []interface{}{
		"value1",
		[]interface{}{
			"value2",
		},
		map[string]interface{}{
			"key1": "value3",
		},
	},
	"intkey": int64(1),
}}

func ExpectParserToEvaluate(p parsley.Parser) func(string, interface{}) {
	return func(input string, expected interface{}) {
		f := text.NewFile("testfile", []byte(input))
		fs := parsley.NewFileSet(f)
		r := text.NewReader(f)
		parseCtx := parsley.NewContext(fs, r)
		evalCtx := ocl.NewContext(vp, &FunctionRegistry{}, &BlockRegistry{})
		val, err := parsley.Evaluate(parseCtx, combinator.Sentence(p), evalCtx)

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
		f := text.NewFile("testfile", []byte(input))
		fs := parsley.NewFileSet(f)
		r := text.NewReader(f)
		ctx := parsley.NewContext(fs, r)
		res, err := parsley.Parse(ctx, combinator.Sentence(p))

		Expect(err).To(HaveOccurred(), "input: %s", input)
		Expect(err).To(MatchError(fmt.Errorf("failed to parse the input: %s", expectedErr)), "input: %s", input)
		Expect(res).To(BeNil(), "input: %s", input)
	}
}

func ExpectParserToHaveEvalError(p parsley.Parser) func(string, error) {
	return func(input string, expectedErr error) {
		f := text.NewFile("testfile", []byte(input))
		fs := parsley.NewFileSet(f)
		r := text.NewReader(f)
		parseCtx := parsley.NewContext(fs, r)
		evalCtx := ocl.NewContext(vp, &FunctionRegistry{}, &BlockRegistry{})
		val, err := parsley.Evaluate(parseCtx, combinator.Sentence(p), evalCtx)

		Expect(err).To(HaveOccurred(), "input: %s", input)
		Expect(err).To(MatchError(expectedErr), "input: %s", input)
		Expect(val).To(BeNil(), "input: %s", input)
	}
}

func ExpectParserToReturn(p parsley.Parser, input string, expected parsley.Node) {
	f := text.NewFile("testfile", []byte(input))
	fs := parsley.NewFileSet(f)
	r := text.NewReader(f)
	ctx := parsley.NewContext(fs, r)
	res, err := parsley.Parse(ctx, combinator.Sentence(p))

	Expect(err).ToNot(HaveOccurred())

	node, ok := res.(*ast.NonTerminalNode)
	Expect(ok).To(BeTrue())
	actual := node.Children()[0]

	Expect(actual).To(BeEquivalentTo(expected))
}
