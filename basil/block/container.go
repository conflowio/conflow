package block

import (
	"fmt"
	"reflect"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/parsley/parsley"
)

// Container is a block container
type Container struct {
	id          basil.ID
	block       basil.Block
	extraParams map[basil.ID]interface{}
	interpreter basil.BlockInterpreter
}

// NewContainer creates a new block container instance
func NewContainer(
	id basil.ID,
	block basil.Block,
	interpreter basil.BlockInterpreter,
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

// EvaluateChildNode evaluates
func (c *Container) EvaluateChildNode(ctx *basil.EvalContext, node basil.Node) parsley.Error {
	switch n := node.(type) {
	case basil.BlockNode:
		return c.setBlock(ctx, n)
	case basil.ParameterNode:
		return c.setParam(ctx, n)
	default:
		panic(fmt.Errorf("Invalid node type: %s", reflect.TypeOf(node)))
	}
	return nil
}

func (c *Container) setBlock(ctx *basil.EvalContext, node basil.BlockNode) parsley.Error {
	value, err := node.Value(ctx)
	if err != nil {
		return err
	}
	return c.interpreter.SetBlock(ctx, c.block, node.BlockType(), value)
}

func (c *Container) setParam(ctx *basil.EvalContext, node basil.ParameterNode) parsley.Error {
	if !node.IsDeclaration() {
		return c.interpreter.SetParam(ctx, c.block, node.Name(), node)
	}

	value, err := node.Value(ctx)
	if err != nil {
		return err
	}

	if c.extraParams == nil {
		c.extraParams = make(map[basil.ID]interface{})
	}
	c.extraParams[node.Name()] = value

	return nil
}
