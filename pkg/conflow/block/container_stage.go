// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package block

import (
	"runtime/debug"

	"github.com/conflowio/parsley/parsley"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/conflow/job"
)

type containerStage struct {
	container   *Container
	evalStage   conflow.EvalStage
	name        conflow.ID
	lightweight bool
	f           func() (int64, error)
	jobID       int
	sem         job.Semaphore
	retryConfig conflow.RetryConfig
}

func newContainerStage(
	container *Container,
	name conflow.ID,
	evalStage conflow.EvalStage,
	lightweight bool,
	retryConfig conflow.RetryConfig,
	f func() (int64, error),
) *containerStage {
	return &containerStage{
		container:   container,
		name:        name,
		evalStage:   evalStage,
		lightweight: lightweight,
		retryConfig: retryConfig,
		f:           f,
	}
}

func (c *containerStage) JobName() conflow.ID {
	return c.name
}

func (c *containerStage) JobID() int {
	return c.jobID
}

func (c *containerStage) SetJobID(id int) {
	c.jobID = id
}

func (c *containerStage) Run() {
	if !c.sem.Run() {
		return
	}

	jobID := c.jobID

	defer func() {
		if r := recover(); r != nil {
			c.container.jobTracker.Failed(jobID)
			c.container.errChan <- parsley.NewErrorf(
				c.container.node.Pos(),
				"%s stage panicked in %q: %s\n%s",
				c.name,
				c.container.Node().ID(),
				r,
				string(debug.Stack()),
			)
		}
	}()

	nextStage, err := c.f()
	if err != nil {
		if re, ok := err.(retryError); ok && c.retryConfig.Limit != 0 {
			c.sem.Reset()

			if c.container.jobTracker.Retry(jobID, c.retryConfig.Limit, re.Duration, re.Reason, func() {
				if err := c.container.jobTracker.ScheduleJob(c); err != nil {
					c.container.jobTracker.Failed(jobID)
					c.container.errChan <- parsley.NewError(c.container.node.Pos(), err)
				}
			}) {
				return
			}
		}

		c.container.jobTracker.Failed(jobID)
		c.container.errChan <- parsley.NewError(c.container.node.Pos(), err)
		return
	}

	c.container.jobTracker.Succeeded(jobID)
	c.container.stateChan <- nextStage
}

func (c *containerStage) Cancel() bool {
	if c.sem.Cancel() {
		c.container.jobTracker.Cancelled(c.jobID)
		return true
	}
	return false
}

func (c *containerStage) Lightweight() bool {
	return c.lightweight
}

func (c *containerStage) EvalStage() conflow.EvalStage {
	return c.evalStage
}
