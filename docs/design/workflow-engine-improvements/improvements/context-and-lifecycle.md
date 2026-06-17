---
title: Context and lifecycle improvements
summary: Context propagation, block notifications on terminal states, and graceful shutdown.
parent: index.md
keywords: [context, cancellation, lifecycle, shutdown, timeout]
---

# Context & lifecycle

## I3: Context propagation

### Current behaviour

`NodeContainer.createEvalContext` creates child contexts from **`context.Background()`**, not from the parent workflow context:

```go
// pkg/conflow/node_container.go
if timeout > 0 {
    ctx, cancel = context.WithTimeout(context.Background(), timeout)
} else {
    ctx, cancel = context.WithCancel(context.Background())
}
```

Effects:

- Root cancellation/deadline does not propagate to child node eval contexts automatically.
- `@timeout` on a block replaces parent context entirely rather than constraining it.
- Context values set on the root context are invisible to child eval contexts.

`StaticContainer.createContainer` has the same pattern (`static_container.go`).

### Risk

- `Evaluate(ctx, ...)` cancellation may not stop in-flight block/parameter jobs promptly.
- Per-block timeouts can outlive a cancelled workflow.
- Incorrect semantics for embedders using context for request-scoped data.

### Options

| Option | Description | Trade-offs |
|--------|-------------|------------|
| **A. Derive from parent** | `context.WithTimeout(parentCtx, timeout)` / `WithCancel(parentCtx)` | Correct cancellation tree; `@timeout` becomes min(parent, block) |
| **B. Merge deadlines** | Derive child ctx; if parent deadline is sooner, use it | Slightly more code; clearest semantics |
| **C. Document-only** | Keep behaviour; document that only root ctx matters | No code change; leaves bug class open |

### Recommendation

**Option B**: always derive from `n.ctx.ctx` (parent eval context's underlying context). When `@timeout` is set, apply `WithTimeout(parent, timeout)`. When not set, use `WithCancel(parent)`.

Add tests: root cancel aborts child parameter eval; block timeout fires before parent deadline when shorter.

### Affected files

- `pkg/conflow/node_container.go` — `createEvalContext`
- `pkg/conflow/block/static_container.go` — `createContainer`
- `pkg/conflow/node_container_test.go` — cancellation cases
- `docs/product/language/runtime-directives.md` — `@timeout` semantics

---

## I4: Lifecycle notifications

### Current behaviour

Terminal states in `Container.setState` log but do not notify blocks:

```go
case containerStateSkipped:
    // TODO: notify block about skipped task
case containerStateErrored:
    // TODO: notify block about error
case containerStateAborted:
    // TODO: notify block about abort
```

`Close` runs only on the normal path through `containerStateClose`. Skip, error, and abort may leave resources held unless blocks defensively clean up inside `Run`.

### Risk

- Leaked connections, temp files, or partial state on abort/error/skip.
- Block authors cannot rely on engine for cleanup guarantees.

### Options

| Option | Description | Trade-offs |
|--------|-------------|------------|
| **A. Always call Close** | Invoke `BlockCloser.Close` on all terminal paths (with documented semantics) | Simple; Close must be idempotent; may surprise authors who skip Close on error today |
| **B. New optional interfaces** | `BlockSkipHandler`, `BlockAbortHandler`, `BlockErrorHandler` | Explicit; no change for existing blocks |
| **C. Context cancellation only** | Rely on `ctx.Done()` in Run; no new hooks | No API change; easy to miss in block code |

### Recommendation

**Option A + B hybrid**:

1. On **error/abort**, call `Close` if `BlockCloser` is implemented (document as best-effort cleanup).
2. On **skip** (from `Init`), do **not** call `Close` by default; add optional `BlockSkipNotifier` if authors need skip-specific logic.
3. Replace TODOs with implemented paths and tests.

Resolve open question in [initiative index](../index.md) before implementation.

### Affected files

- `pkg/conflow/block/container.go` — `setState` terminal cases
- `pkg/conflow/block.go` — optional new interfaces
- `pkg/conflow/block/container_test.go`
- `docs/product/concepts/lifecycle-and-stages.md`

---

## I5: Graceful shutdown

### Current behaviour

After `mainLoop` exits, `Container.Run`:

1. Closes all child `NodeContainer` subscriptions.
2. Calls `jobTracker.Stop()` to cancel jobs implementing `Cancellable`.
3. If running jobs remain, enters `shutdownLoop` for `ContainerGracefulTimeoutSec` (10s, package var).
4. Logs child errors during shutdown but **does not propagate** them to the block result.

`shutdownLoop` drains `stateChan` without processing state transitions.

### Risk

- Fixed 10s timeout may be too short for long-close blocks or too long for fast fail.
- Silent loss of shutdown errors.
- No coordination with root context cancellation during shutdown.

### Options

| Option | Description | Trade-offs |
|--------|-------------|------------|
| **A. Configurable timeout** | `RuntimeConfig.ShutdownTimeout` or scheduler-level setting | Flexible; more API surface |
| **B. Context-driven shutdown** | Shutdown until `parentCtx.Done()` or explicit timeout | Aligns with I3 |
| **C. Aggregate errors** | Collect shutdown failures into `Container.err` (multi-error) | Better operability; error type decisions |

### Recommendation

**Option A + B**: derive shutdown deadline from parent context if set; otherwise use configurable timeout (default 10s). **Option C** for errors: append to a slice, return primary error with `%w` chain or parsley multi-error if available.

### Affected files

- `pkg/conflow/block/container.go` — `shutdownLoop`, `Run` defer path
- `pkg/conflow/block_directive.go` — `RuntimeConfig`
- `docs/product/runtime/evaluation-pipeline.md`
