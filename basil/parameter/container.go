// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parameter

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/util"
	"github.com/opsidian/parsley/parsley"
)

var _ basil.ParameterContainer = &Container{}

// Container is a parameter container
type Container struct {
	evalCtx        basil.EvalContext
	node           basil.ParameterNode
	blockContainer basil.BlockContainer
	value          interface{}
	err            parsley.Error
}

// NewContainer creates a new parameter container
func NewContainer(evalCtx basil.EvalContext, node basil.ParameterNode, blockContainer basil.BlockContainer) *Container {
	return &Container{
		evalCtx:        evalCtx,
		node:           node,
		blockContainer: blockContainer,
	}
}

// ID returns with the parameter id
func (c *Container) ID() basil.ID {
	return c.node.ID()
}

// Node returns with the parameter node
func (c *Container) Node() basil.Node {
	return c.node
}

// BlockContainer returns with the parent block container
func (c *Container) BlockContainer() basil.BlockContainer {
	return c.blockContainer
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

// Close does nothing
func (c *Container) Close() {}

// WaitGroups returns nil
func (c *Container) WaitGroups() []*util.WaitGroup {
	return nil
}
