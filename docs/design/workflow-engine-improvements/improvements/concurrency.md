---
title: Concurrency and scheduling improvements
summary: Goroutine model, worker pool usage, and scheduler queue behaviour.
parent: index.md
keywords: [concurrency, scheduler, lightweight, goroutines]
---

# Concurrency & scheduling

## I1: Lightweight goroutine model

### Current behaviour

Almost all runtime jobs return `Lightweight() == true`:

- `block.Container` — always `true` (`container.go`)
- `parameter.Container` — always `true` (`parameter/container.go`)
- `containerStage` — `false` for init/close; main is `false` only when the block has generated children

The scheduler runs lightweight jobs in a **new goroutine** per job, bypassing the bounded worker queue:

```go
// pkg/conflow/job/scheduler.go
if job.Lightweight() {
    go func() { job.Run() }()
    return nil
}
```

Under a wide parallel dependency graph, concurrency is **unbounded goroutines**, not `NumCPU*2` workers as examples suggest (`examples/common/main.go`).

Product documentation states lightweight jobs run "inline on the scheduling goroutine" — this no longer matches the code.

### Risk

- Goroutine storms on large workflows (memory, scheduler overhead, harder cancellation).
- Worker pool is largely unused; tuning `maxWorkers` has little effect for typical graphs.
- Misleading docs cause incorrect capacity planning.

### Options

| Option | Description | Trade-offs |
|--------|-------------|------------|
| **A. Global semaphore** | Keep goroutine-per-job but cap concurrent lightweight jobs via semaphore on `Scheduler` | Minimal API change; embedders configure max concurrency; still one goroutine allocated per job |
| **B. Pool-first** | Default `Lightweight() == false` for block/parameter containers; use worker pool | Bounded concurrency; may increase latency for tiny graphs; requires queue-full handling (I2) |
| **C. Tiered** | Small/sync params stay lightweight; block `Run()` and generator work use pool | Balanced; more complex classification rules |

### Recommendation

**Option A for batch 1** (lowest risk): add optional `MaxLightweightConcurrency` on `Scheduler`, defaulting to `runtime.NumCPU()*2` to match examples. Revisit **Option C** if profiling shows parameter eval dominates.

Also **fix product docs** to describe actual goroutine behaviour.

### Affected files

- `pkg/conflow/job/scheduler.go`
- `pkg/conflow/block/container.go` (`Lightweight`)
- `pkg/conflow/parameter/container.go` (`Lightweight`)
- `docs/product/runtime/job-scheduler.md`
- Tests: new `pkg/conflow/job/scheduler_test.go` (concurrency cap)

---

## I2: Scheduler queue blocking

### Current behaviour

Non-lightweight jobs use a blocking send on `jobQueue`:

```go
select {
case s.jobQueue <- job:
    return nil
case <-s.stoppedChan:
    return errors.New("job scheduler was stopped")
}
```

If the queue is full and the scheduler is not stopped, `ScheduleJob` **blocks indefinitely**. No timeout, no "queue full" error, no backpressure to the container.

### Risk

- Hung workflows when queue saturates.
- `Container.mainLoop` blocked indirectly if scheduling happens on its goroutine (depends on call path).
- No metric or log when saturation occurs.

### Options

| Option | Description | Trade-offs |
|--------|-------------|------------|
| **A. Non-blocking + error** | `select` with `default`; return `ErrQueueFull` | Caller must retry or fail workflow; explicit failure mode |
| **B. Context-aware schedule** | `ScheduleJob(ctx, job)`; block until ctx done or slot available | Requires interface change on `JobScheduler` |
| **C. Unbounded queue** | Remove max queue size | Avoids blocking; memory risk under overload |

### Recommendation

**Option A** alongside I1 semaphore work: return a typed error (`job.ErrQueueFull`) and propagate to `Container.setError` with block position. Document queue sizing guidance for embedders.

If `JobScheduler` interface change is acceptable later, add **Option B** as a follow-up batch without breaking existing embedders (wrapper adapter).

### Affected files

- `pkg/conflow/job/scheduler.go`
- `pkg/conflow/job.go` (optional error type)
- `pkg/conflow/block/container.go` (error propagation from `evaluateChildren` / stage schedule)
- `docs/product/runtime/job-scheduler.md`
