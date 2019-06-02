package basil

import "context"

// BlockContext is passed to the block objects during evaluation
//go:generate counterfeiter . BlockContext
type BlockContext interface {
	Context() context.Context
	UserContext() interface{}
}

// Context defines an interface about creating a new context
type Contexter interface {
	Context(ctx context.Context) (context.Context, context.CancelFunc)
}

// UserContexter defines an interface about creating a new user context
type UserContexter interface {
	UserContext(ctx interface{}) interface{}
}

// NewBlockContext creates a new block context
func NewBlockContext(context context.Context, userContext interface{}) BlockContext {
	return &blockContext{
		context:     context,
		userContext: userContext,
	}
}

type blockContext struct {
	context     context.Context
	userContext interface{}
}

func (b *blockContext) Context() context.Context {
	return b.context
}

func (b *blockContext) UserContext() interface{} {
	return b.userContext
}
