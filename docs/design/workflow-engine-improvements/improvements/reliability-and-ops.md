---
title: Reliability and operations improvements
summary: Error handling, observability, and retry policy centralisation.
parent: index.md
keywords: [reliability, observability, retry, pubsub, errors]
---

# Reliability & operations

## I10: PubSub errors

### Current behaviour

`PubSub.Unsubscribe` panics when the subscription is not found:

```go
panic(fmt.Errorf("unsubscribe unsuccessful, %q was never subscribed for %q", ...))
```

`NodeContainer.SetDependency` panics on unknown dependency IDs.

### Risk

- Programming errors or double-close crash the entire process.
- Inconsistent with the rest of the engine, which returns `parsley.Error` with source positions.

### Recommendation

Return errors through `parent.SetError` instead of panicking. Keep panic only for truly invariant violations in development builds if desired (`//go:build dev`).

### Affected files

- `pkg/conflow/pubsub.go`
- `pkg/conflow/node_container.go` — `SetDependency`
- New tests in `pkg/conflow/pubsub_test.go`

---

## I11: Observability

### Current behaviour

Debug logging exists on scheduler, tracker, and container (zerolog events: block id, state, active_jobs). There is no:

- Distributed tracing / span per block stage
- Metrics (queue depth, active jobs, retry counts, stage latency)
- Structured workflow execution summary on completion

### Risk

Production question "why is this workflow stuck?" requires log archaeology.

### Options

| Option | Description | Trade-offs |
|--------|-------------|------------|
| **A. OpenTelemetry hooks** | Optional tracer on `EvalContext`; spans at stage boundaries | Standard; optional dependency |
| **B. Metrics interface** | `type RuntimeMetrics interface { IncActiveJobs(...); ... }` on scheduler | Embedder-friendly; no OTel required |
| **C. Logging only** | Structured completion event with duration per block | Minimal; no aggregation |

### Recommendation

**Option B + C in batch 5**: define small optional interfaces on `EvalContext` / `Scheduler` for counters and histograms (no-op default). Add **Option A** as separate batch if OTel dependency is acceptable in go.mod.

Span attributes: block id, type, eval stage, job id, retry attempt.

### Affected files

- `pkg/conflow/eval_context.go`
- `pkg/conflow/job/scheduler.go`, `tracker.go`
- `pkg/conflow/block/container.go`, `container_stage.go`
- `docs/product/runtime/embedding.md`

---

## I12: Retry policy

### Current behaviour

Retry logic is spread across:

- `container.go` — `runMainStage` interprets `Retryable` errors and `Result.Retry()`
- `container_stage.go` — panic recovery, `retryError`, schedules retry via tracker
- `tracker.go` — exponential backoff (`defaultBackoff` in container: base 1.57, 1s–15m)
- Runtime config / `@retry` directive — `RetryConfig.Limit` (-1 = infinite for main)

Default main stage uses `RetryConfig{Limit: -1}` when runtime config nil.

### Risk

- Implicit infinite retries on transient failures.
- Backoff parameters not documented for block authors.
- Retry reason not consistently attached to final error.

### Recommendation

1. Centralise policy in `pkg/conflow/retry/policy.go` (or `job/retry.go`): document `Limit` semantics, export default backoff constructor.
2. Add **jitter** to exponential backoff (reduce thundering herd).
3. Include retry count and last reason in terminal error message.
4. Consider safer default: finite limit unless `@retry` or runtime config explicitly sets infinite.

### Affected files

- `pkg/conflow/block/container.go`, `container_stage.go`
- `pkg/conflow/job/tracker.go`
- New: `pkg/conflow/job/retry_policy.go`
- `docs/product/language/runtime-directives.md` — `@retry`
