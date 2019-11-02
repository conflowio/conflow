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
	name      basil.ID
	scheduler basil.JobScheduler
	logger    basil.Logger
	jobs      map[int]*state
	running   int
	pending   int
	stopped   bool
	mu        *sync.Mutex
}

// NewManager creates a new job manager
func NewManager(name basil.ID, scheduler basil.JobScheduler, logger basil.Logger) *Manager {
	return &Manager{
		name:      name,
		scheduler: scheduler,
		logger:    logger.With().ID("jobManager", name).Logger(),
		jobs:      make(map[int]*state),
		mu:        &sync.Mutex{},
	}
}

// Stop cancels all remaining jobs
// It can be called multiple times
// It will return with the number of running jobs which couldn't be cancelled
func (m *Manager) Stop() int {
	m.mu.Lock()

	if !m.stopped {
		m.debug().Msg("job manager stopping")
		m.stopped = true

		for id, job := range m.jobs {
			if j, ok := job.Job.(basil.Cancellable); ok && j.Cancel() {
				m.debug().ID("jobName", job.Job.JobName()).Int("jobID", id).Msg("job cancelled")
				m.running--
			}
		}
		m.debug().Msg("job manager stopped")
	}

	running := m.running
	m.mu.Unlock()
	return running
}

// Schedule schedules a new job
func (m *Manager) ScheduleJob(job basil.Job, pending bool) {
	m.mu.Lock()
	if m.stopped {
		m.mu.Unlock()
		return
	}

	m.scheduler.ScheduleJob(job)

	j := m.jobs[job.JobID()]
	if j == nil {
		j = &state{Job: job}
		m.jobs[job.JobID()] = j
	} else if j.RetryTimer != nil {
		j.RetryTimer.Stop()
		j.RetryTimer = nil
	}

	if pending {
		m.pending--
	}
	m.running++
	j.Tries++

	m.debug().ID("jobName", job.JobName()).Int("jobID", job.JobID()).Msg("job scheduled")
	m.mu.Unlock()
}

// Finished must be called when a process finished
func (m *Manager) Finished(id int) {
	m.done(id, "job finished")
}

func (m *Manager) Failed(id int) {
	m.done(id, "job failed")
}

// Cancelled must be called when a process was cancelled
func (m *Manager) Cancelled(id int) {
	m.done(id, "job cancelled")
}

func (m *Manager) done(id int, msg string) {
	m.mu.Lock()

	j := m.jobs[id]
	if j != nil {
		m.running--
		delete(m.jobs, id)
		m.debug().ID("jobName", j.Job.JobName()).Int("jobID", id).Msg(msg)

	}

	m.mu.Unlock()
	if j == nil {
		panic(fmt.Errorf("job %q wasn't run by the job manager", id))
	}
}

// Retry must be called when a job failed but can be retried.
// It schedules the job again if there are any retries left
func (m *Manager) Retry(id int, maxTries int, retryDelay time.Duration, f func(job basil.Job) func()) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	j, exists := m.jobs[id]
	if !exists {
		panic(fmt.Errorf("job %q wasn't run by the job manager", id))
	}

	m.running--

	if j.Tries >= maxTries {
		m.debug().
			ID("jobName", j.Job.JobName()).
			Int("jobID", id).
			Int("tries", j.Tries).
			Int("max_tries", maxTries).
			Msg("job failed permanently")
		delete(m.jobs, id)
		return false
	}

	m.debug().
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

func (m *Manager) AddPending(cnt int) {
	if cnt == 0 {
		return
	}

	m.mu.Lock()
	m.pending = m.pending + cnt
	m.mu.Unlock()
}

// ActiveJobCount returns with the number of active (pending or running) jobs
func (m *Manager) ActiveJobCount() int {
	m.mu.Lock()
	active := m.running + m.pending
	m.mu.Unlock()
	return active
}

// RunningJobCount returns with the number of running jobs
func (m *Manager) RunningJobCount() int {
	m.mu.Lock()
	running := m.running
	m.mu.Unlock()
	return running
}

// PendingJobCount returns with the number of pending jobs
func (m *Manager) PendingJobCount() int {
	m.mu.Lock()
	pending := m.pending
	m.mu.Unlock()
	return pending
}

func (m *Manager) debug() basil.LogEvent {
	return m.logger.Debug().Int("active_jobs", m.running+m.pending)
}
