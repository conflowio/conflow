package block

import (
	"context"

	"github.com/opsidian/basil/basil"
)

// NewContext creates a new block context
func NewContext(evalCtx basil.EvalContext, context context.Context, userContext interface{}) *Context {
	return &Context{
		evalCtx:     evalCtx,
		context:     context,
		userContext: userContext,
	}
}

type Context struct {
	evalCtx     basil.EvalContext
	context     context.Context
	userContext interface{}
}

func (c *Context) Context() context.Context {
	return c.context
}

func (c *Context) UserContext() interface{} {
	return c.userContext
}

func (c *Context) Publish(block basil.Block) {
	container := NewContainer(c.evalCtx, nil, block)
	c.evalCtx.Publish(container)
}
