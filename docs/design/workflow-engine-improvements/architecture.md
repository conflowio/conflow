---
title: Workflow engine architecture (baseline)
summary: Current runtime components and evaluation data flow before improvements.
parent: index.md
keywords: [architecture, container, scheduler, pubsub]
---

# Workflow engine architecture (baseline)

This document describes the **current** orchestration layer as of the 2026-06-17 audit. Improvement proposals reference these components.

## End-to-end flow

```text
Evaluate (eval.go)
    │
    ├── dependency.Resolver — Tarjan SCC, generator start/finish split, cycle detection
    │
    └── block.Container.Run — per-block state machine
            │
            ├── NodeContainer (per child) — subscribe via PubSub, schedule when deps ready
            ├── job.Tracker — per-block job accounting, retries, pending/running counts
            └── job.Scheduler — worker pool + lightweight goroutine path
```

## Core components

### Entry: `pkg/conflow/eval.go`

`Evaluate()` validates `@input` parameters, creates an `EvalContext`, and calls `parsley.EvaluateNode` on the root block node. Errors are mapped to Conflow paths via `TransformPathErrors`.

### Dependency graph: `pkg/conflow/dependency/resolver.go`

- Builds a graph from workflow nodes.
- Generator blocks get a synthetic `-start` node so dependents on generated children do not create false cycles with the generator's own fields.
- Tarjan's SCC detects circular dependencies at parse/resolve time.

### Dependency wiring: `pkg/conflow/node_container.go`

- Each child node (parameter or block) gets a `NodeContainer`.
- Tracks `missingDeps` / `nilDeps`; subscribes to dependency IDs via `EvalContext.Subscribe`.
- When a dependency publishes, `PubSub.Publish` → `SetDependency` may schedule the node.
- `Run()` on the parent block evaluates all children whose `EvalStage` matches the current parent stage.

### Block orchestration: `pkg/conflow/block/container.go`

State machine (numeric `int64` states, advanced via `containerStateNext` increment):

```text
Start → Init → PreMain → Main → PreClose → Close → Finished
                  │                    │
            Skipped / Errored / Aborted (terminal)
```

Responsibilities:

- Create child `NodeContainer` instances at `Start`.
- At each pre-stage (`PreMain`, `PreClose`), set `evalStage` and call `evaluateChildren()`.
- Schedule `Init` / `Main` / `Close` via `containerStage` jobs on the block's `job.Tracker`.
- Receive completed children on `resultChan`; bind values via `setChild`; publish to dependents.
- `PublishBlock` for generators: schedule generated child, optional callback, wait on `WaitGroup` for backpressure.

### Pub/sub: `pkg/conflow/pubsub.go`

- Per-`EvalContext` subscription lists keyed by dependency ID.
- `Publish(container)` notifies all subscribers synchronously under read lock.
- `Unsubscribe` panics if the subscription is not found.

### Job scheduling: `pkg/conflow/job/`

| Type | Role |
|------|------|
| `Scheduler` | Bounded worker pool (`jobQueue`) + stopped detection |
| `Tracker` | Per-block wrapper: schedule, succeed/fail/cancel, retry with backoff, pending/running counts |
| `containerStage` | Runs one lifecycle stage (`init`/`main`/`close`); handles panic recovery and retry |

### Static path: `pkg/conflow/block/static_container.go`

Synchronous evaluation for blocks without dynamic child scheduling. Uses `context.Background()` for child contexts. Separate code path from async `Container`.

## Lightweight vs worker-pool jobs

`Job.Lightweight()` controls scheduling:

| Job kind | `Lightweight()` | Actual behaviour |
|----------|-----------------|------------------|
| `block.Container` | always `true` | New goroutine via scheduler |
| `parameter.Container` | always `true` | New goroutine via scheduler |
| `containerStage` init/close | `false` | Worker pool queue |
| `containerStage` main | `false` only if block has generated children | Worker pool for generators; goroutine otherwise |

**Note:** Product docs describe lightweight jobs as running "inline on the scheduling goroutine"; the implementation spawns a goroutine for all lightweight jobs (`scheduler.go`).

## Context hierarchy

- Root `EvalContext` wraps the caller's `context.Context` with cancel.
- Each `NodeContainer.createEvalContext` creates a **new** context from `context.Background()` (with optional `@timeout`), not derived from the parent workflow context.
- `EvalContext.Run()` / `Cancel()` use an atomic semaphore so each eval context runs at most once.

## Concurrency patterns in `Container`

- `mainLoop` selects on `stateChan`, `resultChan`, `errChan`, and `evalCtx.Done()`.
- `SetChild` and `SetError` send to channels via **new goroutines** (non-blocking to callers).
- `resultChan` buffer: 8; `errChan` buffer: 1.
- Stage advancement when `jobTracker.ActiveJobCount() == 0` after a child completes.

## Key files reference

| File | Lines of interest |
|------|-------------------|
| `pkg/conflow/block/container.go` | State machine, `mainLoop`, `evaluateChildren`, `PublishBlock` |
| `pkg/conflow/block/container_stage.go` | Stage execution, retry, panic recovery |
| `pkg/conflow/node_container.go` | `createEvalContext`, dependency scheduling |
| `pkg/conflow/job/scheduler.go` | Lightweight vs queued scheduling |
| `pkg/conflow/job/tracker.go` | Retry backoff, active job counts |
| `pkg/conflow/pubsub.go` | Subscribe / publish / unsubscribe |
| `pkg/conflow/eval_context.go` | Context tree, pub/sub ownership |
