package basil

import (
	"context"
	"fmt"
)

// EvalStage means an evaluation stage (default, pre or post)
type EvalStage int8

// Evaluation stages
const (
	EvalStageInit  EvalStage = -1
	EvalStageMain  EvalStage = 0
	EvalStageClose EvalStage = 1
)

// EvalContext is the evaluation context
type EvalContext struct {
	blockContext BlockContext
	containers   map[ID]BlockContainer
}

// NewEvalContext returns with a new evaluation context
func NewEvalContext(
	ctx context.Context,
	userCtx interface{},
) *EvalContext {
	return &EvalContext{
		blockContext: &blockContext{
			context:     ctx,
			userContext: userCtx,
		},
		containers: make(map[ID]BlockContainer, 16),
	}
}

// ProcessContext returns with the Go context
func (e *EvalContext) New(blockContextOverride BlockContextOverride) *EvalContext {
	blockContext := &blockContext{}

	if blockContextOverride.Context != nil {
		blockContext.context = blockContextOverride.Context
	} else {
		blockContext.context = e.blockContext.Context()
	}

	if blockContextOverride.UserContext != nil {
		blockContext.userContext = blockContextOverride.UserContext
	} else {
		blockContext.userContext = e.blockContext.UserContext()
	}

	return &EvalContext{
		blockContext: blockContext,
		containers:   e.containers,
	}
}

// ProcessContext returns with the Go context
func (e *EvalContext) BlockContext() BlockContext {
	return e.blockContext
}

// BlockContainer returns with the given block container instance if it exists
func (e *EvalContext) BlockContainer(id ID) (BlockContainer, bool) {
	node, ok := e.containers[id]
	return node, ok
}

// AddBlockContainer adds a new block container instance
// It returns with an error if a block with the same id was already registered
func (e *EvalContext) AddBlockContainer(b BlockContainer) error {
	id := b.ID()
	if _, exists := e.containers[id]; exists {
		return fmt.Errorf("%q is already defined, please use a globally unique identifier", id)
	}

	e.containers[id] = b

	return nil
}
