// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package job

import (
	"fmt"
	"sync"
	"sync/atomic"
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
	lastID    uint64
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

// GenerateJobID generates a unique job identifier
func (m *Manager) GenerateJobID(id basil.ID) basil.ID {
	return basil.ID(fmt.Sprintf("%s@%d", id, atomic.AddUint64(&m.lastID, 1)))
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
			if job.Job.Cancel() {
				m.debug().ID("jobID", id).Msg("job cancelled")
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
func (m *Manager) Schedule(job basil.Job, decreasePending bool) {
	m.mu.Lock()
	if m.stopped {
		m.mu.Unlock()
		return
	}

	j := m.jobs[job.JobID()]
	if j == nil {
		j = &state{Job: job}
		m.jobs[job.JobID()] = j
	} else if j.RetryTimer != nil {
		j.RetryTimer.Stop()
		j.RetryTimer = nil
	}

	if decreasePending {
		m.pending--
	}
	m.running++
	j.Tries++

	m.debug().ID("jobID", job.JobID()).Msg("job scheduled")
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
		m.debug().ID("jobID", j.Job.JobID()).Msg("job finished")
	}

	m.mu.Unlock()
	if j == nil {
		panic(fmt.Errorf("job %q wasn't run by the job manager", id))
	}
}

// Retry must be called when a job failed but can be retried.
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
			ID("jobID", j.Job.JobID()).
			Int("tries", j.Tries).
			Int("max_tries", maxTries).
			Msg("job failed permanently")
		delete(m.jobs, id)
		return false
	}

	m.debug().
		ID("jobID", j.Job.JobID()).
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
