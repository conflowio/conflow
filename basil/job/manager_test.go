// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package job_test

import (
	"io/ioutil"
	"sync/atomic"
	"time"

	"github.com/opsidian/basil/logger"
	"github.com/rs/zerolog"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/basilfakes"
	"github.com/opsidian/basil/basil/job"
)

var _ = Describe("Manager", func() {
	var manager *job.Manager
	var scheduler *basilfakes.FakeScheduler

	BeforeEach(func() {
		scheduler = &basilfakes.FakeScheduler{}
		l := logger.NewZeroLogLogger(zerolog.New(zerolog.ConsoleWriter{Out: ioutil.Discard}))
		manager = job.NewManager("test_manager", scheduler, l)
	})

	AfterEach(func() {
		manager.Stop()
	})

	When("a job is pending", func() {
		BeforeEach(func() {
			manager.Pending("job_id")
		})

		It("should increase the pending count", func() {
			Expect(manager.PendingJobCount()).To(Equal(1))
		})

		It("should increase the remaining count", func() {
			Expect(manager.RemainingJobCount()).To(Equal(1))
		})

		When("scheduled", func() {
			BeforeEach(func() {
				job := &basilfakes.FakeJob{}
				job.IDReturns("job_id")
				manager.Schedule(job)
			})

			It("should decrease the pending count", func() {
				Expect(manager.PendingJobCount()).To(Equal(0))
			})

			It("should increase the running count", func() {
				Expect(manager.RunningJobCount()).To(Equal(1))
			})

			It("should keep the remaining count", func() {
				Expect(manager.RemainingJobCount()).To(Equal(1))
			})
		})

		When("a second job is pending", func() {
			It("should further increase the pending count", func() {
				manager.Pending("job_id_2")
				Expect(manager.PendingJobCount()).To(Equal(2))
			})
		})
	})

	When("a job is scheduled", func() {
		var job *basilfakes.FakeJob

		BeforeEach(func() {
			job = &basilfakes.FakeJob{}
			job.IDReturns("job_id")
			manager.Schedule(job)
		})

		It("should call the scheduler", func() {
			Expect(scheduler.ScheduleCallCount()).To(Equal(1))
			passedJob := scheduler.ScheduleArgsForCall(0)
			Expect(passedJob).To(Equal(job))
		})

		It("should increase the running count", func() {
			Expect(manager.RunningJobCount()).To(Equal(1))
		})

		When("a second job is scheduled", func() {
			It("should further increase the running count", func() {
				job2 := &basilfakes.FakeJob{}
				job2.IDReturns("job_id_2")
				manager.Schedule(job2)
				Expect(manager.RunningJobCount()).To(Equal(2))
			})
		})

		When("finished", func() {
			BeforeEach(func() {
				manager.Finished("job_id")
			})
			It("should decrease the running count", func() {
				Expect(manager.RunningJobCount()).To(Equal(0))
			})
		})

		When("failed with no retry", func() {
			BeforeEach(func() {
				manager.Finished("job_id")
			})

			It("should decrease the running count", func() {
				Expect(manager.RunningJobCount()).To(Equal(0))
			})

			It("should not increase the pending count", func() {
				Expect(manager.PendingJobCount()).To(Equal(0))
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

			It("should increase the pending count", func() {
				Expect(manager.PendingJobCount()).To(Equal(1))
			})

			It("should retry", func() {
				Expect(retried1).To(BeTrue())
				Eventually(func() int64 { return atomic.LoadInt64(&tries) }).Should(Equal(int64(1)))
			})

			When("the job is scheduled again", func() {
				BeforeEach(func() {
					manager.Schedule(job)
				})

				It("should decrease the pending count", func() {
					Expect(manager.PendingJobCount()).To(Equal(0))
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
						Expect(manager.PendingJobCount()).To(Equal(0))
					})
				})
			})
		})
	})

	When("stopped", func() {
		var remaining int

		JustBeforeEach(func() {
			remaining = manager.Stop()
		})

		It("should return with 0 jobs remaining", func() {
			Expect(remaining).To(Equal(0))
		})

		When("there is a job scheduled but not running", func() {
			var job *basilfakes.FakeJob

			BeforeEach(func() {
				job = &basilfakes.FakeJob{}
				job.IDReturns("job_id")
				job.CancelReturns(true)
				manager.Schedule(job)
			})

			It("should successfully cancel it", func() {
				Expect(job.CancelCallCount()).To(Equal(1))
				Expect(remaining).To(Equal(0))
			})
		})

		When("there is a job running", func() {
			var job *basilfakes.FakeJob

			BeforeEach(func() {
				job = &basilfakes.FakeJob{}
				job.IDReturns("job_id")
				job.CancelReturns(false)
				manager.Schedule(job)
			})

			It("should not decrease the running count", func() {
				Expect(job.CancelCallCount()).To(Equal(1))
				Expect(remaining).To(Equal(1))
			})
		})

		When("there is a pending job", func() {
			BeforeEach(func() {
				manager.Pending("job_id")
			})

			It("should remove all pending jobs", func() {
				Expect(manager.PendingJobCount()).To(Equal(0))
			})

		})

		When("there is a job that will fail", func() {
			BeforeEach(func() {
				job := &basilfakes.FakeJob{}
				job.IDReturns("job_id")
				manager.Schedule(job)
			})

			When("the retry timer expired", func() {
				var retried chan struct{}

				BeforeEach(func() {
					retried = make(chan struct{})
					manager.Retry("job_id", 2, 0, func(j basil.Job) func() {
						return func() {
							retried <- struct{}{}
						}
					})
					<-retried
				})

				It("should not decrease the pending count", func() {
					Expect(manager.PendingJobCount()).To(Equal(1))
				})
			})

			When("the retry timer is active", func() {
				BeforeEach(func() {
					manager.Retry("job_id", 2, 1*time.Minute, func(j basil.Job) func() {
						return func() {}
					})
				})

				It("should decrease the pending count", func() {
					Expect(manager.PendingJobCount()).To(Equal(0))
				})
			})
		})

		When("pending is called", func() {
			JustBeforeEach(func() {
				manager.Pending("test_id")
			})

			It("should not increase the pending count", func() {
				Expect(manager.PendingJobCount()).To(Equal(0))
			})
		})

		When("schedule is called", func() {
			JustBeforeEach(func() {
				manager.Schedule(&basilfakes.FakeJob{})
			})

			It("should not schedule the job", func() {
				Expect(scheduler.ScheduleCallCount()).To(Equal(0))
			})

			It("should not increase the running count", func() {
				Expect(manager.PendingJobCount()).To(Equal(0))
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

	When("failedRetry is called for an unknown job", func() {
		It("should panic", func() {
			Expect(func() { manager.Retry("non_existing", 1, 0, nil) }).To(Panic())
		})
	})
})
