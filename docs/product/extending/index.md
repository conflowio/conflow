---
title: Extending Conflow
summary: Index for defining custom blocks and functions in Go and running code generation.
parent: index.md
keywords: [extending, custom blocks, codegen]
---

# Extending Conflow

Extension authors implement **blocks** (structs) and **functions** (Go functions) in plain Go, annotate them with comments, run **`conflow generate`**, and register generated interpreters on the parse context.

## Topics

| Document | Summary |
|----------|---------|
| [Go blocks](./go-blocks.md) | Struct blocks, interfaces, dependencies |
| [Go functions](./go-functions.md) | `@function` and interpreter registration |
| [Annotations reference](./annotations.md) | Comment directive syntax |
| [Code generation workflow](./codegen-workflow.md) | CLI, `*.cf.go`, `go generate` |

## Minimal path

1. Write struct with `// @block "task"` and field annotations.
2. Implement `ID()`, `Run()`, optional `Init`/`Close`.
3. Run `conflow generate` (or `go generate`).
4. Register `YourBlockInterpreter{}` on main's `ParseContextOverride`.
5. Parse `.cf` files and `Evaluate`.

Hello-world reference: `examples/helloworld/`.

## Generator packages

| Package | Role |
|---------|------|
| `pkg/conflow/block/generator/` | Block interpreter generation |
| `pkg/conflow/function/generator/` | Function interpreter generation |
| `pkg/conflow/generator/` | CLI driver, file walk, templates |

## See also

- [Schema annotations](../reference/schema-annotations.md)
- [Embedding](../runtime/embedding.md)
- [Built-in blocks](../reference/built-in-blocks.md) as templates
