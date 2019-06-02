package basil

// Scheduler is the job scheduler
type Scheduler interface {
	Start()
	Stop()
	Schedule(Job)
}

// Worker is an interface for a job queue processor
type Worker interface {
	Start()
	Stop()
}
