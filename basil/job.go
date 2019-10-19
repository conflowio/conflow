// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

import "sync/atomic"

const (
	jobStateWaiting int64 = iota
	jobStateRunning
	jobStateCanceled
)

// Job is a unit of work the scheduler can schedule and run
//go:generate counterfeiter . Job
type Job interface {
	JobID() ID
	Run()
	Cancel() bool
	Lightweight() bool
}

type job struct {
	ctx         *EvalContext
	id          ID
	state       int64
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

// JobID returns the job id
func (j *job) JobID() ID {
	return j.id
}

// Run runs the wrapped job if it wasn't cancelled
func (j *job) Run() {
	if !atomic.CompareAndSwapInt64(&j.state, jobStateWaiting, jobStateRunning) {
		return
	}

	j.f()
}

func (j *job) Cancel() bool {
	j.ctx.Cancel()
	return atomic.CompareAndSwapInt64(&j.state, jobStateWaiting, jobStateCanceled)
}

// Lightweight returns true if the job doesn't need to be scheduled on the main job queue
// These jobs will usually run in a goroutine.
func (j *job) Lightweight() bool {
	return j.lightweight
}
