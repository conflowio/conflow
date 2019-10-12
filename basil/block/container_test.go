// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package block_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/basilfakes"
	"github.com/opsidian/basil/basil/block"
	"github.com/opsidian/basil/basil/job"
	"github.com/opsidian/basil/logger/zerolog"
	"github.com/opsidian/basil/test/testfakes"
	"github.com/opsidian/parsley/parsley"
)

var _ = Describe("Container", func() {
	var container *block.Container
	var evalCtx *basil.EvalContext
	var blockNode *basilfakes.FakeBlockNode
	var b basil.Block
	var ctx context.Context
	var cancel context.CancelFunc
	var scheduler basil.Scheduler
	var interpreter *basilfakes.FakeBlockInterpreter
	var value interface{}
	var err parsley.Error

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		interpreter = &basilfakes.FakeBlockInterpreter{}
		blockNode = &basilfakes.FakeBlockNode{}
		blockNode.IDReturns("test_id")
		blockNode.BlockTypeReturns("test_type")
		blockNode.InterpreterReturns(interpreter)
	})

	AfterEach(func() {
		cancel()
		scheduler.Stop()
	})

	JustBeforeEach(func() {
		interpreter.CreateBlockReturns(b)

		logger := zerolog.NewDisabledLogger()
		scheduler = job.NewScheduler(logger, 1, 10)
		scheduler.Start()
		evalCtx = basil.NewEvalContext(ctx, nil, logger, scheduler)

		container = block.NewContainer(evalCtx, blockNode, nil, nil, nil)
		container.Run()
		value, err = container.Value()
	})

	Context("when a node has no children", func() {
		var fakeBlock *basilfakes.FakeBlock

		BeforeEach(func() {
			fakeBlock = &basilfakes.FakeBlock{}
			fakeBlock.IDReturns("test_id")
			b = fakeBlock
		})

		It("should have no error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return with the created block", func() {
			Expect(value).To(Equal(fakeBlock))
		})
	})

	Context("when a node has an init method", func() {
		var fakeBlock *testfakes.FakeBlockWithInit

		When("it has no error", func() {
			BeforeEach(func() {
				fakeBlock = &testfakes.FakeBlockWithInit{}
				fakeBlock.IDReturns("test_id")
				fakeBlock.InitReturns(false, nil)
				b = fakeBlock
			})

			It("should return with no error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should call the init method on the block", func() {
				Expect(fakeBlock.InitCallCount()).To(Equal(1))
			})

			It("should return with the created block", func() {
				Expect(value).To(Equal(fakeBlock))
			})
		})

		When("the block should be skipped", func() {
			BeforeEach(func() {
				fakeBlock = &testfakes.FakeBlockWithInit{}
				fakeBlock.IDReturns("test_id")
				fakeBlock.InitReturns(true, nil)
				b = fakeBlock
			})

			It("should return with no error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should return with a nil value", func() {
				Expect(value).To(BeNil())
			})
		})

		When("it has an error", func() {
			BeforeEach(func() {
				fakeBlock = &testfakes.FakeBlockWithInit{}
				fakeBlock.IDReturns("test_id")
				fakeBlock.InitReturns(false, errors.New("init error"))
				b = fakeBlock
			})

			It("should return with the error", func() {
				Expect(err).To(MatchError("init error"))
			})
		})

	})

	Context("when a node has an main method", func() {
		var fakeBlock *testfakes.FakeBlockWithMain

		When("it has no error", func() {
			BeforeEach(func() {
				fakeBlock = &testfakes.FakeBlockWithMain{}
				fakeBlock.IDReturns("test_id")
				fakeBlock.MainReturns(nil)
				b = fakeBlock
			})

			It("should return with no error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should call the main method on the block", func() {
				Expect(fakeBlock.MainCallCount()).To(Equal(1))
			})

			It("should return with the created block", func() {
				Expect(value).To(Equal(fakeBlock))
			})
		})

		When("it has an error", func() {
			BeforeEach(func() {
				fakeBlock = &testfakes.FakeBlockWithMain{}
				fakeBlock.IDReturns("test_id")
				fakeBlock.MainReturns(errors.New("main error"))
				b = fakeBlock
			})

			It("should return with the error", func() {
				Expect(err).To(MatchError("main error"))
			})
		})
	})

	Context("when a node has an close method", func() {
		var fakeBlock *testfakes.FakeBlockWithClose

		When("it has no error", func() {
			BeforeEach(func() {
				fakeBlock = &testfakes.FakeBlockWithClose{}
				fakeBlock.IDReturns("test_id")
				fakeBlock.CloseReturns(nil)
				b = fakeBlock
			})

			It("should return with no error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should call the close method on the block", func() {
				Expect(fakeBlock.CloseCallCount()).To(Equal(1))
			})

			It("should return with the created block", func() {
				Expect(value).To(Equal(fakeBlock))
			})
		})

		When("it has an error", func() {
			BeforeEach(func() {
				fakeBlock = &testfakes.FakeBlockWithClose{}
				fakeBlock.IDReturns("test_id")
				fakeBlock.CloseReturns(errors.New("close error"))
				b = fakeBlock
			})

			It("should return with the error", func() {
				Expect(err).To(MatchError("close error"))
			})
		})
	})
})
