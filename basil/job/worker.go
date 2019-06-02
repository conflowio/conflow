package job

import (
	"github.com/opsidian/basil/basil"
)

// Worker is a goroutine processing jobs
type Worker struct {
	workerPool chan chan basil.Job
	jobChannel chan basil.Job
	quit       chan bool
}

// NewWorker creates a new worker instance
func NewWorker(workerPool chan chan basil.Job) Worker {
	return Worker{
		workerPool: workerPool,
		jobChannel: make(chan basil.Job),
		quit:       make(chan bool),
	}
}

// Start starts the worker goroutine
func (w Worker) Start() {
	go func() {
		for {
			w.workerPool <- w.jobChannel
			select {
			case job := <-w.jobChannel:
				job.Run()
			case <-w.quit:
				return
			}
		}
	}()
}

// Stop stops the worker goroutine
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
