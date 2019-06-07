// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

// Scheduler is the job scheduler
type Scheduler interface {
	Start()
	Stop()
	Schedule(Job)
}

// Worker is an interface for a job queue processor
type Worker interface {
	Start()
	Stop()
}
