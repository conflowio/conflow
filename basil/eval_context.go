// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync/atomic"
	"time"
)

// EvalStage means an evaluation stage (default, pre or post)
type EvalStage int8

func (e EvalStage) String() string {
	switch e {
	case EvalStageUndefined:
		return "undefined"
	case EvalStageParse:
		return "parse"
	case EvalStageResolve:
		return "resolve"
	case EvalStageInit:
		return "init"
	case EvalStageMain:
		return "main"
	case EvalStageClose:
		return "close"
	case EvalStageIgnore:
		return "ignore"
	}
	panic(fmt.Errorf("unknown eval stage: %d", e))
}

// Evaluation stages
const (
	EvalStageUndefined EvalStage = iota
	EvalStageParse
	EvalStageResolve
	EvalStageInit
	EvalStageMain
	EvalStageClose
	EvalStageIgnore
)

// EvalStages returns with the evaluation stages
var EvalStages = map[string]EvalStage{
	"parse":   EvalStageParse,
	"resolve": EvalStageResolve,
	"init":    EvalStageInit,
	"main":    EvalStageMain,
	"close":   EvalStageClose,
	"ignore":  EvalStageIgnore,
}

// EvalContext is the evaluation context
type EvalContext struct {
	ctx          context.Context
	cancel       context.CancelFunc
	UserContext  interface{}
	Logger       Logger
	jobScheduler JobScheduler
	parentCtx    *EvalContext
	pubSub       *PubSub
	dependencies map[ID]BlockContainer
	sem          int64
	StdOut       io.Writer
	InputParams  map[ID]interface{}
}

// NewEvalContext returns with a new evaluation context
func NewEvalContext(
	ctx context.Context,
	userContext interface{},
	logger Logger,
	jobScheduler JobScheduler,
	dependencies map[ID]BlockContainer,
) *EvalContext {
	ctx, cancel := context.WithCancel(ctx)
	return &EvalContext{
		ctx:          ctx,
		cancel:       cancel,
		UserContext:  userContext,
		Logger:       logger,
		jobScheduler: jobScheduler,
		pubSub:       NewPubSub(),
		dependencies: dependencies,
		StdOut:       os.Stdout,
	}
}

// New creates a new eval context by copying the parent and overriding the provided values
func (e *EvalContext) New(
	ctx context.Context,
	cancel context.CancelFunc,
	dependencies map[ID]BlockContainer,
) *EvalContext {
	return &EvalContext{
		ctx:          ctx,
		cancel:       cancel,
		dependencies: dependencies,
		UserContext:  e.UserContext,
		Logger:       e.Logger,
		jobScheduler: e.jobScheduler,
		pubSub:       e.pubSub,
		parentCtx:    e,
		StdOut:       e.StdOut,
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

func (e *EvalContext) Subscribe(container *NodeContainer, id ID) {
	e.pubSub.Subscribe(container, id)
}

func (e *EvalContext) Unsubscribe(container *NodeContainer, id ID) {
	e.pubSub.Unsubscribe(container, id)
}

func (e *EvalContext) Publish(c Container) {
	e.pubSub.Publish(c)
}

func (e *EvalContext) HasSubscribers(id ID) bool {
	return e.pubSub.HasSubscribers(id)
}

func (e *EvalContext) Deadline() (deadline time.Time, ok bool) {
	return e.ctx.Deadline()
}

func (e *EvalContext) Done() <-chan struct{} {
	return e.ctx.Done()
}

func (e *EvalContext) Err() error {
	return e.ctx.Err()
}

func (e *EvalContext) Value(key interface{}) interface{} {
	return e.ctx.Value(key)
}

func (e *EvalContext) Run() bool {
	return atomic.CompareAndSwapInt64(&e.sem, 0, 1)
}

func (e *EvalContext) Cancel() bool {
	swapped := atomic.CompareAndSwapInt64(&e.sem, 0, 2)
	e.cancel()
	return swapped
}

func (e *EvalContext) JobScheduler() JobScheduler {
	return e.jobScheduler
}

func (e *EvalContext) SetStdOut(stdOut io.Writer) {
	e.StdOut = stdOut
}
