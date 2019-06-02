package test

import "github.com/opsidian/basil/basil"

// Scheduler is a test scheduler, it will simply run the given job in a goroutine in the background
type Scheduler struct{}

func (s Scheduler) Start() {}
func (s Scheduler) Stop()  {}
func (s Scheduler) Schedule(job basil.Job) {
	go job.Run()
}
