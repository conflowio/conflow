package basil

import (
	"context"
)

// EvalStage means an evaluation stage (default, pre or post)
type EvalStage uint8

// Evaluation stages
const (
	EvalStageDefault EvalStage = iota
	EvalStagePre
	EvalStagePost
)

// EvalStageAware defines an interface to return with the evaluation stage
type EvalStageAware interface {
	EvalStage() EvalStage
}

// EvalContext is the evaluation context
type EvalContext struct {
	processCtx             context.Context
	userCtx                interface{}
	blockContainerRegistry BlockContainerRegistry
}

// EvalContextAware defines an interface to return a child evaluation content
type EvalContextAware interface {
	EvalContext(*EvalContext) *EvalContext
}

// EvalContextConfig is used to pass overrides when creating a child context
type EvalContextConfig struct {
	ProcessCtx context.Context
	UserCtx    interface{}
}

// NewEvalContext returns with a new context
func NewEvalContext(
	processCtx context.Context,
	userCtx interface{},
	blockContainerRegistry BlockContainerRegistry,
) *EvalContext {
	return &EvalContext{
		processCtx:             processCtx,
		userCtx:                userCtx,
		blockContainerRegistry: blockContainerRegistry,
	}
}

// New creates a new child context
func (e *EvalContext) New(config EvalContextConfig) *EvalContext {
	ctx := &EvalContext{
		blockContainerRegistry: e.blockContainerRegistry,
	}
	if config.ProcessCtx != nil {
		ctx.processCtx = config.ProcessCtx
	} else {
		ctx.processCtx = e.processCtx
	}
	if config.UserCtx != nil {
		ctx.userCtx = config.UserCtx
	} else {
		ctx.userCtx = e.userCtx
	}

	return ctx
}

// ProcessContext returns with the Go context
func (e *EvalContext) ProcessContext() context.Context {
	return e.processCtx
}

// UserContext returns with a user-specific context
func (e *EvalContext) UserContext() interface{} {
	return e.userCtx
}

// BlockContainerRegistry returns with the block container registry
func (e *EvalContext) BlockContainerRegistry() BlockContainerRegistry {
	return e.blockContainerRegistry
}
