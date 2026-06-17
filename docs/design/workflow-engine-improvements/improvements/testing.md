---
title: Workflow engine testing strategy
summary: Test gaps and verification approach for runtime improvements.
parent: index.md
keywords: [testing, race detector, integration, concurrency]
---

# Testing strategy

## I13: Test coverage skew

### Current state

| Area | Coverage | Gap |
|------|----------|-----|
| `block/container_test.go` | Ginkgo + counterfeiter; init/error/retry paths | No multi-child parallel scenarios |
| `node_container_test.go` | Dependency wiring, directives, triggers | No real scheduler integration |
| `job/tracker_test.go` | Schedule, retry, stop | Scheduler queue saturation untested |
| `dependency/resolver_test.go` | Cycles, generators | Adequate for static graph |
| Integration | Examples as manual runs | No automated parallel workflow fixtures |

No `go test -race` job identified in CI for orchestration packages.

### Goal

Each implementation batch adds tests that fail without the fix and pass with it. Batch 7 adds cross-cutting integration tests.

## Test categories to add

### Unit (per batch)

- **Concurrency (batch 1–2):** semaphore cap enforced; `ErrQueueFull` returned when queue saturated.
- **Context (batch 3):** root cancel stops child eval; block `@timeout` respects parent deadline when shorter.
- **Lifecycle (batch 4):** `Close` called on error path; skip behaviour documented and tested.
- **PubSub (batch 5):** unsubscribe mismatch returns error, no panic.

### Integration fixtures (batch 7)

Deterministic workflows under `pkg/conflow/test/fixtures/workflows/` (or `test/integration/`):

| Fixture | Validates |
|---------|-----------|
| `fan_out_fan_in` | Parallel blocks, join on dependency |
| `generator_stream` | `PublishBlock` backpressure, subscriber wait |
| `abort_mid_run` | Root ctx cancel during `Run` |
| `retry_exhaustion` | Finite retry limit surfaces error |
| `stage_skip_init` | Init skip propagates nil deps |

Run with:

```bash
go test -race ./pkg/conflow/... ./pkg/conflow/block/... ./pkg/conflow/job/...
```

### Regression gate

After every batch: `go test ./...` green. Before merging batch 7: mandatory `-race` on packages touched by concurrency/orchestration changes.

## Fake vs real scheduler

Prefer **real** `job.Scheduler` with 1–2 workers in integration tests; keep counterfeiter for isolated container logic unit tests.

## Verification checklist (per batch)

1. New tests fail on baseline SHA without implementation.
2. Tests pass after implementation.
3. No new panics in error paths covered by I10.
4. Product doc updated when user-visible semantics change.

## References

- Skills: `test-driven-development`, `condition-based-waiting` (for async tests), `verification-before-completion`
- Existing: `pkg/conflow/block/container_test.go`, `pkg/conflow/node_container_test.go`
