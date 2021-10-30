// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package job_test

import (
	"sync/atomic"
	"time"

	"github.com/opsidian/conflow/conflow"
	"github.com/opsidian/conflow/conflow/job/jobfakes"

	"github.com/opsidian/conflow/conflow/job"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/opsidian/conflow/conflow/conflowfakes"
	"github.com/opsidian/conflow/loggers/zerolog"
)

//counterfeiter:generate . pendingJob
type pendingJob interface {
	conflow.Job
	job.Pending
}

//counterfeiter:generate . cancellableJob
type cancellableJob interface {
	conflow.Job
	job.Cancellable
}

var _ = Describe("Scheduler", func() {
	var tracker *job.Tracker
	var scheduler *conflowfakes.FakeJobScheduler

	BeforeEach(func() {
		scheduler = &conflowfakes.FakeJobScheduler{}
		logger := zerolog.NewDisabledLogger()
		backoff := job.ExponentialRetryBackoff(1.5, 1*time.Millisecond, 1*time.Second)
		tracker = job.NewTracker("test_tracker", scheduler, logger, backoff)
	})

	AfterEach(func() {
		tracker.Stop()
	})

	When("a job is scheduled", func() {
		var job *conflowfakes.FakeJob

		BeforeEach(func() {
			job = &conflowfakes.FakeJob{}
			job.JobIDReturns(1)
			tracker.ScheduleJob(job)
		})

		It("should call the scheduler", func() {
			Expect(scheduler.ScheduleJobCallCount()).To(Equal(1))
			passedJob := scheduler.ScheduleJobArgsForCall(0)
			Expect(passedJob).To(Equal(job))
		})

		It("should increase the running job count", func() {
			Expect(tracker.RunningJobCount()).To(Equal(1))
		})

		It("should increase the active job count", func() {
			Expect(tracker.ActiveJobCount()).To(Equal(1))
		})

		When("a second job is scheduled", func() {
			It("should further increase the running count", func() {
				job2 := &conflowfakes.FakeJob{}
				job2.JobIDReturns(1)
				tracker.ScheduleJob(job2)
				Expect(tracker.RunningJobCount()).To(Equal(2))
			})
		})

		When("finished", func() {
			BeforeEach(func() {
				tracker.Succeeded(1)
			})
			It("should decrease the running job count", func() {
				Expect(tracker.RunningJobCount()).To(Equal(0))
			})
		})

		When("failed with no retry", func() {
			BeforeEach(func() {
				tracker.Succeeded(1)
			})

			It("should decrease the running job count", func() {
				Expect(tracker.RunningJobCount()).To(Equal(0))
			})
		})

		When("failed with a retry", func() {
			var tries int64
			var retried1 bool
			BeforeEach(func() {
				retried1 = tracker.Retry(1, 1, 1*time.Millisecond, "test", func() {
					atomic.AddInt64(&tries, 1)
				})
			})

			It("should retry", func() {
				Expect(retried1).To(BeTrue())
				Eventually(func() int64 { return atomic.LoadInt64(&tries) }).Should(Equal(int64(1)))
			})

			When("the job is scheduled again", func() {
				BeforeEach(func() {
					tracker.ScheduleJob(job)
				})

				When("failing the second time", func() {
					var retried2 bool
					BeforeEach(func() {
						retried2 = tracker.Retry(1, 1, 1*time.Millisecond, "test", func() {
							atomic.AddInt64(&tries, 1)
						})
					})

					PIt("should not retry", func() {
						Expect(retried2).To(BeFalse())
					})
				})
			})
		})
	})

	When("a pending job is scheduled", func() {
		BeforeEach(func() {
			tracker.AddPending(2)
			job := &jobfakes.FakePendingJob{}
			job.JobIDReturns(1)
			job.PendingReturns(true)
			tracker.ScheduleJob(job)
		})
		It("should decrease the pending count", func() {
			Expect(tracker.PendingJobCount()).To(Equal(1))
		})
	})

	When("stopped", func() {
		var active int

		JustBeforeEach(func() {
			active = tracker.Stop()
		})

		It("should return with 0 jobs active", func() {
			Expect(active).To(Equal(0))
		})

		When("there is a job scheduled but not running", func() {
			var job *jobfakes.FakeCancellableJob

			BeforeEach(func() {
				job = &jobfakes.FakeCancellableJob{}
				job.JobIDReturns(1)
				job.CancelReturns(true)
				tracker.ScheduleJob(job)
			})

			It("should successfully cancel it", func() {
				Expect(job.CancelCallCount()).To(Equal(1))
				Expect(active).To(Equal(0))
			})
		})

		When("there is a job running", func() {
			var job *jobfakes.FakeCancellableJob

			BeforeEach(func() {
				job = &jobfakes.FakeCancellableJob{}
				job.JobIDReturns(1)
				job.CancelReturns(false)
				tracker.ScheduleJob(job)
			})

			It("should not decrease the active job count", func() {
				Expect(job.CancelCallCount()).To(Equal(1))
				Expect(active).To(Equal(1))
			})
		})

		When("schedule is called", func() {
			JustBeforeEach(func() {
				tracker.ScheduleJob(&conflowfakes.FakeJob{})
			})

			It("should not schedule the job", func() {
				Expect(scheduler.ScheduleJobCallCount()).To(Equal(0))
			})
		})
	})

	When("finished is called for an unknown job", func() {
		It("should panic", func() {
			Expect(func() { tracker.Succeeded(999) }).To(Panic())
		})
	})

	When("failed is called for an unknown job", func() {
		It("should panic", func() {
			Expect(func() { tracker.Succeeded(999) }).To(Panic())
		})
	})

	When("retry is called for an unknown job", func() {
		It("should panic", func() {
			Expect(func() { tracker.Retry(999, 1, 0, "test", nil) }).To(Panic())
		})
	})

	When("pending jobs are added", func() {
		BeforeEach(func() {
			tracker.AddPending(2)
		})
		It("should increase the pending jobs count", func() {
			Expect(tracker.PendingJobCount()).To(Equal(2))
		})
		It("should increase the active jobs count", func() {
			Expect(tracker.ActiveJobCount()).To(Equal(2))
		})
	})
})
