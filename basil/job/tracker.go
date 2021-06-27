// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package job

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/opsidian/basil/basil"
)

type jobState struct {
	Tries      int
	Job        basil.Job
	RetryTimer *time.Timer
}

type Pending interface {
	Pending() bool
}

type Cancellable interface {
	Cancel() bool
}

// Tracker schedules and tracks jobs
type Tracker struct {
	name      basil.ID
	scheduler basil.JobScheduler
	logger    basil.Logger
	jobs      map[int]*jobState
	running   int
	pending   int
	stopped   bool
	mu        *sync.Mutex
}

// NewTracker creates a new job tracker
func NewTracker(name basil.ID, scheduler basil.JobScheduler, logger basil.Logger) *Tracker {
	return &Tracker{
		name:      name,
		scheduler: scheduler,
		logger:    logger.With().ID("scheduler", name).Logger(),
		jobs:      make(map[int]*jobState),
		mu:        &sync.Mutex{},
	}
}

// Stop cancels all remaining jobs
// It can be called multiple times
// It will return with the number of running jobs which couldn't be cancelled
func (t *Tracker) Stop() int {
	t.mu.Lock()

	if !t.stopped {
		t.debug().Msg("job manager stopping")
		t.stopped = true

		for id, job := range t.jobs {
			if jc, ok := job.Job.(Cancellable); ok && jc.Cancel() {
				t.debug().ID("jobName", job.Job.JobName()).Int("jobID", id).Msg("job cancelled")
				t.running--
			}
		}
		t.debug().Msg("job manager stopped")
	}

	running := t.running
	t.mu.Unlock()
	return running
}

// Schedule schedules a new job
func (t *Tracker) ScheduleJob(job basil.Job) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.stopped {
		return errors.New("job tracker was stopped")
	}

	var pending bool
	if jp, ok := job.(Pending); ok {
		pending = jp.Pending()
	}

	if err := t.scheduler.ScheduleJob(job); err != nil {
		return err
	}

	j := t.jobs[job.JobID()]
	if j == nil {
		j = &jobState{Job: job}
		t.jobs[job.JobID()] = j
	} else if j.RetryTimer != nil {
		j.RetryTimer.Stop()
		j.RetryTimer = nil
	}

	if pending {
		t.pending--
	}
	t.running++
	j.Tries++

	t.debug().ID("jobName", job.JobName()).Int("jobID", job.JobID()).Msg("job scheduled")

	return nil
}

// Finished must be called when a process finished
func (t *Tracker) Finished(id int) {
	t.done(id, "job finished")
}

func (t *Tracker) Failed(id int) {
	t.done(id, "job failed")
}

// Cancelled must be called when a process was cancelled
func (t *Tracker) Cancelled(id int) {
	t.done(id, "job cancelled")
}

func (t *Tracker) done(id int, msg string) {
	t.mu.Lock()

	j := t.jobs[id]
	if j != nil {
		t.running--
		delete(t.jobs, id)
		t.debug().ID("jobName", j.Job.JobName()).Int("jobID", id).Msg(msg)

	}

	t.mu.Unlock()
	if j == nil {
		panic(fmt.Errorf("job %q wasn't run by the job manager", id))
	}
}

// Retry must be called when a job failed but can be retried.
// It schedules the job again if there are any retries left
func (t *Tracker) Retry(id int, maxTries int, retryDelay time.Duration, f func(job basil.Job) func()) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	j, exists := t.jobs[id]
	if !exists {
		panic(fmt.Errorf("job %q wasn't run by the job manager", id))
	}

	t.running--

	if j.Tries >= maxTries {
		t.debug().
			ID("jobName", j.Job.JobName()).
			Int("jobID", id).
			Int("tries", j.Tries).
			Int("max_tries", maxTries).
			Msg("job failed permanently")
		delete(t.jobs, id)
		return false
	}

	t.debug().
		ID("jobName", j.Job.JobName()).
		Int("jobID", id).
		Int("tries", j.Tries).
		Int("max_tries", maxTries).
		Dur("retry_delay", retryDelay).
		Msg("job failed, retrying")

	// make sure we cancel any previous timers just to be sure
	if j.RetryTimer != nil {
		j.RetryTimer.Stop()
	}
	j.RetryTimer = time.AfterFunc(retryDelay, f(j.Job))
	return true
}

func (t *Tracker) AddPending(cnt int) {
	if cnt == 0 {
		return
	}

	t.mu.Lock()
	t.pending = t.pending + cnt
	t.mu.Unlock()
}

func (t *Tracker) RemovePending(cnt int) {
	if cnt == 0 {
		return
	}

	t.mu.Lock()
	t.pending = t.pending - cnt
	t.mu.Unlock()
}

// ActiveJobCount returns with the number of active (pending or running) jobs
func (t *Tracker) ActiveJobCount() int {
	t.mu.Lock()
	active := t.running + t.pending
	t.mu.Unlock()
	return active
}

// RunningJobCount returns with the number of running jobs
func (t *Tracker) RunningJobCount() int {
	t.mu.Lock()
	running := t.running
	t.mu.Unlock()
	return running
}

// PendingJobCount returns with the number of pending jobs
func (t *Tracker) PendingJobCount() int {
	t.mu.Lock()
	pending := t.pending
	t.mu.Unlock()
	return pending
}

func (t *Tracker) debug() basil.LogEvent {
	return t.logger.Debug().Int("active_jobs", t.running+t.pending)
}
