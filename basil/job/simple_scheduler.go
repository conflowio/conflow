package job

import "github.com/opsidian/basil/basil"

type SimpleScheduler struct{}

func (s SimpleScheduler) ScheduleJob(job basil.Job) error {
	go func() {
		job.Run()
	}()

	return nil
}
