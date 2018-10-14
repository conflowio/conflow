package test

import (
	"fmt"

	. "github.com/onsi/gomega"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
)

func ExpectParserToEvaluate(p parsley.Parser) func(string, interface{}) {
	return func(input string, expected interface{}) {
		f := text.NewFile("testfile", []byte(input))
		fs := parsley.NewFileSet(f)
		r := text.NewReader(f)
		ctx := parsley.NewContext(fs, r)
		val, err := parsley.Evaluate(ctx, combinator.Sentence(p), nil)

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
		ctx := parsley.NewContext(fs, r)
		val, err := parsley.Evaluate(ctx, combinator.Sentence(p), nil)

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
