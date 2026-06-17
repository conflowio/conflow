---
title: Workflow engine improvement areas
summary: Index of prioritized runtime improvements grouped by concern.
parent: ../index.md
keywords: [improvements, concurrency, orchestration, testing]
---

# Improvement areas

Each document below describes one or more related gaps found in the runtime audit, with **current behaviour**, **options**, **recommended direction**, and **affected files**.

## By priority

| ID | Title | Priority | Document |
|----|-------|----------|----------|
| I1 | Lightweight goroutine model | P1 | [Concurrency & scheduling](./concurrency.md#i1-lightweight-goroutine-model) |
| I2 | Scheduler queue blocking | P1 | [Concurrency & scheduling](./concurrency.md#i2-scheduler-queue-blocking) |
| I3 | Context propagation | P1 | [Context & lifecycle](./context-and-lifecycle.md#i3-context-propagation) |
| I4 | Block lifecycle notifications | P2 | [Context & lifecycle](./context-and-lifecycle.md#i4-lifecycle-notifications) |
| I5 | Graceful shutdown | P2 | [Context & lifecycle](./context-and-lifecycle.md#i5-graceful-shutdown) |
| I6 | State machine opacity | P2 | [Orchestration](./orchestration.md#i6-state-machine) |
| I7 | Async channel delivery | P3 | [Orchestration](./orchestration.md#i7-async-channel-delivery) |
| I8 | Deadlock detection scope | P3 | [Orchestration](./orchestration.md#i8-deadlock-detection) |
| I9 | Static vs dynamic paths | P3 | [Orchestration](./orchestration.md#i9-static-vs-dynamic) |
| I10 | PubSub panic on unsubscribe | P3 | [Reliability & operations](./reliability-and-ops.md#i10-pubsub-errors) |
| I11 | Observability gaps | P2 | [Reliability & operations](./reliability-and-ops.md#i11-observability) |
| I12 | Retry policy scatter | P3 | [Reliability & operations](./reliability-and-ops.md#i12-retry-policy) |
| I13 | Test coverage skew | P2 | [Testing strategy](./testing.md) |

## By document

| Document | Improvement IDs |
|----------|-----------------|
| [Concurrency & scheduling](./concurrency.md) | I1, I2 |
| [Context & lifecycle](./context-and-lifecycle.md) | I3, I4, I5 |
| [Orchestration](./orchestration.md) | I6, I7, I8, I9 |
| [Reliability & operations](./reliability-and-ops.md) | I10, I11, I12 |
| [Testing strategy](./testing.md) | I13 |

Implementation order is defined in [Batch roadmap](../batch-roadmap.md).
