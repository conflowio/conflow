// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

// Job is a unit of work the scheduler can schedule and run
type Job interface {
	Run()
	Lightweight() bool
}

type job struct {
	ctx         *EvalContext
	id          ID
	lightweight bool
	f           func()
}

// NewJob creates a new job
func NewJob(ctx *EvalContext, id ID, lightweight bool, f func()) Job {
	return &job{
		ctx:         ctx,
		id:          id,
		lightweight: lightweight,
		f:           f,
	}
}

// ID returns the job id
func (j *job) ID() ID {
	return j.id
}

// Run runs the wrapped job if it wasn't cancelled
func (j *job) Run() {
	if j.ctx.Context.Err() != nil {
		return
	}

	j.f()
}

// Lightweight returns true if the job doesn't need to be scheduled on the main job queue
// These jobs will usually run in a goroutine.
func (j *job) Lightweight() bool {
	return j.lightweight
}
