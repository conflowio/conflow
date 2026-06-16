---
title: Lifecycle and evaluation stages
summary: Block init-main-close lifecycle, evaluation stages, lazy evaluation, and conditional skip.
parent: concepts/index.md
keywords: [lifecycle, init, run, close, eval stages]
---

# Lifecycle and evaluation stages

## Block instance lifecycle

Each block instance follows a multi-stage lifecycle:

| Stage | Go interface | Purpose |
|-------|--------------|---------|
| Init | `BlockInitialiser.Init(ctx) (skipped bool, err)` | Pre-main setup; `skipped=true` skips the block |
| Main | `BlockRunner.Run(ctx) (Result, error)` | Primary business logic |
| Close | `BlockCloser.Close(ctx) error` | Cleanup after main |

Not every block implements all three. Missing interfaces are no-ops.

### Conditional execution

`Init` returning `skipped=true` means the block (and typically its dependents' evaluation for its outputs) is skipped — useful for feature flags or optional steps.

## Evaluation stages (system)

Defined in `pkg/conflow/eval_context.go`:

| Stage | Constant | Typical use |
|-------|----------|-------------|
| `parse` | `EvalStageParse` | Import modules, parse-time registration |
| `resolve` | `EvalStageResolve` | Directives like `@triggers` |
| `init` | `EvalStageInit` | Runtime directives `@retry`, block `Init` |
| `main` | `EvalStageMain` | Parameter evaluation, `Run` |
| `close` | `EvalStageClose` | `Close` |
| `ignore` | `EvalStageIgnore` | Directives that only affect schema/metadata |

Field and directive annotation `@eval_stage` assigns when a parameter or directive is processed.

## Lazy evaluation

Parameters and child blocks inside a block are **lazy**:

- They are evaluated only when the matching **stage** is active.
- They are evaluated only when **dependencies** are ready.

This avoids running expensive inputs before `Init` decides to skip, or before upstream blocks produce outputs.

## Parameter `@eval_stage`

On struct fields, `@eval_stage "init"` defers evaluation until the init stage (e.g. fields read only in `Init()`).

## Block interfaces summary

```go
// Optional
func (b *Sample) Init(ctx context.Context) (bool, error) { return false, nil }

// Required for tasks
func (b *Sample) Run(ctx context.Context) (conflow.Result, error) { return nil, nil }

// Optional
func (b *Sample) Close(ctx context.Context) error { return nil }
```

## Hello example stages

`examples/helloworld/hello.go`:

- `Init` — seeds random number generator
- `Run` — picks greeting into `@read_only` field `greeting`

## See also

- [Block types](./block-types.md)
- [Dependencies and evaluation order](./dependencies-and-order.md)
- [Evaluation pipeline](../runtime/evaluation-pipeline.md)
