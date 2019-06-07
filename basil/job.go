// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

// Job is a unit of work the scheduler can schedule and run
type Job interface {
	Run()
}

// JobFunc is a function which implements the job interface
type JobFunc func()

func (f JobFunc) Run() {
	f()
}

type job struct {
	ctx EvalContext
	id  ID
	job Job
}

// NewJob creates a new job
func NewJob(ctx EvalContext, id ID, j Job) Job {
	return &job{
		ctx: ctx,
		id:  id,
		job: j,
	}
}

// ID returns the job id
func (j *job) ID() ID {
	return j.id
}

// Run runs the wrapped job if it wasn't cancelled
func (j *job) Run() {
	if j.ctx.BlockContext().Context().Err() != nil {
		return
	}

	j.job.Run()
}
