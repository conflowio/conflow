package block

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/parsley/parsley"
)

// Container is a block container
type Container struct {
	id          basil.ID
	block       basil.Block
	extraParams map[basil.ID]interface{}
	interpreter Interpreter
}

// NewContainer creates a new block container instance
func NewContainer(
	id basil.ID,
	block basil.Block,
	interpreter Interpreter,
) *Container {
	return &Container{
		id:          id,
		block:       block,
		interpreter: interpreter,
	}
}

// Block returns with the block
func (c *Container) ID() basil.ID {
	return c.id
}

// Block returns with the block
func (c *Container) Block() basil.Block {
	return c.block
}

// Param returns with the parameter value
func (c *Container) Param(name basil.ID) interface{} {
	if c.extraParams != nil {
		if val, ok := c.extraParams[name]; ok {
			return val
		}
	}

	return c.interpreter.Param(c.block, name)
}

// SetParam sets the parameter value
func (c *Container) SetParam(ctx *basil.EvalContext, name basil.ID, node parsley.Node) parsley.Error {
	if _, exists := c.interpreter.Params()[name]; exists {
		return c.interpreter.SetParam(ctx, c.block, name, node)
	}
	value, err := node.Value(ctx)
	if err != nil {
		return err
	}

	if c.extraParams == nil {
		c.extraParams = make(map[basil.ID]interface{})
	}
	c.extraParams[name] = value

	return nil
}
