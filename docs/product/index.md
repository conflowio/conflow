---
title: Conflow product
summary: Product-level overview of Conflow as a configuration and workflow DSL platform.
parent: ../index.md
keywords: [conflow, product, overview, dsl, workflow]
---

# Conflow product

Conflow is a **technology preview** (see [README](../README.md)) for building **strongly typed configuration and workflow languages** in Go.

## What problem it solves

| Problem | Conflow approach |
|---------|------------------|
| YAML/JSON config lacks types and validation | JSON Schema–aligned types with static checking before evaluation |
| Workflow tools are often GUI or opaque DSLs | **Workflow-as-code**: Go structs define blocks; Conflow text defines graphs |
| Dependency ordering in pipelines is manual | **Parallel evaluation**: code runs when dependencies are satisfied |
| Runtime reflection in DSL hosts is slow/fragile | **Code generation** produces interpreters; no reflection at runtime |

## Core value proposition

1. **Define your language in Go** — structs with comment annotations become block types; functions become callable builtins.
2. **Write programs in Conflow** — block graphs with parameters, expressions, and imports.
3. **Generate interpreters** — `conflow generate` emits `*.cf.go` next to your Go sources.
4. **Parse, validate, and evaluate** — embed `pkg/conflow` in an application; schedule block work on a job pool.

## Product surface

| Surface | Description | Doc |
|---------|-------------|-----|
| Conflow language | Text syntax for blocks, parameters, expressions | [Language](./language/index.md) |
| Block model | Types, lifecycle, generators, dependencies | [Concepts](./concepts/index.md) |
| Go extension API | `@block`, `@function`, field annotations | [Extending](./extending/index.md) |
| Runtime | Parse context, evaluation, job scheduler | [Runtime](./runtime/index.md) |
| Integrations | JSON Schema files, OpenAPI 3 definitions | [Integrations](./integrations/index.md) |
| Built-ins | Standard blocks and functions | [Reference](./reference/index.md) |
| Examples | Runnable projects in `examples/` | [Examples](./examples/index.md) |

## Typical workflows

### Configuration DSL

1. Define block structs representing config sections.
2. Run `conflow generate`.
3. Register interpreters on a `main` block's `ParseContextOverride`.
4. Parse `.cf` files and `Evaluate` with optional `@input` parameters.

### Workflow-as-code

1. Define **task** blocks with `Run()` business logic.
2. Use **generator** blocks for iterators, tickers, or dynamic children.
3. Compose workflows in Conflow; dependencies determine order and parallelism.
4. Apply runtime directives (`@retry`, `@timeout`, `@triggers`).

### OpenAPI / schema-driven config

1. Define OpenAPI or JSON Schema in Conflow or import JSON Schema files.
2. Use `conflow openapi generate` for Go/JSON/YAML output.
3. Validate structured config against generated schemas.

## Status and license

- **Status**: Technology preview; APIs and syntax may change.
- **License**: Mozilla Public License 2.0 (see file headers and `vendor/`).

## Drill-down

- [Agent guide](./agent-guide.md) — task → doc/code mapping for LLM agents
- [Overview](./overview.md) — extended product narrative
- [Concepts](./concepts/index.md) — blocks, dependencies, lifecycle, generators
- [Language](./language/index.md) — syntax and directives
- [Extending](./extending/index.md) — building custom blocks and functions
- [Runtime](./runtime/index.md) — embedding and evaluation
- [Integrations](./integrations/index.md) — JSON Schema and OpenAPI
- [Reference](./reference/index.md) — built-in blocks, functions, annotations
- [Examples](./examples/index.md) — example catalog
