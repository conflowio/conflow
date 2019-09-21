// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

import "context"

// EvalStage means an evaluation stage (default, pre or post)
type EvalStage int8

// Evaluation stages
const (
	EvalStageInit  EvalStage = -1
	EvalStageMain  EvalStage = 0
	EvalStageClose EvalStage = 1
)

// EvalStages returns with the evaluation stages
var EvalStages = map[string]EvalStage{
	"init":  EvalStageInit,
	"main":  EvalStageMain,
	"close": EvalStageClose,
}

// EvalContext is the evaluation context
type EvalContext struct {
	Context      context.Context
	Cancel       context.CancelFunc
	UserContext  interface{}
	Logger       Logger
	parentCtx    *EvalContext
	scheduler    Scheduler
	dependencies map[ID]BlockContainer
	pubsub       *PubSub
}

// NewEvalContext returns with a new evaluation context
func NewEvalContext(
	ctx context.Context,
	userContext interface{},
	logger Logger,
	scheduler Scheduler,
) *EvalContext {
	ctx, cancel := context.WithCancel(ctx)
	return &EvalContext{
		Context:     ctx,
		Cancel:      cancel,
		UserContext: userContext,
		Logger:      logger,
		scheduler:   scheduler,
		pubsub:      NewPubSub(),
	}
}

// WithDependencies returns a copy of parent with the given dependencies
func (e *EvalContext) WithDependencies(dependencies map[ID]BlockContainer) *EvalContext {
	ctx, cancel := context.WithCancel(context.Background())
	return &EvalContext{
		Context:      ctx,
		Cancel:       cancel,
		UserContext:  e.UserContext,
		Logger:       e.Logger,
		parentCtx:    e,
		scheduler:    e.scheduler,
		pubsub:       e.pubsub,
		dependencies: dependencies,
	}
}

// BlockContainer returns with the given block container instance if it exists
func (e *EvalContext) BlockContainer(id ID) (BlockContainer, bool) {
	if container, ok := e.dependencies[id]; ok {
		return container, true
	}

	if e.parentCtx != nil {
		return e.parentCtx.BlockContainer(id)
	}

	return nil, false
}

func (e *EvalContext) ScheduleJob(job Job) {
	e.scheduler.Schedule(job)
}

func (e *EvalContext) Subscribe(container *NodeContainer, id ID) {
	e.pubsub.Subscribe(container, id)
}

func (e *EvalContext) Unsubscribe(container *NodeContainer, id ID) {
	e.pubsub.Unsubscribe(container, id)
}

func (e *EvalContext) Publish(c Container) {
	e.pubsub.Publish(c)
}
