// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conflow_test

import (
	"context"
	"errors"
	"time"

	"github.com/opsidian/conflow/conflow"
	"github.com/opsidian/conflow/util"

	"github.com/opsidian/parsley/parsley"

	"github.com/opsidian/conflow/loggers/zerolog"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/opsidian/conflow/conflow/conflowfakes"
)

var _ = Describe("NodeContainer", func() {
	var n *conflow.NodeContainer
	var nerr parsley.Error
	var evalCtx *conflow.EvalContext
	var node, parentNode *conflowfakes.FakeNode
	var parentContainer *conflowfakes.FakeBlockContainer
	var ctx context.Context
	var cancel context.CancelFunc
	var scheduler *conflowfakes.FakeJobScheduler

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		parentNode = &conflowfakes.FakeNode{}
		parentNode.IDReturns("parent_node_id")

		node = &conflowfakes.FakeNode{}
		node.IDReturns("node_id")

		parentContainer = &conflowfakes.FakeBlockContainer{}
		parentContainer.NodeReturns(parentNode)

		logger := zerolog.NewDisabledLogger()
		scheduler = &conflowfakes.FakeJobScheduler{}
		evalCtx = conflow.NewEvalContext(ctx, nil, logger, scheduler, nil)
	})

	AfterEach(func() {
		cancel()
	})

	JustBeforeEach(func() {
		n, nerr = conflow.NewNodeContainer(evalCtx, parentContainer, node, scheduler)
	})

	It("should have no error", func() {
		Expect(nerr).ToNot(HaveOccurred())
	})

	When("a resolve directive has an error", func() {
		var err parsley.Error

		BeforeEach(func() {
			err = parsley.NewError(0, errors.New("directive error"))
			directiveBlock := &conflowfakes.FakeBlockNode{}
			directiveBlock.ValueReturns(nil, err)
			directiveBlock.EvalStageReturns(conflow.EvalStageResolve)
			node.DirectivesReturns([]conflow.BlockNode{directiveBlock})
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
				node.CreateContainerReturns(&conflowfakes.FakeJobContainer{})
			})

			It("will create a container", func() {
				Expect(node.CreateContainerCallCount()).To(Equal(1))
				_, _, parent, _, wgs, pending := node.CreateContainerArgsForCall(0)
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
				dep1 := &conflowfakes.FakeVariableNode{}
				dep1.IDReturns("other_node.dep1")
				dep1.ParentIDReturns("other_node")
				node.DependenciesReturns(conflow.Dependencies{
					"other_node.dep1": dep1,
				})
				parentNode.DependenciesReturns(conflow.Dependencies{
					"other_node.dep1": dep1,
				})

				node.CreateContainerReturns(&conflowfakes.FakeJobContainer{})
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
				dep1 := &conflowfakes.FakeVariableNode{}
				dep1.IDReturns("other_node.dep1")
				dep1.ParentIDReturns("other_node")
				node.DependenciesReturns(conflow.Dependencies{
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
					depNode := &conflowfakes.FakeNode{}
					depNode.IDReturns("other_node")
					dep := &conflowfakes.FakeBlockContainer{}
					dep.NodeReturns(depNode)
					dep.ValueReturns("foo", nil)
					evalCtx.Publish(dep)
				})

				It("should create a container with a pending status", func() {
					Expect(node.CreateContainerCallCount()).To(Equal(1))
					_, _, _, _, _, pending := node.CreateContainerArgsForCall(0)
					Expect(pending).To(BeTrue())
				})
			})
		})

		When("the node is skipped", func() {
			BeforeEach(func() {
				directive := &conflowfakes.FakeBlockDirective{}
				directive.ApplyToRuntimeConfigStub = func(config *conflow.RuntimeConfig) {
					config.Skip = util.BoolPtr(true)
				}
				directiveBlock := &conflowfakes.FakeBlockNode{}
				directiveBlock.ValueReturns(directive, nil)
				directiveBlock.EvalStageReturns(conflow.EvalStageInit)
				node.DirectivesReturns([]conflow.BlockNode{directiveBlock})
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
				directive := &conflowfakes.FakeBlockDirective{}
				directive.ApplyToRuntimeConfigStub = func(config *conflow.RuntimeConfig) {
					config.Timeout = util.TimeDurationPtr(1 * time.Second)
				}
				directiveBlock := &conflowfakes.FakeBlockNode{}
				directiveBlock.ValueReturns(directive, nil)
				directiveBlock.EvalStageReturns(conflow.EvalStageInit)
				node.DirectivesReturns([]conflow.BlockNode{directiveBlock})
			})

			It("should pass a context with a timeout", func() {
				Expect(node.CreateContainerCallCount()).To(Equal(1))
				evalCtx, _, _, _, _, _ := node.CreateContainerArgsForCall(0)
				deadline, ok := evalCtx.Deadline()
				Expect(ok).To(BeTrue(), "was expecting a context with a deadline")
				Expect(deadline).To(BeTemporally(">", time.Time{}))
			})
		})

		When("an init directive has an error", func() {
			var err parsley.Error

			BeforeEach(func() {
				err = parsley.NewError(0, errors.New("directive error"))
				directiveBlock := &conflowfakes.FakeBlockNode{}
				directiveBlock.ValueReturns(nil, err)
				directiveBlock.EvalStageReturns(conflow.EvalStageInit)
				node.DirectivesReturns([]conflow.BlockNode{directiveBlock})
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
		var dependency conflow.Container
		var wg *conflowfakes.FakeWaitGroup

		JustBeforeEach(func() {
			// This way we'll test whether we properly subscribe to dependencies
			evalCtx.Publish(dependency)
		})

		When("the node still has unmet dependencies", func() {
			BeforeEach(func() {
				depNode := &conflowfakes.FakeNode{}
				depNode.IDReturns("dep1")
				dep := &conflowfakes.FakeBlockContainer{}
				dep.NodeReturns(depNode)
				wg = &conflowfakes.FakeWaitGroup{}
				dep.WaitGroupsReturns([]conflow.WaitGroup{wg})
				dep.ValueReturns("foo", nil)
				dependency = dep

				dep1 := &conflowfakes.FakeVariableNode{}
				dep1.IDReturns("dep1.param1")
				dep1.ParentIDReturns("dep1")
				dep2 := &conflowfakes.FakeVariableNode{}
				dep2.IDReturns("dep2.param2")
				dep2.ParentIDReturns("dep2")
				node.DependenciesReturns(conflow.Dependencies{
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
				depNode := &conflowfakes.FakeNode{}
				depNode.IDReturns("dep1")
				dep := &conflowfakes.FakeBlockContainer{}
				dep.NodeReturns(depNode)
				wg = &conflowfakes.FakeWaitGroup{}
				dep.WaitGroupsReturns([]conflow.WaitGroup{wg})
				dep.ValueReturns("foo", nil)
				dependency = dep

				dep1 := &conflowfakes.FakeVariableNode{}
				dep1.IDReturns("dep1.param1")
				dep1.ParentIDReturns("dep1")
				node.DependenciesReturns(conflow.Dependencies{
					"dep1.param1": dep1,
				})

				node.CreateContainerReturns(&conflowfakes.FakeJobContainer{})
			})

			It("should run", func() {
				Expect(node.CreateContainerCallCount()).To(Equal(1))
				Expect(scheduler.ScheduleJobCallCount()).To(Equal(1))
			})

			It("should pass the dependencies and the wait groups", func() {
				evalCtx, _, _, _, wgs, _ := node.CreateContainerArgsForCall(0)
				passedDep, _ := evalCtx.BlockContainer("dep1")
				Expect(passedDep).To(Equal(dependency))
				Expect(wgs).To(ConsistOf(wg))
			})

			When("the parent container doesn't have the same eval stage", func() {
				BeforeEach(func() {
					parentContainer.EvalStageReturns(conflow.EvalStageInit)
					node.EvalStageReturns(conflow.EvalStageMain)
				})

				It("should not schedule the container", func() {
					Expect(node.CreateContainerCallCount()).To(Equal(0))
					Expect(scheduler.ScheduleJobCallCount()).To(Equal(0))
				})
			})
		})

		When("the node has triggers set", func() {
			var triggers []conflow.ID

			BeforeEach(func() {
				directive := &conflowfakes.FakeBlockDirective{}
				directive.ApplyToRuntimeConfigStub = func(config *conflow.RuntimeConfig) {
					config.Triggers = triggers
				}
				directiveBlock := &conflowfakes.FakeBlockNode{}
				directiveBlock.ValueReturns(directive, nil)
				directiveBlock.EvalStageReturns(conflow.EvalStageResolve)
				node.DirectivesReturns([]conflow.BlockNode{directiveBlock})

				depNode := &conflowfakes.FakeNode{}
				depNode.IDReturns("dep1")
				dep := &conflowfakes.FakeBlockContainer{}
				dep.NodeReturns(depNode)
				wg = &conflowfakes.FakeWaitGroup{}
				dep.WaitGroupsReturns([]conflow.WaitGroup{wg})
				dep.ValueReturns("foo", nil)
				dependency = dep

				dep1 := &conflowfakes.FakeVariableNode{}
				dep1.IDReturns("dep1.param1")
				dep1.ParentIDReturns("dep1")
				node.DependenciesReturns(conflow.Dependencies{
					"dep1.param1": dep1,
				})

				node.CreateContainerReturns(&conflowfakes.FakeJobContainer{})
			})

			When("the dependency is not a trigger", func() {
				BeforeEach(func() {
					triggers = []conflow.ID{"dep2"}
				})

				It("should run the first time", func() {
					Expect(node.CreateContainerCallCount()).To(Equal(1))
					Expect(scheduler.ScheduleJobCallCount()).To(Equal(1))
				})

				It("should not add the wait groups", func() {
					_, _, _, _, wgs, _ := node.CreateContainerArgsForCall(0)
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
					triggers = []conflow.ID{"dep1"}
				})

				It("should run", func() {
					Expect(node.CreateContainerCallCount()).To(Equal(1))
					Expect(scheduler.ScheduleJobCallCount()).To(Equal(1))
				})

				It("should add the wait groups", func() {
					_, _, _, _, wgs, _ := node.CreateContainerArgsForCall(0)
					Expect(wgs).To(ConsistOf(wg))
				})
			})
		})

		When("the dependency is a sibling parameter", func() {
			BeforeEach(func() {
				depNode := &conflowfakes.FakeNode{}
				depNode.IDReturns("parent_node_id.sibling")
				dep := &conflowfakes.FakeParameterContainer{}
				dep.NodeReturns(depNode)
				dep.BlockContainerReturns(parentContainer)
				dep.ValueReturns("foo", nil)
				dependency = dep

				dep1 := &conflowfakes.FakeVariableNode{}
				dep1.IDReturns("parent_node_id.sibling")
				dep1.ParentIDReturns("parent_node_id")
				node.DependenciesReturns(conflow.Dependencies{
					"parent_node_id.sibling": dep1,
				})

				node.CreateContainerReturns(&conflowfakes.FakeJobContainer{})
			})

			It("should run", func() {
				Expect(node.CreateContainerCallCount()).To(Equal(1))
				evalCtx, _, _, _, _, _ := node.CreateContainerArgsForCall(0)
				passedDep, _ := evalCtx.BlockContainer("parent_node_id")
				Expect(passedDep).To(Equal(parentContainer))
				Expect(scheduler.ScheduleJobCallCount()).To(Equal(1))
			})
		})

	})
})
