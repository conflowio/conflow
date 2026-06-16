---
title: Conflow documentation
summary: Root index for Conflow repository documentation.
keywords: [conflow, docs, index]
---

# Conflow documentation

Documentation for [Conflow](https://github.com/conflowio/conflow): a strongly typed configuration and workflow language implemented in Go.

## Sections

| Section | Purpose |
|---------|---------|
| [Product](./product/index.md) | What Conflow is, how it works, and how to use it |
| [Agent guide](./product/agent-guide.md) | Task-oriented navigation for LLM agents |
| [README](../README.md) | Quick introduction and hello-world walkthrough |

## Audience

- **Humans** building DSLs, workflows, or configuration systems on Conflow
- **LLM agents** navigating the repo to implement blocks, parse/evaluate workflows, or extend the language

## Repository map (high level)

| Path | Role |
|------|------|
| `pkg/conflow/` | Core runtime: parse context, evaluation, blocks, jobs |
| `pkg/parsers/` | Conflow language parser |
| `pkg/schema/` | JSON Schema–aligned type system and validation |
| `pkg/blocks/` | Built-in workflow blocks (`exec`, `sleep`, `import`, …) |
| `pkg/functions/` | Built-in functions (`len`, `str_*`, `json_*`, …) |
| `pkg/directives/` | Runtime directives (`@retry`, `@timeout`, `@input`, …) |
| `pkg/openapi/` | OpenAPI 3 definition blocks and code generation |
| `cmd/conflow/` | CLI (`generate`, `openapi generate`) |
| `examples/` | Runnable demonstrations |

## Conventions in this documentation

- **Block** — a typed unit of configuration or workflow logic
- **Parameter** — a named input, output, or child block on a block
- **Interpreter** — generated Go code (`*.cf.go`) that connects structs/functions to the parser
- **Main** — the root block; Conflow source files describe the body of `main`

Start with the [product overview](./product/index.md).
