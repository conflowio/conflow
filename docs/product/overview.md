---
title: Product overview
summary: Extended narrative of Conflow's purpose, architecture at the product level, and design choices.
parent: index.md
keywords: [overview, design, parallel evaluation, static typing]
---

# Product overview

## Positioning

Conflow sits between **static config formats** (YAML, TOML, JSON) and **general-purpose workflow engines**. Authors define a **domain-specific surface** in Go; users write **Conflow programs** that the runtime parses, type-checks, and executes as a **dependency-driven parallel graph**.

The README describes it as:

> A strongly typed configuration language focusing on simplicity and usability … able to **generate, parse and evaluate your own domain specific language (DSL)**.

## Design pillars

### 1. Strong static typing

Types derive from Go struct fields and JSON Schema rules (`@required`, `@minimum`, `@format`, etc.). Validation runs when values are assigned and when `@input` parameters are passed to `Evaluate`. Most type errors surface before block `Run()` executes.

Implementation: `pkg/schema/`, field directives in `pkg/schema/directives/`.

### 2. Parallel evaluation by dependency

Conflow is not a sequential script interpreter. Each block instance runs when:

- Its evaluation stage is active (`init`, `main`, `close`, …)
- All referenced parameters and child blocks are ready

The dependency resolver (`pkg/conflow/dependency/`) builds an order and detects cycles (Tarjan's SCC). A job scheduler (`pkg/conflow/job/`) runs block work concurrently.

### 3. Code generation, not reflection

Block and function hosts are **generated interpreters** (`StructNameInterpreter`, `FunctionNameInterpreter`) in `*.cf.go` files. The generator (`pkg/conflow/generator/`, CLI `conflow generate`) scans `// @block` and `// @function` comments.

### 4. Go-native extension model

Block authors write normal Go:

- `Init(ctx) (skipped bool, err)` — optional pre-main hook; `skipped=true` skips the block
- `Run(ctx) (Result, error)` — main logic
- `Close(ctx) error` — optional cleanup

Generators publish child blocks via `BlockPublisher.PublishBlock`.

### 5. Familiar syntax

Conflow syntax is intentionally close to Go where practical: blocks look like struct literals, `:=` for user parameters, expressions for arithmetic and calls. See [Language](./language/index.md).

## Product components (logical)

```text
┌─────────────────────────────────────────────────────────────┐
│  Author Go structs/functions + annotations                  │
└──────────────────────────┬──────────────────────────────────┘
                           │ conflow generate
                           ▼
┌─────────────────────────────────────────────────────────────┐
│  Generated interpreters (*.cf.go)                           │
└──────────────────────────┬──────────────────────────────────┘
                           │ register on ParseContext
                           ▼
┌─────────────────────────────────────────────────────────────┐
│  Conflow source (.cf)  ──►  Parser (pkg/parsers)            │
└──────────────────────────┬──────────────────────────────────┘
                           │ AST + schema
                           ▼
┌─────────────────────────────────────────────────────────────┐
│  Dependency resolution  ──►  Job scheduler  ──►  Block.Run  │
└─────────────────────────────────────────────────────────────┘
```

## Block type taxonomy

| Type constant | Role | Typical use |
|---------------|------|-------------|
| `main` | Root program | Entry block; body is the `.cf` file |
| `task` | Runnable unit | Business logic, I/O, side effects |
| `generator` | Emits child blocks | Iterators, tickers, exec stdout/stderr |
| `configuration` | Structured data | OpenAPI paths, schema objects |
| `directive` | Metadata / runtime config | `@doc`, `@retry`, schema `@required` |

Details: [Block types](./concepts/block-types.md).

## Comparison (informal)

| Aspect | YAML + scripts | Conflow |
|--------|----------------|---------|
| Types | External schema optional | Built into block definitions |
| Execution order | Manual orchestration | Dependency graph |
| Custom types | Ad hoc | Go structs + generate |
| Parallelism | External tool | Job scheduler in runtime |
| Learning curve | Low for simple config | Higher; pays off for workflows |

## When to use Conflow

**Good fit:**

- Typed configuration with validation and custom block types
- Workflow graphs with clear dependencies and optional parallelism
- Embedding a small language in a Go service
- OpenAPI-first API definition with Conflow as the authoring format

**Less ideal:**

- Single flat config file with no custom types
- Workflows requiring visual designer-first UX out of the box
- Non-Go ecosystems without a Go host process

## Next steps

- [Core concepts](./concepts/index.md)
- [Embedding guide](./runtime/embedding.md)
- [Example catalog](./examples/catalog.md)
