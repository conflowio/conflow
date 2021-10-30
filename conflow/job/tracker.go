// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package job

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/opsidian/conflow/conflow"
)

type jobState struct {
	Tries      int
	Job        conflow.Job
	RetryTimer *time.Timer
}

type Pending interface {
	Pending() bool
}

type Cancellable interface {
	Cancel() bool
}

type RetryBackoff func(i int) time.Duration

func ExponentialRetryBackoff(base float64, baseDelay, maxDelay time.Duration) RetryBackoff {
	return func(i int) time.Duration {
		backoff := float64(baseDelay.Nanoseconds()) * math.Pow(base, float64(i))
		if backoff > math.MaxInt64 {
			return maxDelay
		}

		res := time.Duration(backoff)
		if res > maxDelay {
			return maxDelay
		}

		return res
	}
}

// Tracker schedules and tracks jobs
type Tracker struct {
	name         conflow.ID
	scheduler    conflow.JobScheduler
	logger       conflow.Logger
	retryBackoff RetryBackoff
	jobs         map[int]*jobState
	running      int
	pending      int
	stopped      bool
	stoppedChan  chan struct{}
	mu           *sync.Mutex
}

// NewTracker creates a new job tracker
func NewTracker(name conflow.ID, scheduler conflow.JobScheduler, logger conflow.Logger, retryBackoff RetryBackoff) *Tracker {
	return &Tracker{
		name:         name,
		scheduler:    scheduler,
		logger:       logger.With().ID("scheduler", name).Logger(),
		retryBackoff: retryBackoff,
		jobs:         make(map[int]*jobState),
		mu:           &sync.Mutex{},
		stoppedChan:  make(chan struct{}),
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
		close(t.stoppedChan)

		for id, job := range t.jobs {
			if jc, ok := job.Job.(Cancellable); ok && jc.Cancel() {
				t.debug().ID("job_name", job.Job.JobName()).Int("job_id", id).Msg("job cancelled")
				t.running--
			}
		}
		t.debug().Msg("job manager stopped")
	}

	running := t.running
	t.mu.Unlock()
	return running
}

// ScheduleJob schedules a new job
func (t *Tracker) ScheduleJob(job conflow.Job) error {
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

	t.debug().ID("job_name", job.JobName()).Int("job_id", job.JobID()).Msg("job scheduled")

	return nil
}

// Succeeded must be called when a process finished successfully
func (t *Tracker) Succeeded(id int) {
	t.done(id, "job succeeded")
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
	defer t.mu.Unlock()

	j := t.jobs[id]
	if j == nil {
		panic(fmt.Errorf("job %q wasn't run by the job manager", id))
	}

	t.running--
	delete(t.jobs, id)
	t.debug().ID("job_name", j.Job.JobName()).Int("job_id", id).Msg(msg)
}

// Retry must be called when a job failed but can be retried.
// It schedules the job again if there are any retries left
// If Retry returns false, you are expected to call Failed()
func (t *Tracker) Retry(id int, limit int, delay time.Duration, reason string, f func()) bool {
	if limit == 0 {
		panic("Retry should not be called with limit = 0")
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	j, exists := t.jobs[id]
	if !exists {
		panic(fmt.Errorf("job %q wasn't run by the job manager", id))
	}

	// make sure we cancel any previous timers just to be sure
	if j.RetryTimer != nil {
		j.RetryTimer.Stop()
		j.RetryTimer = nil
	}

	if limit > 0 && j.Tries >= limit+1 {
		return false
	}

	t.running--

	if delay <= 0 && t.retryBackoff != nil {
		delay = t.retryBackoff(j.Tries - 1)
	}

	t.debug().
		ID("job_name", j.Job.JobName()).
		Int("job_id", id).
		Int("tries", j.Tries).
		Int("retry_limit", limit).
		Str("retry_delay", delay.String()).
		Str("retry_reason", reason).
		Msg("job failed, retrying")

	if delay <= 0 {
		go f()
	} else {
		j.RetryTimer = time.AfterFunc(delay, f)
	}

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

func (t *Tracker) debug() conflow.LogEvent {
	return t.logger.Debug().Int("active_jobs", t.running+t.pending)
}
