// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

import (
	"context"
	"time"
)

// BlockContext is passed to the block objects during evaluation
//go:generate counterfeiter . BlockContext
type BlockContext interface {
	context.Context
	UserContext() interface{}
	Logger() Logger
	PublishBlock(Block, func() error) (bool, error)
}

// Context defines an interface about creating a new context
type Contexter interface {
	Context(ctx context.Context) (context.Context, context.CancelFunc)
}

// UserContexter defines an interface about creating a new user context
type UserContexter interface {
	UserContext(ctx interface{}) interface{}
}

// UserContexter defines an interface about creating a new logger
type Loggerer interface {
	Logger(logger Logger) Logger
}

// NewBlockContext creates a new block context
func NewBlockContext(
	evalContext *EvalContext,
	container BlockContainer,
) BlockContext {
	return &blockContext{
		evalContext: evalContext,
		container:   container,
	}
}

type blockContext struct {
	evalContext *EvalContext
	container   BlockContainer
}

func (b *blockContext) Deadline() (deadline time.Time, ok bool) {
	return b.evalContext.Deadline()
}

func (b *blockContext) Done() <-chan struct{} {
	return b.evalContext.Done()
}

func (b *blockContext) Err() error {
	return b.evalContext.Err()
}

func (b *blockContext) Value(key interface{}) interface{} {
	return b.evalContext.Value(key)
}

func (b *blockContext) UserContext() interface{} {
	return b.evalContext.UserContext
}

func (b *blockContext) Logger() Logger {
	return b.evalContext.Logger
}

func (b *blockContext) PublishBlock(block Block, f func() error) (bool, error) {
	return b.container.PublishBlock(block, f)
}
