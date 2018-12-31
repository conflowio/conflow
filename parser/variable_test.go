package parser_test

import (
	"context"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/basilfakes"
	"github.com/opsidian/basil/block"
	"github.com/opsidian/basil/block/blockfakes"
	"github.com/opsidian/basil/function"
	"github.com/opsidian/basil/identifier"
	"github.com/opsidian/basil/test"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/basil/parser"
)

var _ = Describe("Variable", func() {

	var p = parser.Variable()
	var parseCtx *parsley.Context
	var evalCtx interface{}
	var res parsley.Node
	var parseErr, evalErr error
	var value interface{}
	var input string
	var blockNodeRegistry *basilfakes.FakeBlockNodeRegistry
	var blockContainerRegistry block.ContainerRegistry

	BeforeEach(func() {
		parseCtx = nil
		evalCtx = nil
		res = nil
		parseErr = nil
		evalErr = nil
		value = nil
		blockNodeRegistry = &basilfakes.FakeBlockNodeRegistry{}
		blockContainerRegistry = block.NewContainerRegistry()
	})

	JustBeforeEach(func() {
		parseCtx = test.ParseCtx(input, nil, nil)
		parseCtx.SetUserContext(basil.NewParseContext(
			block.InterpreterRegistry{},
			function.InterpreterRegistry{},
			identifier.NewRegistry(8, 16),
			blockNodeRegistry,
		))
		evalCtx = basil.NewEvalContext(context.Background(), nil, blockContainerRegistry)
		res, parseErr = parsley.Parse(parseCtx, combinator.Sentence(p))

		if parseErr == nil {
			value, evalErr = res.Value(evalCtx)
		}
	})

	Context("when referencing a root block", func() {
		var rootNode *basilfakes.FakeBlockNode
		var rootBlock *basilfakes.FakeBlock
		var rootBlockInterpreter *blockfakes.FakeInterpreter

		BeforeEach(func() {
			rootNode = &basilfakes.FakeBlockNode{}
			blockNodeRegistry.BlockNodeReturnsOnCall(0, rootNode, true)

			rootBlock = &basilfakes.FakeBlock{}
			rootBlock.IDReturns("root")
			rootBlockInterpreter = &blockfakes.FakeInterpreter{}
			rootBlockInterpreter.ParamReturnsOnCall(0, "bar")

			rootBlockContainer := block.NewContainer(rootBlock, rootBlockInterpreter)
			blockContainerRegistry.AddBlockContainer(rootBlockContainer)

		})

		Context("with an existing parameter", func() {
			BeforeEach(func() {
				input = "param1"
				rootNode.ParamTypeReturnsOnCall(0, "string", true)
			})

			It("should evaluate successfully", func() {
				Expect(parseErr).ToNot(HaveOccurred())
				Expect(evalErr).ToNot(HaveOccurred())
				Expect(value).To(Equal("bar"))

				Expect(blockNodeRegistry.BlockNodeArgsForCall(0)).To(Equal(basil.ID("root")))
				Expect(rootNode.ParamTypeArgsForCall(0)).To(Equal(basil.ID("param1")))
				passedBlock, passedParam := rootBlockInterpreter.ParamArgsForCall(0)
				Expect(passedBlock).To(Equal(rootBlock))
				Expect(passedParam).To(Equal(basil.ID("param1")))
			})
		})

		Context("with a nonexisting parameter", func() {
			BeforeEach(func() {
				input = "param1"
				rootNode.ParamTypeReturnsOnCall(0, "", false)
			})

			It("should return a parse error", func() {
				Expect(parseErr).To(MatchError("parameter \"param1\" does not exist at testfile:1:1"))
			})
		})
	})

	Context("when referencing a block module parameter", func() {
		var blockNode *basilfakes.FakeBlockNode
		var fooBlock *basilfakes.FakeBlock
		var fooBlockInterpreter *blockfakes.FakeInterpreter

		BeforeEach(func() {
			blockNode = &basilfakes.FakeBlockNode{}
			blockNodeRegistry.BlockNodeReturnsOnCall(0, blockNode, true)

			fooBlock = &basilfakes.FakeBlock{}
			fooBlock.IDReturns("foo")
			fooBlockInterpreter = &blockfakes.FakeInterpreter{}
			fooBlockInterpreter.ParamReturnsOnCall(0, "bar")

			blockContainer := block.NewContainer(fooBlock, fooBlockInterpreter)
			blockContainerRegistry.AddBlockContainer(blockContainer)
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

				Expect(blockNodeRegistry.BlockNodeArgsForCall(0)).To(Equal(basil.ID("foo")))
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
			blockNodeRegistry.BlockNodeReturnsOnCall(0, nil, false)
			input = "foo.param1"
		})

		It("should return a parse error", func() {
			Expect(parseErr).To(MatchError("block \"foo\" does not exist at testfile:1:1"))
		})
	})

})
