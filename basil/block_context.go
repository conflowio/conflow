package basil

import "context"

// BlockContext is passed to the block objects during evaluation
//go:generate counterfeiter . BlockContext
type BlockContext interface {
	Context() context.Context
	UserContext() interface{}
}

// BlockContextOverride contains block context override values
type BlockContextOverride struct {
	Context     context.Context
	UserContext interface{}
}

// BlockContextOverrider defines an interface to override a block context
type BlockContextOverrider interface {
	BlockContextOverride(BlockContext) BlockContextOverride
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
