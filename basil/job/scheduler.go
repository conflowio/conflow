// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package job

import (
	"github.com/opsidian/basil/basil"
)

// Scheduler handles workers and schedules jobs
type Scheduler struct {
	logger       basil.Logger
	maxWorkers   int
	maxQueueSize int
	workers      []basil.Worker
	workerPool   chan chan basil.Job
	jobQueue     chan basil.Job
	quit         chan bool
}

// NewScheduler creates a new scheduler instance
func NewScheduler(logger basil.Logger, maxWorkers int, maxQueueSize int) *Scheduler {
	return &Scheduler{
		logger:       logger,
		workers:      make([]basil.Worker, maxWorkers),
		workerPool:   make(chan chan basil.Job, maxWorkers),
		maxWorkers:   maxWorkers,
		maxQueueSize: maxQueueSize,
		jobQueue:     make(chan basil.Job, maxQueueSize),
		quit:         make(chan bool),
	}
}

// Start creates and starts the workers
func (s *Scheduler) Start() {
	s.logger.Debug().
		Int("workers", s.maxWorkers).
		Int("queueSize", s.maxQueueSize).
		Msg("job scheduler starting")

	for i := 0; i < s.maxWorkers; i++ {
		s.workers[i] = NewWorker(s.logger, s.workerPool)
		s.workers[i].Start()
	}

	go s.dispatch()
}

// Stop stops all the workers and the dispatcher process
func (s *Scheduler) Stop() {
	go func() {
		s.quit <- true
	}()

	for i := 0; i < s.maxWorkers; i++ {
		s.workers[i].Stop()
	}

	s.logger.Debug().Msg("job scheduler stopped")
}

// Schedule schedules a new job
func (s *Scheduler) Schedule(job basil.Job) {
	if job.Lightweight() {
		s.logger.Debug().
			ID("jobID", job.JobID()).
			Msg("job running")
		go job.Run()
		return
	}
	s.jobQueue <- job
}

func (s *Scheduler) dispatch() {
	for {
		select {
		case job := <-s.jobQueue:
			go func(job basil.Job) {
				jobChannel := <-s.workerPool
				jobChannel <- job
			}(job)
		case <-s.quit:
			return
		}
	}
}
