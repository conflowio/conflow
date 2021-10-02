// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

// JobScheduler is the job scheduler
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . JobScheduler
type JobScheduler interface {
	ScheduleJob(job Job) error
}

// Job is a unit of work the scheduler can schedule and run
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Job
type Job interface {
	JobName() ID
	JobID() int
	SetJobID(int)
	Run()
	Lightweight() bool
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . JobContainer
type JobContainer interface {
	Job
	Container
	Cancel() bool
	Pending() bool
	EvalStage() EvalStage
}
