// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conflow

import (
	"context"
	"io"
)

// Context defines an interface about creating a new context
type Contexter interface {
	Context(ctx context.Context) (context.Context, context.CancelFunc)
}

// UserContexter defines an interface about creating a new user context
type UserContexter interface {
	UserContext(ctx interface{}) interface{}
}

// Loggerer defines an interface about creating a new logger
type Loggerer interface {
	Logger(logger Logger) Logger
}

// NewBlockContext creates a new block context
func NewBlockContext(
	evalContext *EvalContext,
	blockPublisher BlockPublisher,
) *BlockContext {
	return &BlockContext{
		evalContext:    evalContext,
		blockPublisher: blockPublisher,
	}
}

type BlockContext struct {
	evalContext    *EvalContext
	blockPublisher BlockPublisher
}

func (b *BlockContext) BlockPublisher() BlockPublisher {
	return b.blockPublisher
}

func (b *BlockContext) JobScheduler() JobScheduler {
	return b.evalContext.jobScheduler
}

func (b *BlockContext) Logger() Logger {
	return b.evalContext.Logger
}

func (b *BlockContext) Stdout() io.Writer {
	return b.evalContext.Stdout
}

func (b *BlockContext) UserContext() interface{} {
	return b.evalContext.UserContext
}
