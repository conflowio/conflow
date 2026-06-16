---
title: Example catalog
summary: Per-example descriptions, features demonstrated, and key files.
parent: examples/index.md
keywords: [examples, catalog, helloworld, iterator]
---

# Example catalog

| Example | Path | Demonstrates |
|---------|------|--------------|
| helloworld | `examples/helloworld/` | Custom `@block "task"`, `Init`/`Run`, `go generate`, minimal host |
| iterator | `examples/iterator/` | `iterator` generator, nested loops, array indexing |
| exec | `examples/exec/` | `exec` block, stdout/stderr streams, `println` |
| streams | `examples/streams/` | Stream processing patterns |
| ticker | `examples/ticker/` | `ticker` generator, timed events |
| sleep | `examples/timeout/` | `sleep` block, durations (directory name `timeout`) |
| retry | `examples/retry/` | `@retry`, failing `fail` block |
| triggers | `examples/triggers/` | `@triggers` conditional scheduling |
| inputs | `examples/inputs/` | `@input` runtime parameters |
| modules | `examples/modules/` | `import` block, submodule `sum/` |
| multifile | `examples/multifile/` | Multiple `.cf` files in one program |
| jsonschema | `examples/jsonschema/` | External `person.json`, validated objects |
| openapi | `examples/openapi/` | Full OpenAPI 3 petstore in Conflow |
| licensify | `examples/licensify/` | File walker generator, real-world-ish workflow |
| benchmark | `examples/benchmark/` | Performance benchmarking setup |
| common | `examples/common/` | Shared parse context and main runner (library) |
| retry/fail | `examples/retry/fail.go` | `fail` block implementation for retry demo |

## helloworld

**Files:** `hello.go`, `main.go`, `main.cf`

- Defines `Hello` task block with random greeting
- Registers `hello`, `print`, `println` on main
- Conflow: `hello { to = "World" }`

**Docs:** [Go blocks](../extending/go-blocks.md), [Embedding](../runtime/embedding.md)

## iterator

**Files:** `main.cf`, uses `examples/common/iterator`

- Nested `iterator` over `colors` and `shapes` arrays
- Shows `:=` binding from `i1.value`

**Docs:** [Generators](../concepts/generators.md)

## exec

**Files:** `main.cf`

- Shell command with staged stdout
- Separate `println` blocks for stdout and stderr streams

**Docs:** [Built-in blocks: exec](../reference/built-in-blocks.md)

## modules

**Files:** `main.cf`, `modules/sum/sum.cf`

```conflow
sum import "./sum"
one_plus_two sum { a = 1, b = 2 }
```

**Docs:** [Built-in blocks: import](../reference/built-in-blocks.md)

## jsonschema

**Files:** `main.cf`, `person.json`

- `@doc` directive
- Person/spouse/pet types from JSON Schema file

**Docs:** [JSON Schema integration](../integrations/json-schema.md)

## openapi

**Files:** `petstore.cf`, `Makefile`

- Complete OpenAPI 3 definition
- `conflow openapi generate` targets

**Docs:** [OpenAPI integration](../integrations/openapi.md)

## licensify

**Files:** `licensify.go`, `file_walker.go`, `main.cf`

- Custom generators walking files
- Larger workflow composition

## triggers

**Files:** `main.cf`

- Two iterators + sleeps with different durations
- `println` only triggered by `sleep2` completion

**Docs:** [Runtime directives](../language/runtime-directives.md)

## Agent navigation tips

When implementing a feature similar to an example:

1. Read the example's `main.cf` for Conflow usage
2. Read `main.go` for interpreter registration
3. Read custom `*.go` (non-`cf.go`) for block logic
4. Cross-check [Built-in blocks](../reference/built-in-blocks.md) if using standard blocks

## See also

- [Examples index](./index.md)
- [Product overview](../overview.md)
