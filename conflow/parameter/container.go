// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parameter

import (
	"github.com/opsidian/conflow/conflow"
	"github.com/opsidian/parsley/parsley"
)

var _ conflow.ParameterContainer = &Container{}

// Container is a parameter container
type Container struct {
	evalCtx *conflow.EvalContext
	node    conflow.ParameterNode
	parent  conflow.BlockContainer
	value   interface{}
	err     parsley.Error
	jobID   int
	wgs     []conflow.WaitGroup
	pending bool
}

// NewContainer creates a new parameter container
func NewContainer(
	evalCtx *conflow.EvalContext,
	node conflow.ParameterNode,
	parent conflow.BlockContainer,
	value interface{},
	wgs []conflow.WaitGroup,
	pending bool,
) *Container {
	return &Container{
		evalCtx: evalCtx,
		node:    node,
		parent:  parent,
		value:   value,
		wgs:     wgs,
		pending: pending,
	}
}

// Node returns with the parameter node
func (c *Container) Node() conflow.Node {
	return c.node
}

// JobName returns with the job name
func (c *Container) JobName() conflow.ID {
	return c.node.ID()
}

// ID returns with the block id
func (c *Container) JobID() int {
	return c.jobID
}

// SetJobID sets the job id
func (c *Container) SetJobID(id int) {
	c.jobID = id
}

func (c *Container) Lightweight() bool {
	return true
}

func (c *Container) Cancel() bool {
	return c.evalCtx.Cancel()
}

func (c *Container) EvalStage() conflow.EvalStage {
	return c.node.EvalStage()
}

func (c *Container) Pending() bool {
	return c.pending
}

// BlockContainer returns with the parent block container
func (c *Container) BlockContainer() conflow.BlockContainer {
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
	defer func() {
		if c.parent != nil {
			c.parent.SetChild(c)
		}

		c.evalCtx.Cancel()
	}()

	if !c.evalCtx.Run() || c.value != nil {
		return
	}

	c.value, c.err = parsley.EvaluateNode(c.evalCtx, c.node)
}

// Close notifies all wait groups
func (c *Container) Close() {
	for _, wg := range c.wgs {
		wg.Done(c.err)
	}
}

// WaitGroups returns nil
func (c *Container) WaitGroups() []conflow.WaitGroup {
	return c.wgs
}
