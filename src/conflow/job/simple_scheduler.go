// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package job

import "github.com/conflowio/conflow/src/conflow"

type SimpleScheduler struct{}

func (s SimpleScheduler) ScheduleJob(job conflow.Job) error {
	go func() {
		job.Run()
	}()

	return nil
}
