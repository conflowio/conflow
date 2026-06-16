---
title: Defining functions in Go
summary: Custom Conflow-callable functions via @function and generated interpreters.
parent: extending/index.md
keywords: [functions, @function]
---

# Defining functions in Go

Functions extend the expression language with callable operations backed by Go implementations.

## Basic function

```go
// @function
func Lower(s string) string {
    return strings.ToLower(s)
}
```

After generation, register in the function registry (often alongside defaults from `functions.DefaultRegistry()`).

## Annotation

Single comment line before the function:

```go
// @function
```

Optional metadata can be added similar to blocks (see generator parser in `pkg/conflow/generator/parser/`).

## Registration

Function interpreters live in `function.InterpreterRegistry` (`pkg/conflow/function/`). Built-in registry:

`pkg/functions/registry.go` — merged when building `ParseContext`.

Custom hosts typically:

1. Start from `functions.DefaultRegistry()`
2. Add package-specific interpreters
3. Attach to parse context via main block or global parser setup

## Naming in Conflow

Generated function names follow Go function names, often with prefixes for packages (built-ins use `str_`, `json_`, etc.).

## Testing

Function packages include `*_test.go` beside implementations (e.g. `pkg/functions/strings/lower_test.go`).

## See also

- [Built-in functions](../reference/built-in-functions.md)
- [Code generation workflow](./codegen-workflow.md)
