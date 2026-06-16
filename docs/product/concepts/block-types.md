---
title: Block types
summary: main, task, generator, configuration, and directive block types and their roles.
parent: concepts/index.md
keywords: [block types, main, task, generator, directive]
---

# Block types

Block types are declared in Go with `// @block` and an optional `type` metadata field. Constants are defined in `pkg/conflow/block.go`.

## Summary table

| Type | Constant | Evaluated? | Typical Go pattern |
|------|----------|------------|-------------------|
| `main` | `BlockTypeMain` | Yes | Root struct; registers child interpreters |
| `task` | `BlockTypeTask` | Yes | `Run()` business logic |
| `generator` | `BlockTypeGenerator` | Yes | `Run()` + `PublishBlock` for children |
| `configuration` | `BlockTypeConfiguration` | Usually parse-only | OpenAPI/schema trees |
| `directive` | `BlockTypeDirective` | Metadata | `ApplyToRuntimeConfig` / `ApplyToParameterConfig` |

## `main`

- **One per program** (per module).
- The `.cf` file content is the body of `main`.
- Often implements `ParseContextOverride()` to register custom block interpreters.
- Entry point for `conflow.Evaluate(..., "main", inputParams)`.

Example: `examples/helloworld/main.go`.

## `task`

- Standard runnable block.
- Implements `BlockRunner` (`Run`).
- Optional `BlockInitialiser` (`Init`) and `BlockCloser` (`Close`).
- Examples: `exec`, `sleep`, `println`, custom `Hello` block.

```go
// @block "task"
type Hello struct {
    // @id
    id conflow.ID
    // @required
    to string
}
```

## `generator`

- Emits **multiple instances** of **generated** child block types during `Run`.
- Generated block types are declared in the generator body in Conflow (e.g. `i1 it`).
- Uses `BlockPublisher.PublishBlock` to schedule dependent blocks.
- Dependency resolver splits generator start/finish to avoid cycles (`pkg/conflow/dependency/resolver.go`).

Examples:

- `iterator` / `it` — `pkg/blocks/iterator.go`, `examples/iterator`
- `ticker` / `tick` — `pkg/blocks/ticker.go`, `examples/ticker`
- `exec` stdout/stderr streams — `pkg/blocks/exec.go`

## `configuration`

- Structured definition blocks without runtime task semantics.
- Used heavily in **OpenAPI** (`pkg/openapi/`) and JSON Schema loading.
- Often evaluated at parse or resolve stages only.

## `directive`

Two families:

1. **Runtime directives** (`pkg/directives/`) — `@retry`, `@timeout`, `@skip`, `@triggers`, `@input`, `@doc`, …
2. **Schema directives** (`pkg/schema/directives/`) — `@required`, `@format`, `@generated`, `@dependency`, …

Directives are blocks with `eval_stage` often set to `ignore`, `resolve`, or `init` so they configure metadata rather than run workflow logic.

## Declaring type in Go

Single-line:

```go
// @block "task"
```

Multi-line metadata:

```go
//	@block {
//	  type = "directive"
//	  eval_stage = "init"
//	}
```

## See also

- [Generators](./generators.md)
- [Lifecycle and evaluation stages](./lifecycle-and-stages.md)
- [Runtime directives](../language/runtime-directives.md)
