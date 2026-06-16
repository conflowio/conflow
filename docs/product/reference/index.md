---
title: Reference
summary: Index for built-in blocks, functions, and schema annotation reference.
parent: index.md
keywords: [reference, builtins]
---

# Reference

Quick lookup for standard library surface area shipped in this repository.

## Documents

| Document | Content |
|----------|---------|
| [Built-in blocks](./built-in-blocks.md) | `pkg/blocks/` task and utility blocks |
| [Built-in functions](./built-in-functions.md) | `pkg/functions/` callable functions |
| [Schema annotations](./schema-annotations.md) | Field and schema directives |

## Registries (source of truth)

| Registry | File |
|----------|------|
| Functions | `pkg/functions/registry.go` |
| Runtime directives | `pkg/directives/directives.go` |
| Schema directives | `pkg/schema/directives/directives.go` |

## Version

`conflow.Version` in `pkg/conflow/version.go`; CLI `conflow version`.

## See also

- [Extending](../extending/index.md) — custom blocks/functions
- [Language](../language/index.md) — syntax for using builtins
