// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil_test

import (
	"context"
	"errors"
	"time"

	"github.com/opsidian/parsley/parsley"

	"github.com/opsidian/basil/loggers/zerolog"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/basilfakes"
)

var _ = Describe("NodeContainer", func() {
	var n *basil.NodeContainer
	var nerr parsley.Error
	var evalCtx *basil.EvalContext
	var node, parentNode *basilfakes.FakeNode
	var parentContainer *basilfakes.FakeBlockContainer
	var ctx context.Context
	var cancel context.CancelFunc
	var scheduler *basilfakes.FakeJobScheduler

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		parentNode = &basilfakes.FakeNode{}
		parentNode.IDReturns("parent_node_id")

		node = &basilfakes.FakeNode{}
		node.IDReturns("node_id")

		parentContainer = &basilfakes.FakeBlockContainer{}
		parentContainer.NodeReturns(parentNode)

		logger := zerolog.NewDisabledLogger()
		scheduler = &basilfakes.FakeJobScheduler{}
		evalCtx = basil.NewEvalContext(ctx, nil, logger, scheduler, nil)
	})

	AfterEach(func() {
		cancel()
	})

	JustBeforeEach(func() {
		n, nerr = basil.NewNodeContainer(evalCtx, parentContainer, node, scheduler)
	})

	It("should have no error", func() {
		Expect(nerr).ToNot(HaveOccurred())
	})

	When("a resolve directive has an error", func() {
		var err parsley.Error

		BeforeEach(func() {
			err = parsley.NewError(0, errors.New("directive error"))
			directiveBlock := &basilfakes.FakeBlockNode{}
			directiveBlock.ValueReturns(nil, err)
			directiveBlock.EvalStageReturns(basil.EvalStageResolve)
			node.DirectivesReturns([]basil.BlockNode{directiveBlock})
		})

		It("should return the error", func() {
			Expect(nerr).To(MatchError(err))
		})
	})

	Context("Run", func() {
		var isPending bool
		var runErr parsley.Error

		JustBeforeEach(func() {
			isPending, runErr = n.Run()
		})

		When("the node is ready to run", func() {
			BeforeEach(func() {
				node.CreateContainerReturns(&basilfakes.FakeJobContainer{})
			})

			It("will create a container", func() {
				Expect(node.CreateContainerCallCount()).To(Equal(1))
				_, parent, _, wgs, pending := node.CreateContainerArgsForCall(0)
				Expect(parent).To(Equal(parentContainer))
				Expect(wgs).To(BeEmpty())
				Expect(pending).To(BeFalse())
			})

			It("calls schedule", func() {
				Expect(scheduler.ScheduleJobCallCount()).To(Equal(1))
				scheduledContainer := scheduler.ScheduleJobArgsForCall(0)
				Expect(scheduledContainer).To(Equal(scheduledContainer))
			})

			When("the scheduler can run the container", func() {
				BeforeEach(func() {
					scheduler.ScheduleJobReturns(nil)
				})

				It("should schedule a new container", func() {
					Expect(runErr).ToNot(HaveOccurred())
					Expect(isPending).To(BeFalse())
				})
			})

			When("the scheduler had an error", func() {
				var err error
				BeforeEach(func() {
					err = errors.New("some error")
					scheduler.ScheduleJobReturns(err)
				})

				It("should return the error", func() {
					Expect(runErr).To(MatchError(parsley.NewError(0, err)))
				})
			})
		})

		When("the node dependency is already met in the parent", func() {
			BeforeEach(func() {
				dep1 := &basilfakes.FakeVariableNode{}
				dep1.IDReturns("other_node.dep1")
				dep1.ParentIDReturns("other_node")
				node.DependenciesReturns(basil.Dependencies{
					"other_node.dep1": dep1,
				})
				parentNode.DependenciesReturns(basil.Dependencies{
					"other_node.dep1": dep1,
				})

				node.CreateContainerReturns(&basilfakes.FakeJobContainer{})
				scheduler.ScheduleJobReturns(nil)
			})

			It("will schedule the container", func() {
				Expect(node.CreateContainerCallCount()).To(Equal(1))
				Expect(scheduler.ScheduleJobCallCount()).To(Equal(1))
				Expect(runErr).ToNot(HaveOccurred())
				Expect(isPending).To(BeFalse())
			})
		})

		When("the node has unmet dependencies", func() {
			BeforeEach(func() {
				dep1 := &basilfakes.FakeVariableNode{}
				dep1.IDReturns("other_node.dep1")
				dep1.ParentIDReturns("other_node")
				node.DependenciesReturns(basil.Dependencies{
					"other_node.dep1": dep1,
				})
			})

			It("should not run", func() {
				Expect(runErr).ToNot(HaveOccurred())
				Expect(isPending).To(BeTrue())
				Expect(node.CreateContainerCallCount()).To(Equal(0))
				Expect(scheduler.ScheduleJobCallCount()).To(Equal(0))
			})

			When("the missing dependency is met", func() {
				JustBeforeEach(func() {
					depNode := &basilfakes.FakeNode{}
					depNode.IDReturns("other_node")
					dep := &basilfakes.FakeBlockContainer{}
					dep.NodeReturns(depNode)
					dep.ValueReturns("foo", nil)
					evalCtx.Publish(dep)
				})

				It("should create a container with a pending status", func() {
					Expect(node.CreateContainerCallCount()).To(Equal(1))
					_, _, _, _, pending := node.CreateContainerArgsForCall(0)
					Expect(pending).To(BeTrue())
				})
			})
		})

		When("the node is skipped", func() {
			BeforeEach(func() {
				directive := &basilfakes.FakeDirective{}
				directive.RuntimeConfigReturns(basil.RuntimeConfig{Skip: true})
				directiveBlock := &basilfakes.FakeBlockNode{}
				directiveBlock.ValueReturns(directive, nil)
				directiveBlock.EvalStageReturns(basil.EvalStageInit)
				node.DirectivesReturns([]basil.BlockNode{directiveBlock})
			})

			It("should return not pending", func() {
				Expect(runErr).ToNot(HaveOccurred())
				Expect(isPending).To(BeFalse())
			})

			It("should not be scheduled", func() {
				Expect(node.CreateContainerCallCount()).To(Equal(0))
				Expect(scheduler.ScheduleJobCallCount()).To(Equal(0))
			})
		})

		When("the node has timeout", func() {
			BeforeEach(func() {
				directive := &basilfakes.FakeDirective{}
				directive.RuntimeConfigReturns(basil.RuntimeConfig{Timeout: 1 * time.Second})
				directiveBlock := &basilfakes.FakeBlockNode{}
				directiveBlock.ValueReturns(directive, nil)
				directiveBlock.EvalStageReturns(basil.EvalStageInit)
				node.DirectivesReturns([]basil.BlockNode{directiveBlock})
			})

			It("should pass a context with a timeout", func() {
				Expect(node.CreateContainerCallCount()).To(Equal(1))
				evalCtx, _, _, _, _ := node.CreateContainerArgsForCall(0)
				deadline, ok := evalCtx.Deadline()
				Expect(ok).To(BeTrue(), "was expecting a context with a deadline")
				Expect(deadline).To(BeTemporally(">", time.Time{}))
			})
		})

		When("an init directive has an error", func() {
			var err parsley.Error

			BeforeEach(func() {
				err = parsley.NewError(0, errors.New("directive error"))
				directiveBlock := &basilfakes.FakeBlockNode{}
				directiveBlock.ValueReturns(nil, err)
				directiveBlock.EvalStageReturns(basil.EvalStageInit)
				node.DirectivesReturns([]basil.BlockNode{directiveBlock})
			})

			It("should return the error", func() {
				Expect(runErr).To(MatchError(err))
			})

			It("should not be scheduled", func() {
				Expect(scheduler.ScheduleJobCallCount()).To(Equal(0))
			})
		})
	})

	Context("SetDependency", func() {
		var dependency basil.Container
		var wg *basilfakes.FakeWaitGroup

		JustBeforeEach(func() {
			// This way we'll test whether we properly subscribe to dependencies
			evalCtx.Publish(dependency)
		})

		When("the node still has unmet dependencies", func() {
			BeforeEach(func() {
				depNode := &basilfakes.FakeNode{}
				depNode.IDReturns("dep1")
				dep := &basilfakes.FakeBlockContainer{}
				dep.NodeReturns(depNode)
				wg = &basilfakes.FakeWaitGroup{}
				dep.WaitGroupsReturns([]basil.WaitGroup{wg})
				dep.ValueReturns("foo", nil)
				dependency = dep

				dep1 := &basilfakes.FakeVariableNode{}
				dep1.IDReturns("dep1.param1")
				dep1.ParentIDReturns("dep1")
				dep2 := &basilfakes.FakeVariableNode{}
				dep2.IDReturns("dep2.param2")
				dep2.ParentIDReturns("dep2")
				node.DependenciesReturns(basil.Dependencies{
					"dep1.param1": dep1,
					"dep2.param2": dep2,
				})
			})

			It("should not run", func() {
				Expect(node.CreateContainerCallCount()).To(Equal(0))
				Expect(scheduler.ScheduleJobCallCount()).To(Equal(0))
			})

			It("should increase any passed wait groups", func() {
				Expect(wg.AddCallCount()).To(Equal(1))
				Expect(wg.AddArgsForCall(0)).To(Equal(1))
			})
		})

		When("the node has all dependencies", func() {
			BeforeEach(func() {
				depNode := &basilfakes.FakeNode{}
				depNode.IDReturns("dep1")
				dep := &basilfakes.FakeBlockContainer{}
				dep.NodeReturns(depNode)
				wg = &basilfakes.FakeWaitGroup{}
				dep.WaitGroupsReturns([]basil.WaitGroup{wg})
				dep.ValueReturns("foo", nil)
				dependency = dep

				dep1 := &basilfakes.FakeVariableNode{}
				dep1.IDReturns("dep1.param1")
				dep1.ParentIDReturns("dep1")
				node.DependenciesReturns(basil.Dependencies{
					"dep1.param1": dep1,
				})

				node.CreateContainerReturns(&basilfakes.FakeJobContainer{})
			})

			It("should run", func() {
				Expect(node.CreateContainerCallCount()).To(Equal(1))
				Expect(scheduler.ScheduleJobCallCount()).To(Equal(1))
			})

			It("should pass the dependencies and the wait groups", func() {
				evalCtx, _, _, wgs, _ := node.CreateContainerArgsForCall(0)
				passedDep, _ := evalCtx.BlockContainer("dep1")
				Expect(passedDep).To(Equal(dependency))
				Expect(wgs).To(ConsistOf(wg))
			})

			When("the parent container doesn't have the same eval stage", func() {
				BeforeEach(func() {
					parentContainer.EvalStageReturns(basil.EvalStageInit)
					node.EvalStageReturns(basil.EvalStageMain)
				})

				It("should not schedule the container", func() {
					Expect(node.CreateContainerCallCount()).To(Equal(0))
					Expect(scheduler.ScheduleJobCallCount()).To(Equal(0))
				})
			})
		})

		When("the node has triggers set", func() {
			var triggers []basil.ID

			BeforeEach(func() {
				directive := &basilfakes.FakeDirective{}
				directive.RuntimeConfigReturns(basil.RuntimeConfig{Triggers: triggers})
				directiveBlock := &basilfakes.FakeBlockNode{}
				directiveBlock.ValueReturns(directive, nil)
				directiveBlock.EvalStageReturns(basil.EvalStageResolve)
				node.DirectivesReturns([]basil.BlockNode{directiveBlock})

				depNode := &basilfakes.FakeNode{}
				depNode.IDReturns("dep1")
				dep := &basilfakes.FakeBlockContainer{}
				dep.NodeReturns(depNode)
				wg = &basilfakes.FakeWaitGroup{}
				dep.WaitGroupsReturns([]basil.WaitGroup{wg})
				dep.ValueReturns("foo", nil)
				dependency = dep

				dep1 := &basilfakes.FakeVariableNode{}
				dep1.IDReturns("dep1.param1")
				dep1.ParentIDReturns("dep1")
				node.DependenciesReturns(basil.Dependencies{
					"dep1.param1": dep1,
				})

				node.CreateContainerReturns(&basilfakes.FakeJobContainer{})
			})

			When("the dependency is not a trigger", func() {
				BeforeEach(func() {
					triggers = []basil.ID{"dep2"}
				})

				It("should run the first time", func() {
					Expect(node.CreateContainerCallCount()).To(Equal(1))
					Expect(scheduler.ScheduleJobCallCount()).To(Equal(1))
				})

				It("should not add the wait groups", func() {
					_, _, _, wgs, _ := node.CreateContainerArgsForCall(0)
					Expect(wgs).To(BeEmpty())
				})

				When("setting the dependency the second time", func() {
					JustBeforeEach(func() {
						evalCtx.Publish(dependency)
					})

					It("should not run the second time", func() {
						Expect(node.CreateContainerCallCount()).To(Equal(1))
						Expect(scheduler.ScheduleJobCallCount()).To(Equal(1))
					})
				})
			})

			When("the dependency is a trigger", func() {
				BeforeEach(func() {
					triggers = []basil.ID{"dep1"}
				})

				It("should run", func() {
					Expect(node.CreateContainerCallCount()).To(Equal(1))
					Expect(scheduler.ScheduleJobCallCount()).To(Equal(1))
				})

				It("should add the wait groups", func() {
					_, _, _, wgs, _ := node.CreateContainerArgsForCall(0)
					Expect(wgs).To(ConsistOf(wg))
				})
			})
		})

		When("the dependency is a sibling parameter", func() {
			BeforeEach(func() {
				depNode := &basilfakes.FakeNode{}
				depNode.IDReturns("parent_node_id.sibling")
				dep := &basilfakes.FakeParameterContainer{}
				dep.NodeReturns(depNode)
				dep.BlockContainerReturns(parentContainer)
				dep.ValueReturns("foo", nil)
				dependency = dep

				dep1 := &basilfakes.FakeVariableNode{}
				dep1.IDReturns("parent_node_id.sibling")
				dep1.ParentIDReturns("parent_node_id")
				node.DependenciesReturns(basil.Dependencies{
					"parent_node_id.sibling": dep1,
				})

				node.CreateContainerReturns(&basilfakes.FakeJobContainer{})
			})

			It("should run", func() {
				Expect(node.CreateContainerCallCount()).To(Equal(1))
				evalCtx, _, _, _, _ := node.CreateContainerArgsForCall(0)
				passedDep, _ := evalCtx.BlockContainer("parent_node_id")
				Expect(passedDep).To(Equal(parentContainer))
				Expect(scheduler.ScheduleJobCallCount()).To(Equal(1))
			})
		})

	})
})
