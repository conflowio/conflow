---
title: Examples
summary: Index for runnable examples in the examples/ directory.
parent: index.md
keywords: [examples, demos]
---

# Examples

The `examples/` directory contains runnable programs demonstrating Conflow features. Each example typically has:

- `main.go` — Go host and block registration
- `main.cf` — Conflow program
- `Makefile` — build/run shortcuts
- `*.cf.go` — generated interpreters

## Catalog

See [Example catalog](./catalog.md) for per-example descriptions and learning paths.

## Running an example

```bash
cd examples/helloworld
make   # or: go run .
```

Requirements: Go 1.26+, generated `*.cf.go` present (run `conflow generate` if missing).

## Shared utilities

`examples/common/` — `NewParseContext()`, `Main()` with scheduler and logging.

## Suggested learning order

1. `helloworld` — custom block + main registration
2. `iterator` — generators and dependencies
3. `exec` — subprocess + streams
4. `modules` — import and reusable modules
5. `inputs` — runtime `@input`
6. `retry` / `timeout` / `triggers` — runtime directives
7. `jsonschema` — external schema
8. `openapi` — API definition export

## See also

- [Embedding](../runtime/embedding.md)
- [README](../../README.md) hello-world section
