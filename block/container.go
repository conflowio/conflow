package block

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/parsley/parsley"
)

// Container is a block container
type Container struct {
	block       basil.Block
	extraParams map[basil.ID]interface{}
	interpreter Interpreter
}

// NewContainer creates a new block container instance
func NewContainer(
	block basil.Block,
	interpreter Interpreter,
) *Container {
	return &Container{
		block:       block,
		interpreter: interpreter,
	}
}

// Block returns with the block
func (c *Container) Block() basil.Block {
	return c.block
}

// Param returns with the parameter value
func (c *Container) Param(name basil.ID) interface{} {
	if val, ok := c.extraParams[name]; ok {
		return val
	}

	return c.interpreter.Param(c.block, name)
}

// SetParam sets the parameter value
func (c *Container) SetParam(ctx *basil.EvalContext, name basil.ID, node parsley.Node) parsley.Error {
	if _, exists := c.interpreter.Params()[name]; exists {
		return c.interpreter.Update(ctx, c.block, name, node)
	}
	value, err := node.Value(ctx)
	if err != nil {
		return err
	}

	c.extraParams[name] = value
	return nil
}
