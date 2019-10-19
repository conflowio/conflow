// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parameter

import (
	"sync/atomic"

	"github.com/opsidian/basil/basil/job"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/util"
	"github.com/opsidian/parsley/parsley"
)

var _ basil.ParameterContainer = &Container{}

const (
	containerStateWaiting int64 = iota
	containerStateRunning
	containerStateCancelled
)

// Container is a parameter container
type Container struct {
	evalCtx *basil.EvalContext
	node    basil.ParameterNode
	parent  basil.BlockContainer
	value   interface{}
	err     parsley.Error
	state   int64
	jobID   basil.ID
}

// NewContainer creates a new parameter container
func NewContainer(
	evalCtx *basil.EvalContext,
	node basil.ParameterNode,
	parent basil.BlockContainer,
	jobManager *job.Manager,
) *Container {
	return &Container{
		evalCtx: evalCtx,
		node:    node,
		parent:  parent,
		jobID:   jobManager.GenerateJobID(node.ID()),
	}
}

// ID returns with the parameter id
func (c *Container) ID() basil.ID {
	return c.node.ID()
}

// ID returns with the block id
func (c *Container) JobID() basil.ID {
	return c.jobID
}

// Node returns with the parameter node
func (c *Container) Node() basil.Node {
	return c.node
}

// BlockContainer returns with the parent block container
func (c *Container) BlockContainer() basil.BlockContainer {
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
	if !atomic.CompareAndSwapInt64(&c.state, containerStateWaiting, containerStateRunning) {
		return
	}

	c.value, c.err = c.node.Value(c.evalCtx)
	c.parent.SetChild(c)
}

func (c *Container) Cancel() bool {
	c.evalCtx.Cancel()
	return atomic.CompareAndSwapInt64(&c.state, containerStateWaiting, containerStateCancelled)
}

func (c *Container) Lightweight() bool {
	return true
}

// Close does nothing
func (c *Container) Close() {}

// WaitGroups returns nil
func (c *Container) WaitGroups() []*util.WaitGroup {
	return nil
}
