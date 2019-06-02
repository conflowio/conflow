package job

import (
	"github.com/opsidian/basil/basil"
)

// Scheduler handles workers and schedules jobs
type Scheduler struct {
	maxWorkers int
	workers    []basil.Worker
	workerPool chan chan basil.Job
	jobQueue   chan basil.Job
	quit       chan bool
}

// NewScheduler creates a new scheduler instance
func NewScheduler(maxWorkers int, maxQueueSize int) *Scheduler {
	return &Scheduler{
		workers:    make([]basil.Worker, maxWorkers),
		workerPool: make(chan chan basil.Job, maxWorkers),
		maxWorkers: maxWorkers,
		jobQueue:   make(chan basil.Job, maxQueueSize),
		quit:       make(chan bool),
	}
}

// Start creates and starts the workers
func (s *Scheduler) Start() {
	for i := 0; i < s.maxWorkers; i++ {
		s.workers[i] = NewWorker(s.workerPool)
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
}

// Schedule schedules a new job
func (s *Scheduler) Schedule(job basil.Job) {
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
