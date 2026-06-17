---
title: Workflow engine batch roadmap
summary: Independent, mergeable implementation batches derived from the improvement design.
parent: index.md
keywords: [roadmap, batches, implementation]
---

# Batch roadmap

Each batch is **independently mergeable** with its own tests and plan file in `docs/plans/`. Create one plan file per batch **when starting that batch** (see the writing-plans skill for task structure), not upfront for all batches.

Plans use the header format:

```markdown
**Design Doc:** docs/design/workflow-engine-improvements/index.md
```

## Overview

| Batch | Name | Improvements | Mergeable alone? |
|-------|------|--------------|------------------|
| **1** | Lightweight concurrency cap | I1 (partial), docs fix | Yes |
| **2** | Scheduler queue errors | I2 | Yes |
| **3** | Context propagation | I3 | Yes |
| **4** | Lifecycle hooks | I4, I5 (partial) | Yes |
| **5** | Reliability & ops | I10, I11 (minimal), I12 (partial) | Yes |
| **6** | State machine refactor | I6 | Yes |
| **7** | Integration & race tests | I13, I7 (validation), I8 (debug) | Yes |
| **8** | Static/dynamic unification | I9 | Yes (after 3–4) |

## Batch 1: Lightweight concurrency cap

**Goal:** Bound goroutine concurrency for lightweight jobs without changing block author APIs.

**Includes:**

- Optional `MaxLightweightConcurrency` on `job.Scheduler` (semaphore)
- Default aligned with `runtime.NumCPU()*2`
- Unit tests for cap enforcement
- Fix `docs/product/runtime/job-scheduler.md` lightweight description

**Excludes:**

- Changing `Lightweight()` return values on containers (defer to later profiling)
- `JobScheduler` interface signature change

**Key files:** `pkg/conflow/job/scheduler.go`, new `scheduler_test.go`, product doc.

**Plan file (when ready):** `docs/plans/2026-06-17-workflow-engine-batch-1.md`

---

## Batch 2: Scheduler queue errors

**Goal:** Never block forever when the job queue is full.

**Includes:**

- `job.ErrQueueFull` (or typed sentinel)
- Non-blocking send with `default` branch
- Propagate error to container with block position
- Tests: saturated queue returns error

**Excludes:**

- Context-aware `ScheduleJob(ctx)` (future batch)

**Key files:** `pkg/conflow/job/scheduler.go`, `pkg/conflow/block/container.go`

**Depends on:** Batch 1 recommended (shared scheduler tests) but not strictly required.

**Plan file:** `docs/plans/2026-06-17-workflow-engine-batch-2.md`

---

## Batch 3: Context propagation

**Goal:** Child eval contexts derive from the parent workflow context; `@timeout` constrains rather than replaces.

**Includes:**

- Fix `NodeContainer.createEvalContext`
- Fix `StaticContainer.createContainer` (same semantics)
- Tests: root cancel, nested timeout vs parent deadline
- Document `@timeout` in runtime-directives

**Excludes:**

- Lifecycle hook changes (batch 4)

**Key files:** `pkg/conflow/node_container.go`, `pkg/conflow/block/static_container.go`

**Plan file:** `docs/plans/2026-06-17-workflow-engine-batch-3.md`

---

## Batch 4: Lifecycle hooks

**Goal:** Defined behaviour on skip, error, and abort; improved shutdown error handling.

**Includes:**

- Implement TODOs in `setState` (see [Context & lifecycle](./improvements/context-and-lifecycle.md#i4-lifecycle-notifications))
- `Close` on error/abort when `BlockCloser` implemented
- Configurable shutdown timeout on `RuntimeConfig`
- Tests for each terminal path

**Excludes:**

- New block author interfaces beyond optional `BlockSkipNotifier` if chosen
- Full multi-error aggregation (minimal: primary error + log others)

**Key files:** `pkg/conflow/block/container.go`, `pkg/conflow/block.go`

**Depends on:** Batch 3 (shutdown should respect parent ctx).

**Plan file:** `docs/plans/2026-06-17-workflow-engine-batch-4.md`

---

## Batch 5: Reliability & operations

**Goal:** Replace panics with errors; optional metrics hooks; centralise retry defaults.

**Includes:**

- PubSub / SetDependency error returns (I10)
- Optional `RuntimeMetrics` interface + no-op default (I11 minimal)
- Extract retry backoff to `job/retry_policy.go`; document limits (I12)
- Unit tests for each

**Excludes:**

- OpenTelemetry dependency (optional follow-up)
- Default retry limit behaviour change (document only unless explicitly approved)

**Key files:** `pkg/conflow/pubsub.go`, `pkg/conflow/job/tracker.go`, new retry policy module

**Plan file:** `docs/plans/2026-06-17-workflow-engine-batch-5.md`

---

## Batch 6: State machine refactor

**Goal:** Typed container states and explicit transitions for debuggability.

**Includes:**

- `ContainerState` type with `String()`
- Explicit transition helper (replace `containerStateNext` increment)
- Log state names in debug events
- Refactor tests unchanged behaviour

**Excludes:**

- Table-driven FSM for pause/checkpoint
- Async delivery refactor (I7)

**Key files:** `pkg/conflow/block/container.go`, optional `state.go`

**Plan file:** `docs/plans/2026-06-17-workflow-engine-batch-6.md`

---

## Batch 7: Integration & race tests

**Goal:** Confidence in concurrent behaviour under the race detector.

**Includes:**

- Integration fixtures (fan-out, generator, abort, retry, skip)
- CI or documented command: `go test -race` on orchestration packages
- Debug logging for stuck workflows (I8 minimal)

**Excludes:**

- Full dispatcher refactor (I7 Option C)
- Production heuristic abort on stuck detection

**Key files:** new `test/integration/` or fixtures under `pkg/conflow/`

**Depends on:** Batches 1–4 recommended.

**Plan file:** `docs/plans/2026-06-17-workflow-engine-batch-7.md`

---

## Batch 8: Static/dynamic unification

**Goal:** Single evaluation semantics for static and dynamic containers.

**Includes:**

- Audit `NewStaticContainer` call sites
- Shared internal eval helper with `EvalMode` flag
- Parity tests: static vs dynamic produce same binding for fixture blocks

**Excludes:**

- Removing static container entirely

**Depends on:** Batches 3–4 (context + lifecycle parity).

**Plan file:** `docs/plans/2026-06-17-workflow-engine-batch-8.md`

---

## Verification strategy (all batches)

- Unit tests: bind policy per change; aliasing/race where relevant.
- Regression: `go test ./...` after each batch.
- From batch 7 onward: `go test -race` on `pkg/conflow/block`, `pkg/conflow/job`, root `pkg/conflow` tests.

## Creating batch plans

When the user asks to implement a batch:

1. Read this design doc and the relevant improvement leaf doc.
2. Use the writing-plans skill to create **one** plan file for that batch only.
3. Include `## Review` section per writing-plans skill.
4. Do not combine multiple batches into one plan file.
