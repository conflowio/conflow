// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package job_test

import (
	"sync/atomic"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/basilfakes"
	"github.com/opsidian/basil/basil/job"
	"github.com/opsidian/basil/logger/zerolog"
)

var _ = Describe("Manager", func() {
	var manager *job.Manager
	var scheduler *basilfakes.FakeScheduler

	BeforeEach(func() {
		scheduler = &basilfakes.FakeScheduler{}
		logger := zerolog.NewDisabledLogger()
		manager = job.NewManager("test_manager", scheduler, logger)
	})

	AfterEach(func() {
		manager.Stop()
	})

	When("a job is scheduled", func() {
		var job *basilfakes.FakeJob

		BeforeEach(func() {
			job = &basilfakes.FakeJob{}
			job.JobIDReturns("job_id")
			manager.Schedule(job)
		})

		It("should call the scheduler", func() {
			Expect(scheduler.ScheduleCallCount()).To(Equal(1))
			passedJob := scheduler.ScheduleArgsForCall(0)
			Expect(passedJob).To(Equal(job))
		})

		It("should increase the running job count", func() {
			Expect(manager.RunningJobCount()).To(Equal(1))
		})

		It("should increase the active job count", func() {
			Expect(manager.ActiveJobCount()).To(Equal(1))
		})

		When("a second job is scheduled", func() {
			It("should further increase the running count", func() {
				job2 := &basilfakes.FakeJob{}
				job2.JobIDReturns("job_id_2")
				manager.Schedule(job2)
				Expect(manager.RunningJobCount()).To(Equal(2))
			})
		})

		When("finished", func() {
			BeforeEach(func() {
				manager.Finished("job_id")
			})
			It("should decrease the running job count", func() {
				Expect(manager.RunningJobCount()).To(Equal(0))
			})
		})

		When("failed with no retry", func() {
			BeforeEach(func() {
				manager.Finished("job_id")
			})

			It("should decrease the running job count", func() {
				Expect(manager.RunningJobCount()).To(Equal(0))
			})
		})

		When("failed with a retry", func() {
			var tries int64
			var retried1 bool
			BeforeEach(func() {
				retried1 = manager.Retry("job_id", 2, 1*time.Millisecond, func(j basil.Job) func() {
					return func() {
						atomic.AddInt64(&tries, 1)
					}
				})
			})

			It("should retry", func() {
				Expect(retried1).To(BeTrue())
				Eventually(func() int64 { return atomic.LoadInt64(&tries) }).Should(Equal(int64(1)))
			})

			When("the job is scheduled again", func() {
				BeforeEach(func() {
					manager.Schedule(job)
				})

				When("failing the second time", func() {
					var retried2 bool
					BeforeEach(func() {
						retried2 = manager.Retry("job_id", 2, 1*time.Millisecond, func(j basil.Job) func() {
							return func() {
								atomic.AddInt64(&tries, 1)
							}
						})
					})

					It("should not retry", func() {
						Expect(retried2).To(BeFalse())
					})
				})
			})
		})
	})

	When("stopped", func() {
		var active int

		JustBeforeEach(func() {
			active = manager.Stop()
		})

		It("should return with 0 jobs active", func() {
			Expect(active).To(Equal(0))
		})

		When("there is a job scheduled but not running", func() {
			var job *basilfakes.FakeJob

			BeforeEach(func() {
				job = &basilfakes.FakeJob{}
				job.JobIDReturns("job_id")
				job.CancelReturns(true)
				manager.Schedule(job)
			})

			It("should successfully cancel it", func() {
				Expect(job.CancelCallCount()).To(Equal(1))
				Expect(active).To(Equal(0))
			})
		})

		When("there is a job running", func() {
			var job *basilfakes.FakeJob

			BeforeEach(func() {
				job = &basilfakes.FakeJob{}
				job.JobIDReturns("job_id")
				job.CancelReturns(false)
				manager.Schedule(job)
			})

			It("should not decrease the active job count", func() {
				Expect(job.CancelCallCount()).To(Equal(1))
				Expect(active).To(Equal(1))
			})
		})

		When("schedule is called", func() {
			JustBeforeEach(func() {
				manager.Schedule(&basilfakes.FakeJob{})
			})

			It("should not schedule the job", func() {
				Expect(scheduler.ScheduleCallCount()).To(Equal(0))
			})
		})
	})

	When("finished is called for an unknown job", func() {
		It("should panic", func() {
			Expect(func() { manager.Finished("non_existing") }).To(Panic())
		})
	})

	When("failed is called for an unknown job", func() {
		It("should panic", func() {
			Expect(func() { manager.Finished("non_existing") }).To(Panic())
		})
	})

	When("retry is called for an unknown job", func() {
		It("should panic", func() {
			Expect(func() { manager.Retry("non_existing", 1, 0, nil) }).To(Panic())
		})
	})

	When("pending jobs are added", func() {
		BeforeEach(func() {
			manager.AddPending(2)
		})
		It("should increase the pending jobs count", func() {
			Expect(manager.PendingJobCount()).To(Equal(2))
		})
		It("should increase the active jobs count", func() {
			Expect(manager.ActiveJobCount()).To(Equal(2))
		})
	})

	It("should generate unique job ids", func() {
		id1 := manager.GenerateJobID("job_id")
		id2 := manager.GenerateJobID("job_id")
		Expect(id1).ToNot(Equal(id2))
	})
})
