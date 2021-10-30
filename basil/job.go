// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

// JobScheduler is the job scheduler
//counterfeiter:generate . JobScheduler
type JobScheduler interface {
	ScheduleJob(job Job) error
}

// Job is a unit of work the scheduler can schedule and run
//counterfeiter:generate . Job
type Job interface {
	JobName() ID
	JobID() int
	SetJobID(int)
	Run()
	Lightweight() bool
}

//counterfeiter:generate . JobContainer
type JobContainer interface {
	Job
	Container
	Cancel() bool
	Pending() bool
	EvalStage() EvalStage
}
