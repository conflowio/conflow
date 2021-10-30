// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package job

import (
	"errors"
	"sync/atomic"

	"github.com/opsidian/conflow/conflow"
)

// Scheduler handles workers and schedules jobs
type Scheduler struct {
	logger       conflow.Logger
	maxWorkers   int
	maxQueueSize int
	workers      []Worker
	jobQueue     chan conflow.Job
	lastID       int64
	stopped      bool
	stoppedChan  chan struct{}
}

// NewScheduler creates a new scheduler instance
func NewScheduler(logger conflow.Logger, maxWorkers int, maxQueueSize int) *Scheduler {
	return &Scheduler{
		logger:       logger,
		workers:      make([]Worker, maxWorkers),
		maxWorkers:   maxWorkers,
		maxQueueSize: maxQueueSize,
		jobQueue:     make(chan conflow.Job, maxQueueSize),
		stoppedChan:  make(chan struct{}),
	}
}

// Start creates and starts the workers
func (s *Scheduler) Start() {
	s.logger.Debug().
		Int("workers", s.maxWorkers).
		Int("queueSize", s.maxQueueSize).
		Msg("job scheduler starting")

	for i := 0; i < s.maxWorkers; i++ {
		s.workers[i] = NewWorker(s.logger, s.jobQueue)
		s.workers[i].Start()
	}
}

// Stop stops all the workers and the dispatcher process
func (s *Scheduler) Stop() {
	s.stopped = true
	close(s.stoppedChan)

	for i := 0; i < s.maxWorkers; i++ {
		s.workers[i].Stop()
	}

	s.logger.Debug().Msg("job scheduler stopped")
}

// ScheduleJob schedules a new job
func (s *Scheduler) ScheduleJob(job conflow.Job) error {
	if s.stopped {
		return errors.New("job scheduler was stopped")
	}

	if job.JobID() == 0 {
		job.SetJobID(int(atomic.AddInt64(&s.lastID, 1)))
	}

	if job.Lightweight() {
		go func() {
			job.Run()
		}()
		return nil
	} else {
		select {
		case s.jobQueue <- job:
			return nil
		case <-s.stoppedChan: // If jobQueue was full, we still want to return if the scheduler was stopped
			return errors.New("job scheduler was stopped")
		}
	}
}
