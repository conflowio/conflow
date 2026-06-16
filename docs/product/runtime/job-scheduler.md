---
title: Job scheduler
summary: Concurrent execution of block jobs via pkg/conflow/job.
parent: runtime/index.md
keywords: [job scheduler, concurrency, parallel]
---

# Job scheduler

## Role

The job scheduler executes block lifecycle work **concurrently** when dependencies allow. It decouples **dependency readiness** from **OS thread assignment**.

## Interfaces

`pkg/conflow/job.go`:

```go
type JobScheduler interface {
    ScheduleJob(job Job) error
}

type Job interface {
    JobName() ID
    JobID() int
    SetJobID(int)
    Run()
    Lightweight() bool
}
```

`JobContainer` adds cancellation and eval stage tracking.

## Default setup (examples)

`examples/common/main.go`:

```go
scheduler := job.NewScheduler(logger, runtime.NumCPU()*2, 100)
scheduler.Start()
defer scheduler.Stop()
```

Worker count scales with CPU; queue depth 100.

## Lightweight jobs

Jobs returning `Lightweight() == true` may run inline on the scheduling goroutine — used for small work to reduce overhead.

## Interaction with generators

`PublishBlock` coordinates with scheduled jobs so dependents finish processing a publication before the generator proceeds — backpressure for streams and iterators.

## Cancellation

Jobs respect context cancellation on `EvalContext`; long-running blocks should use `ctx.Done()` (see `exec` abort path).

## Implementation

Primary package: `pkg/conflow/job/` (scheduler, tracker).

## See also

- [Dependencies and evaluation order](../concepts/dependencies-and-order.md)
- [Generators](../concepts/generators.md)
