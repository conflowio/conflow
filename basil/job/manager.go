// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package job

import (
	"fmt"
	"sync"
	"time"

	"github.com/opsidian/basil/basil"
)

type state struct {
	Tries      int
	Job        basil.Job
	RetryTimer *time.Timer
}

// JobManager tracks and schedules jobs
type Manager struct {
	id        basil.ID
	scheduler basil.Scheduler
	logger    basil.Logger
	jobs      map[basil.ID]*state
	running   int
	pending   int
	stopped   bool
	mu        *sync.Mutex
}

// NewManager creates a new job manager
func NewManager(id basil.ID, scheduler basil.Scheduler, logger basil.Logger) *Manager {
	return &Manager{
		id:        id,
		scheduler: scheduler,
		logger:    logger.With().ID("id", id).Logger(),
		jobs:      make(map[basil.ID]*state),
		mu:        &sync.Mutex{},
	}
}

// Stop cancels all remaining jobs
// It can be called multiple times
// It will return with the number of jobs which couldn't be cancelled
func (m *Manager) Stop() int {
	m.mu.Lock()

	if !m.stopped {
		m.debug().Msg("job manager stopping")
		m.stopped = true

		for _, job := range m.jobs {
			if job == nil {
				m.pending--
				continue
			}
			if job.RetryTimer != nil && job.RetryTimer.Stop() {
				m.pending--
			}
			if job.Job.Cancel() {
				m.debug().ID("jobID", job.Job.ID()).Msg("job cancelled")
				m.running--
			}
		}
		m.debug().Msg("job manager stopped")
	}

	remaining := m.running + m.pending
	m.mu.Unlock()
	return remaining
}

// Pending records a job which will run in the future
func (m *Manager) Pending(id basil.ID) {
	m.mu.Lock()
	if m.stopped {
		m.mu.Unlock()
		return
	}

	m.pending++
	m.jobs[id] = nil

	m.debug().ID("jobID", id).Msg("job pending")

	m.mu.Unlock()
}

// Schedule schedules a new job
func (m *Manager) Schedule(job basil.Job) {
	m.mu.Lock()
	if m.stopped {
		m.mu.Unlock()
		return
	}

	j, exists := m.jobs[job.ID()]

	if exists && j == nil {
		m.pending--
	}

	if j == nil {
		m.jobs[job.ID()] = &state{Job: job}
	} else if j.RetryTimer != nil {
		j.RetryTimer.Stop()
		j.RetryTimer = nil
		m.pending--
	}

	m.jobs[job.ID()].Tries++

	m.running++

	m.debug().ID("jobID", job.ID()).Msg("job scheduled")

	m.scheduler.Schedule(job)

	m.mu.Unlock()
}

// Finished must be called when a process finished
func (m *Manager) Finished(id basil.ID) {
	m.mu.Lock()

	j := m.jobs[id]
	if j != nil {
		m.running--
		delete(m.jobs, id)
		m.debug().ID("jobID", j.Job.ID()).Msg("job finished")
	}

	m.mu.Unlock()
	if j == nil {
		panic(fmt.Errorf("job %q wasn't run by the job manager", id))
	}
}

// FailedRetry must be called when a job failed but can be retried.
// It schedules the job again if there are any retries left
func (m *Manager) Retry(id basil.ID, maxTries int, retryDelay time.Duration, f func(job basil.Job) func()) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	j, exists := m.jobs[id]
	if !exists {
		panic(fmt.Errorf("job %q wasn't run by the job manager", id))
	}

	m.running--
	if j.Tries >= maxTries {
		m.debug().
			ID("jobID", j.Job.ID()).
			Int("tries", j.Tries).
			Int("max_tries", maxTries).
			Msg("job failed permanently")
		delete(m.jobs, id)
		return false
	}

	m.pending++

	m.debug().
		ID("jobID", j.Job.ID()).
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

// PendingJobCount returns with the number of pending jobs
func (m *Manager) PendingJobCount() int {
	m.mu.Lock()
	pending := m.pending
	m.mu.Unlock()
	return pending
}

// RunningJobCount returns with the number of running jobs
func (m *Manager) RunningJobCount() int {
	m.mu.Lock()
	running := m.running
	m.mu.Unlock()
	return running
}

// RemainingJobCount returns with the number of pending and running jobs
func (m *Manager) RemainingJobCount() int {
	m.mu.Lock()
	remaining := m.pending + m.running
	m.mu.Unlock()
	return remaining
}

func (m *Manager) debug() basil.LogEvent {
	e := m.logger.Debug()
	if e.Enabled() {
		e = e.
			Int("pending", m.pending).
			Int("running", m.running)
	}
	return e
}
