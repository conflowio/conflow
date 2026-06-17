---
title: Workflow engine improvements
summary: Design for hardening and evolving Conflow's dependency-driven runtime orchestration.
parent: ../index.md
keywords: [workflow, runtime, concurrency, orchestration]
---

# Workflow engine improvements

**Status:** Draft  
**Date:** 2026-06-17

## Problem

Conflow's workflow engine executes block graphs concurrently using a dependency-driven, stage-based model (`init` → `main` → `close`). The design is sound for a DSL runtime, but a code audit identified gaps between **documented behaviour**, **actual concurrency semantics**, and **operational needs** (cancellation, observability, graceful failure).

Left unaddressed, large parallel workflows risk goroutine pressure, hung scheduling, incorrect timeout/cancellation propagation, and difficult production debugging.

## Goal

Document the current architecture, enumerate improvement opportunities with recommended directions, and define an **independent, mergeable batch roadmap** so each batch can ship with tests without waiting on later work.

## Non-goals (this initiative)

- Changing Conflow language syntax or block author APIs (except optional new lifecycle interfaces).
- Distributed workflow durability (Temporal-style persistence).
- Replacing the dependency resolver or pub/sub model wholesale.
- Performance benchmarking / optimisation beyond concurrency bounds (separate initiative if needed).

## What is already strong

- **Dependency resolution** with generator-aware graph splitting (`pkg/conflow/dependency/resolver.go`).
- **Stage-based lazy evaluation** — parameters and children evaluate only when their stage is active and deps are ready.
- **Generator backpressure** via `PublishBlock` + wait groups.
- **Error path mapping** via `TransformPathErrors` back to Conflow source locations.
- **Clean separation** between parse (`eval.go`), graph (`dependency`), wiring (`node_container.go`), and execution (`block/container.go`).

## Document map

| Document | Contents |
|----------|----------|
| [Architecture (baseline)](./architecture.md) | Current runtime components, data flow, and key source files |
| [Improvements](./improvements/index.md) | Prioritized improvement areas grouped by concern |
| [Batch roadmap](./batch-roadmap.md) | Implementation sequence; what each batch includes and excludes |

## Priority summary

| Priority | Area | Document |
|----------|------|----------|
| P1 | Lightweight goroutine model vs worker pool | [Concurrency & scheduling](./improvements/concurrency.md) |
| P1 | Context propagation for per-node timeouts | [Context & lifecycle](./improvements/context-and-lifecycle.md) |
| P1 | Scheduler queue blocking | [Concurrency & scheduling](./improvements/concurrency.md) |
| P2 | Lifecycle hooks on skip/abort/error | [Context & lifecycle](./improvements/context-and-lifecycle.md) |
| P2 | FSM clarity and extensibility | [Orchestration](./improvements/orchestration.md) |
| P2 | Observability (tracing, metrics) | [Reliability & operations](./improvements/reliability-and-ops.md) |
| P3 | Async channel delivery, pub/sub panics, retry centralisation | [Orchestration](./improvements/orchestration.md), [Reliability & operations](./improvements/reliability-and-ops.md) |
| P3 | Integration / race tests | [Testing strategy](./improvements/testing.md) |
| P3 | Static vs dynamic path unification | [Orchestration](./improvements/orchestration.md) |

## Open questions (resolve before batch 1)

1. **Concurrency ceiling:** Should lightweight jobs remain goroutine-per-job with a global semaphore, or should most work go through the bounded worker pool?
2. **Lifecycle API:** New optional interfaces (`OnSkip`, `OnAbort`) vs always calling `Close` on terminal paths?
3. **Observability backend:** OpenTelemetry only, or also a minimal metrics interface for embedders without OTel?

## References

- Product: [Evaluation pipeline](../../product/runtime/evaluation-pipeline.md), [Job scheduler](../../product/runtime/job-scheduler.md), [Dependencies and order](../../product/concepts/dependencies-and-order.md)
- Source: `pkg/conflow/block/container.go`, `pkg/conflow/node_container.go`, `pkg/conflow/job/`, `pkg/conflow/pubsub.go`, `pkg/conflow/dependency/resolver.go`
