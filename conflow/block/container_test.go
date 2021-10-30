// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package block_test

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/parsley/parsley"

	"github.com/opsidian/conflow/conflow"
	"github.com/opsidian/conflow/conflow/block"
	"github.com/opsidian/conflow/conflow/conflowfakes"
	"github.com/opsidian/conflow/conflow/job"
	"github.com/opsidian/conflow/loggers/zerolog"
	"github.com/opsidian/conflow/test/testfakes"
)

var _ = Describe("Container", func() {
	var container *block.Container
	var evalCtx *conflow.EvalContext
	var blockNode *conflowfakes.FakeBlockNode
	var b conflow.Block
	var ctx context.Context
	var cancel context.CancelFunc
	var scheduler *job.Scheduler
	var interpreter *conflowfakes.FakeBlockInterpreter
	var value interface{}
	var err parsley.Error

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		interpreter = &conflowfakes.FakeBlockInterpreter{}
		blockNode = &conflowfakes.FakeBlockNode{}
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
		evalCtx = conflow.NewEvalContext(ctx, nil, logger, scheduler, nil)

		container = block.NewContainer(evalCtx, conflow.RuntimeConfig{}, blockNode, nil, nil, nil, false)
		container.Run()
		value, err = container.Value()
	})

	Context("when a node has no children", func() {
		var fakeBlock *conflowfakes.FakeIdentifiable

		BeforeEach(func() {
			fakeBlock = &conflowfakes.FakeIdentifiable{}
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
				fakeBlock.InitReturns(false, errors.New("init error"))
				b = fakeBlock
			})

			It("should return with the error", func() {
				Expect(err).To(MatchError("init error"))
			})
		})

	})

	Context("when a node has a run method", func() {
		var fakeBlock *testfakes.FakeBlockWithRun

		When("it has no error", func() {
			BeforeEach(func() {
				fakeBlock = &testfakes.FakeBlockWithRun{}
				fakeBlock.RunReturns(nil, nil)
				b = fakeBlock
			})

			It("should return with no error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should call the main method on the block", func() {
				Expect(fakeBlock.RunCallCount()).To(Equal(1))
			})

			It("should return with the created block", func() {
				Expect(value).To(Equal(fakeBlock))
			})
		})

		When("it has an error", func() {
			BeforeEach(func() {
				fakeBlock = &testfakes.FakeBlockWithRun{}
				fakeBlock.RunReturns(nil, errors.New("main error"))
				b = fakeBlock
			})

			It("should return with the error", func() {
				Expect(err).To(MatchError("main error"))
			})
		})

		When("it returns with a retry result on the first run", func() {
			BeforeEach(func() {
				fakeBlock = &testfakes.FakeBlockWithRun{}
				fakeBlock.RunReturnsOnCall(0, conflow.Retry("test retry"), nil)
				fakeBlock.RunReturnsOnCall(1, nil, nil)
				b = fakeBlock
			})

			It("should return with no error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should call the main method on the block twice", func() {
				Expect(fakeBlock.RunCallCount()).To(Equal(2))
			})

			It("should return with the created block", func() {
				Expect(value).To(Equal(fakeBlock))
			})
		})

		When("it returns with a retry result (with delay) on the first run", func() {
			BeforeEach(func() {
				fakeBlock = &testfakes.FakeBlockWithRun{}
				fakeBlock.RunReturnsOnCall(0, conflow.RetryAfter(1*time.Millisecond, "test retry"), nil)
				fakeBlock.RunReturnsOnCall(1, nil, nil)
				b = fakeBlock
			})

			It("should return with no error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should call the main method on the block twice", func() {
				Expect(fakeBlock.RunCallCount()).To(Equal(2))
			})

			It("should return with the created block", func() {
				Expect(value).To(Equal(fakeBlock))
			})
		})

		When("it returns with a retryable error on the first run", func() {
			BeforeEach(func() {
				fakeBlock = &testfakes.FakeBlockWithRun{}
				fakeBlock.RunReturnsOnCall(0, nil, conflow.RetryableError(errors.New("test error"), 0))
				fakeBlock.RunReturnsOnCall(1, nil, nil)
				b = fakeBlock
			})

			It("should return with no error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should call the main method on the block twice", func() {
				Expect(fakeBlock.RunCallCount()).To(Equal(2))
			})

			It("should return with the created block", func() {
				Expect(value).To(Equal(fakeBlock))
			})
		})

		When("it returns with a retryable error (with duration) on the first run", func() {
			BeforeEach(func() {
				fakeBlock = &testfakes.FakeBlockWithRun{}
				fakeBlock.RunReturnsOnCall(0, nil, conflow.RetryableError(errors.New("test error"), 1*time.Millisecond))
				fakeBlock.RunReturnsOnCall(1, nil, nil)
				b = fakeBlock
			})

			It("should return with no error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should call the main method on the block twice", func() {
				Expect(fakeBlock.RunCallCount()).To(Equal(2))
			})

			It("should return with the created block", func() {
				Expect(value).To(Equal(fakeBlock))
			})
		})
	})

	Context("when a node has an close method", func() {
		var fakeBlock *testfakes.FakeBlockWithClose

		When("it has no error", func() {
			BeforeEach(func() {
				fakeBlock = &testfakes.FakeBlockWithClose{}
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
				fakeBlock.CloseReturns(errors.New("close error"))
				b = fakeBlock
			})

			It("should return with the error", func() {
				Expect(err).To(MatchError("close error"))
			})
		})
	})
})
