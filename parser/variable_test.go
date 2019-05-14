package parser_test

import (
	"context"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/basilfakes"
	"github.com/opsidian/basil/basil/block"
	"github.com/opsidian/basil/test"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/basil/parser"
)

var _ = Describe("Variable", func() {

	var p = parser.Variable()
	var parsleyContext *parsley.Context
	var parseCtx *basil.ParseContext
	var evalCtx *basil.EvalContext
	var res parsley.Node
	var parseErr, evalErr error
	var value interface{}
	var input string
	var blockNode *basilfakes.FakeBlockNode

	BeforeEach(func() {
		parseCtx = basil.NewParseContext(basil.NewIDRegistry(8, 16))
		evalCtx = basil.NewEvalContext(context.Background(), nil)
		parseErr = nil
		evalErr = nil
		value = nil
		blockNode = nil
	})

	JustBeforeEach(func() {
		parsleyContext = test.ParseCtx(input, nil, nil)
		parsleyContext.SetUserContext(parseCtx)

		if blockNode != nil {
			err := parseCtx.AddBlockNode(blockNode)
			Expect(err).ToNot(HaveOccurred())
		}

		res, parseErr = parsley.Parse(parsleyContext, combinator.Sentence(p))
		if parseErr == nil {
			value, evalErr = res.Value(evalCtx)
		}
	})

	Context("when referencing only the block name", func() {
		BeforeEach(func() {
			input = "foo"
		})

		It("should return with a parse error", func() {
			Expect(parseErr).To(HaveOccurred())
		})
	})

	Context("when referencing a block module parameter", func() {
		var fooBlock *basilfakes.FakeBlock
		var fooBlockInterpreter *basilfakes.FakeBlockInterpreter

		BeforeEach(func() {
			blockNode = &basilfakes.FakeBlockNode{}
			blockNode.IDReturns(basil.ID("foo"))

			fooBlock = &basilfakes.FakeBlock{}
			fooBlockInterpreter = &basilfakes.FakeBlockInterpreter{}
			fooBlockInterpreter.ParamReturnsOnCall(0, "bar")

			blockContainer := block.NewContainer(basil.ID("foo"), fooBlock, fooBlockInterpreter)
			err := evalCtx.AddBlockContainer(blockContainer)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("with an existing parameter", func() {
			BeforeEach(func() {
				input = "foo.param1"
				blockNode.ParamTypeReturnsOnCall(0, "string", true)
			})

			It("should evaluate successfully", func() {
				Expect(parseErr).ToNot(HaveOccurred())
				Expect(evalErr).ToNot(HaveOccurred())
				Expect(value).To(Equal("bar"))

				Expect(blockNode.ParamTypeArgsForCall(0)).To(Equal(basil.ID("param1")))
				passedBlock, passedParam := fooBlockInterpreter.ParamArgsForCall(0)
				Expect(passedBlock).To(Equal(fooBlock))
				Expect(passedParam).To(Equal(basil.ID("param1")))
			})
		})

		Context("with a nonexisting parameter", func() {
			BeforeEach(func() {
				input = "foo.param1"
				blockNode.ParamTypeReturnsOnCall(0, "", false)
			})

			It("should return a parse error", func() {
				Expect(parseErr).To(MatchError("parameter \"param1\" does not exist at testfile:1:5"))
			})
		})
	})

	Context("when referencing a non-existing block", func() {
		BeforeEach(func() {
			input = "foo.param1"
		})

		It("should return a parse error", func() {
			Expect(parseErr).To(MatchError("block \"foo\" does not exist at testfile:1:1"))
		})
	})

})
