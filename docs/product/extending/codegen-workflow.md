---
title: Code generation workflow
summary: conflow generate CLI, *.cf.go output, module-mode paths, and go generate integration.
parent: extending/index.md
keywords: [code generation, conflow generate, cf.go]
---

# Code generation workflow

## CLI

```bash
conflow generate [path]
```

Subcommand: `cmd/conflow/generate/cmd.go`

Behavior:

1. Walk directory for `.go` files (skips `vendor/`)
2. Find files with `// @block` or `// @function` in comments
3. Parse Go AST and emit generated code

**Target path:** Defaults to the current working directory, or pass a relative/absolute directory. Works in Go module mode from any checkout location.

Flags:

- `--local` — import grouping for generated files

## Generated artifacts

| Output | Content |
|--------|---------|
| `typename.cf.go` | `TypeNameInterpreter`, schema helpers, object methods |
| Same package or subpath | Interpreter may be generated in `@key` path if configured |

Files include a generated header from `pkg/conflow/generator/template/header.go`.

**Do not edit `*.cf.go` manually** — regenerate after changing structs, fields, or annotations.

## Typical developer loop

```bash
# Edit hello.go
go generate   # if //go:generate conflow generate is set
# or
conflow generate .
go build
```

## What generation produces

For blocks:

- `BlockInterpreter` implementation
- JSON Schema–compatible `Schema()` metadata
- `SetParam` / `SetBlock` / `Param` wiring
- Type validation on assignment

For functions:

- Callable interpreter with typed arguments and return schema

## Generator internals

| Step | Package |
|------|---------|
| AST walk | `pkg/conflow/generator/generate.go` |
| Struct parse | `pkg/conflow/block/generator/` |
| Function parse | `pkg/conflow/function/generator/` |
| Write file | `pkg/conflow/generator/write_file.go` |

## OpenAPI generation (separate command)

```bash
conflow openapi generate go|json|yaml
```

See [OpenAPI integration](../integrations/openapi.md).

## Troubleshooting

| Issue | Check |
|-------|-------|
| No `*.cf.go` emitted | `@block` comment on struct, not pointer type |
| generate fails on path | Pass `.` or a directory containing annotated `.go` files |
| Stale schema | Re-run generate after field changes |
| Import cycles | Interpreter path / package layout |

## See also

- [Defining blocks in Go](./go-blocks.md)
- [CLI reference](../runtime/cli.md)
