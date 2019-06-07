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
func NewBlockContext(context context.Context, userContext interface{}, logger Logger) BlockContext {
	return &blockContext{
		context:     context,
		userContext: userContext,
		logger:      logger,
	}
}

type blockContext struct {
	context     context.Context
	userContext interface{}
	logger      Logger
}

func (b *blockContext) Context() context.Context {
	return b.context
}

func (b *blockContext) UserContext() interface{} {
	return b.userContext
}

func (b *blockContext) Logger() Logger {
	return b.logger
}
