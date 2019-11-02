// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package test

import (
	"sync/atomic"

	"github.com/opsidian/basil/basil"
)

// Scheduler is a test scheduler, it will simply run the given job in a goroutine in the background
type Scheduler struct {
	lastID int64
}

func (s *Scheduler) ScheduleJob(job basil.Job) {
	job.SetJobID(int(atomic.AddInt64(&s.lastID, 1)))
	go job.Run()
}
