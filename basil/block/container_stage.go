// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package block

import (
	"github.com/opsidian/basil/basil/job"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/parsley/parsley"
)

type containerStage struct {
	name        basil.ID
	jobID       int
	sem         job.Semaphore
	lightweight bool
	f           func() (int64, error)
	container   *Container
}

func newContainerStage(
	container *Container,
	name basil.ID,
	lightweight bool,
	f func() (int64, error),
) *containerStage {
	return &containerStage{
		container:   container,
		name:        name,
		lightweight: lightweight,
		f:           f,
	}
}

func (c *containerStage) JobName() basil.ID {
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

	defer func() {
		if r := recover(); r != nil {
			c.container.jobManager.Failed(c.jobID)
			c.container.errChan <- parsley.NewErrorf(
				c.container.node.Pos(),
				"%s stage panicked in %q: %s",
				c.name,
				c.container.ID(),
				r,
			)
		}
	}()

	nextStage, err := c.f()
	if err != nil {
		c.container.jobManager.Failed(c.jobID)
		c.container.errChan <- parsley.NewError(c.container.node.Pos(), err)
		return
	}

	c.container.jobManager.Finished(c.jobID)
	c.container.stateChan <- nextStage
}

func (c *containerStage) Cancel() bool {
	if c.sem.Cancel() {
		c.container.jobManager.Cancelled(c.jobID)
		return true
	}
	return false
}

func (c *containerStage) Lightweight() bool {
	return c.lightweight
}
