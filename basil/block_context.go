// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

import "context"

// BlockContext is passed to the block objects during evaluation
//go:generate counterfeiter . BlockContext
type BlockContext interface {
	Context() context.Context
	UserContext() interface{}
	Logger() Logger
	PublishBlock(Block) error
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

func (b *blockContext) Context() context.Context {
	return b.evalContext.Context
}

func (b *blockContext) UserContext() interface{} {
	return b.evalContext.UserContext
}

func (b *blockContext) Logger() Logger {
	return b.evalContext.Logger
}

func (b *blockContext) PublishBlock(block Block) error {
	return b.container.PublishBlock(block)
}
