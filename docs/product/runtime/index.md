---
title: Runtime
summary: Index for embedding Conflow, evaluation pipeline, job scheduler, and CLI.
parent: index.md
keywords: [runtime, embedding, evaluation]
---

# Runtime

The Conflow **runtime** parses Conflow source into an AST, resolves dependencies, schedules block jobs, and executes Go block implementations. It lives primarily in `pkg/conflow/`.

## Topics

| Document | Summary |
|----------|---------|
| [Embedding](./embedding.md) | ParseContext, Evaluate, minimal host app |
| [Evaluation pipeline](./evaluation-pipeline.md) | Parse → resolve → stages → Run |
| [Job scheduler](./job-scheduler.md) | Concurrent block execution |
| [CLI](./cli.md) | `conflow` binary commands |

## Core types

| Type | Role |
|------|------|
| `ParseContext` | Block registry, file set, directive registry |
| `EvalContext` | Per-evaluation state, dependencies, pub/sub |
| `BlockInterpreter` | Generated bridge between AST and Go struct |
| `JobScheduler` | Async execution queue |

## Dependencies

- **Parsley** (`github.com/conflowio/parsley`) — parser/evaluator framework
- **Zerolog** — logging in examples (`pkg/loggers/zerolog/`)

## See also

- [Core concepts](../concepts/index.md)
- [Extending](../extending/index.md)
