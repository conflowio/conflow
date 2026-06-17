---
title: Orchestration improvements
summary: State machine, async delivery, deadlock detection, and static/dynamic path unification.
parent: index.md
keywords: [orchestration, fsm, channels, deadlock, static]
---

# Orchestration

## I6: State machine

### Current behaviour

Block lifecycle uses numeric `int64` constants and implicit advancement via `containerStateNext` (state + 1):

```go
const (
    containerStatePending = 1  // iota after blank
    containerStateStart
    // ... through containerStateFinished
    containerStateNext      // sentinel: increment current state
)
```

Intermediate states like `PreMain` exist only as numeric successors — not named in switch logging or external APIs.

### Risk

- Hard to debug (log `state=5` vs name).
- Adding a stage requires careful iota ordering; easy to break increment chain.
- No single table describing valid transitions.

### Options

| Option | Description | Trade-offs |
|--------|-------------|------------|
| **A. Named enum + String()** | Replace ints with typed `ContainerState` and explicit transitions | Better logs; moderate refactor |
| **B. Table-driven FSM** | `map[State]map[Event]Transition` with actions | Most extensible; larger change |
| **C. Document only** | Add state diagram to docs; keep ints | No runtime improvement |

### Recommendation

**Option A for batch 4** (internal refactor): typed state, `String()`, explicit transition function. Defer table-driven FSM unless pause/checkpoint features are planned.

### Affected files

- `pkg/conflow/block/container.go`
- Optional: `pkg/conflow/block/state.go` (new)
- `pkg/conflow/block/container_test.go`

---

## I7: Async channel delivery

### Current behaviour

`SetChild` and `SetError` spawn goroutines to send on channels:

```go
func (c *Container) SetChild(container conflow.Container) {
    go func() { c.resultChan <- container }()
}
```

Buffers: `resultChan` = 8, `errChan` = 1. If `mainLoop` stalls, sender goroutines accumulate blocked on send.

Child completion order on `resultChan` is non-deterministic relative to scheduling order.

### Risk

- Goroutine leak under backpressure.
- Potential ordering surprises for dependents (mitigated by pub/sub, but stage advancement uses active job count).
- Harder to reason about under `-race` without dedicated tests.

### Options

| Option | Description | Trade-offs |
|--------|-------------|------------|
| **A. Synchronous send when on loop goroutine** | Detect caller; direct send if safe | Fragile without clear ownership |
| **B. Larger buffer + drop detection** | Increase buffers; metric on blocked sends | Masks problem |
| **C. Central dispatcher** | Single goroutine owns container; jobs callback synchronously | Cleaner model; refactor `Run()` structure |

### Recommendation

**Option C** as a follow-up after concurrency fixes: refactor so `mainLoop` is the only goroutine mutating container state; job completion posts events via non-blocking internal queue processed in loop. Until then, document channel buffer sizes and add `-race` tests (see [Testing](./testing.md)).

### Affected files

- `pkg/conflow/block/container.go`
- `pkg/conflow/parameter/container.go` (defer `SetChild`)

---

## I8: Deadlock detection

### Current behaviour

`evaluateChildren` detects one failure mode:

```go
if pending == total {
    return parsley.NewErrorf(..., "%q is deadlocked as no children could be evaluated", ...)
}
```

This catches all children waiting on deps at a stage start. It does **not** detect:

- Runtime circular wait (pub/sub never fires).
- Generator loops with no subscribers (`PublishBlock` returns `published=false`).
- Jobs stuck in retry/backoff indefinitely (`RetryConfig.Limit: -1` default for main stage).

### Options

| Option | Description | Trade-offs |
|--------|-------------|------------|
| **A. Workflow watchdog** | Optional global timeout on `Evaluate` (already via ctx) + log active nodes on deadline | Uses existing context; needs active-node registry |
| **B. Stuck detection** | If no progress (no publish/job complete) for N seconds, dump graph | Heuristic; false positives on long sleeps |
| **C. Static analysis only** | Rely on resolver cycles | Insufficient for runtime deadlocks |

### Recommendation

**Option A** (enhance embedder guidance) + minimal **Option B** debug mode: when `Logger` debug enabled and active job count unchanged for configurable interval, log block IDs and pending deps. No production abort from heuristic in batch 5.

### Affected files

- `pkg/conflow/block/container.go`
- `pkg/conflow/job/tracker.go` (optional last-activity timestamp)
- `docs/product/runtime/evaluation-pipeline.md`

---

## I9: Static vs dynamic paths

### Current behaviour

Two evaluation paths:

| Path | Type | Concurrency | Context |
|------|------|-------------|---------|
| Dynamic | `block.Container` | Async, pub/sub, job tracker | Background-derived (see I3) |
| Static | `block.StaticContainer` | Sync, sequential children | `context.Background()` |

Features added to one path (timeouts, skip, bind, isolation) must be duplicated manually.

### Risk

- Behaviour drift between static and dynamic evaluation.
- Static path used in codegen/tests may not reflect production runtime.

### Options

| Option | Description | Trade-offs |
|--------|-------------|------------|
| **A. Shared evaluator core** | Extract common bind/eval steps; static calls with `async=false` | Medium refactor; single semantics |
| **B. Deprecate static** | Route all through `Container` | Simplest long-term; may affect perf in tests |
| **C. Test-only static** | Document static as test helper only | Low effort; production still dual |

### Recommendation

**Option A** in a late batch after I3/I4: introduce internal `evaluateNode(ctx, mode EvalMode)` shared by both containers. **Option C** documentation immediately if static is primarily for tests.

Audit call sites of `NewStaticContainer` before batch 6.

### Affected files

- `pkg/conflow/block/static_container.go`
- `pkg/conflow/block/container.go`
- Grep: `NewStaticContainer` call sites
