package parameter

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/parsley/parsley"
)

var _ basil.ParameterContainer = &Container{}

// Container is a parameter container
type Container struct {
	evalCtx basil.EvalContext
	node    basil.ParameterNode
	parent  basil.BlockContainer
	value   interface{}
	err     parsley.Error
}

// NewContainer creates a new parameter container
func NewContainer(evalCtx basil.EvalContext, node basil.ParameterNode, parent basil.BlockContainer) *Container {
	return &Container{
		evalCtx: evalCtx,
		node:    node,
		parent:  parent,
	}
}

// ID returns with the parameter id
func (c *Container) ID() basil.ID {
	return c.node.ID()
}

// Node returns with the parameter node
func (c *Container) Node() basil.ParameterNode {
	return c.node
}

// Parent returns with the parent block container
func (c *Container) Parent() basil.BlockContainer {
	return c.parent
}

// Value returns with the parameter value or an evaluation error
func (c *Container) Value() (interface{}, parsley.Error) {
	if c.err != nil {
		return nil, c.err
	}

	return c.value, nil
}

// Run evaluates the parameter
func (c *Container) Run() {
	c.value, c.err = c.node.Value(c.evalCtx)
}
