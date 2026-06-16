---
title: Conflow language
summary: Index for Conflow surface syntax, value types, and runtime directives.
parent: index.md
keywords: [language, syntax, conflow]
---

# Conflow language

The **Conflow language** is the text format users write in `.cf` files. It describes the body of the `main` block (and imported modules). Syntax is EBNF-oriented and Go-like; full formal spec is TODO in README.

## Topics

| Document | Summary |
|----------|---------|
| [Syntax overview](./syntax-overview.md) | Blocks, assignments, comments, imports |
| [Value types and expressions](./value-types.md) | Literals, collections, operators, calls |
| [Runtime directives](./runtime-directives.md) | `@retry`, `@timeout`, `@input`, `@triggers`, … |

## File conventions

| Extension | Role |
|-----------|------|
| `.cf` | Conflow source |
| `.cf.go` | Generated interpreter (do not edit by hand) |

## Parser

Implemented in `pkg/parsers/`. Main entry parsers:

- `parsers.NewMain("main", MainInterpreter{})` — parse program root
- Block-specific parsers registered via interpreters

## Relationship to Go

| Conflow | Go analogue |
|---------|-------------|
| `block { field = value }` | Struct literal |
| `name := value` | Short variable declaration on block |
| `foo.bar` | Field access across blocks |
| `fn(a, b)` | Function call |

Types and validation come from **generated schema**, not from the Conflow file alone.

## See also

- [Blocks and parameters](../concepts/blocks-and-parameters.md)
- [Defining blocks in Go](../extending/go-blocks.md)
- [Example catalog](../examples/catalog.md)
