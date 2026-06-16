---
title: Agent guide
summary: Navigation guide for LLM agents working in the Conflow repository.
parent: index.md
keywords: [agent, llm, navigation, tasks]
audience: agent
---

# Agent guide

Structured navigation for automated agents implementing or debugging Conflow.

## Product summary (one paragraph)

Conflow is a Go library and CLI for building **typed DSLs**: authors define blocks (structs) and functions in Go with comment annotations, run `conflow generate` to create `*.cf.go` interpreters, register interpreters on a `ParseContext`, parse `.cf` files, and `Evaluate` the `main` block. Execution is **dependency-driven** and **parallel** via a job scheduler. Generators publish dynamic child blocks at runtime.

## Task → start here

| Task | Read first | Key code paths |
|------|------------|----------------|
| Understand what Conflow is | [Overview](./overview.md) | `README.md` |
| Parse/evaluate a program | [Embedding](./runtime/embedding.md) | `pkg/conflow/eval.go`, `examples/common/main.go` |
| Add a custom block | [Go blocks](./extending/go-blocks.md) | `examples/helloworld/hello.go` |
| Add a custom function | [Go functions](./extending/go-functions.md) | `pkg/functions/strings/lower.go` |
| Run code generation | [Codegen workflow](./extending/codegen-workflow.md) | `cmd/conflow/generate/`, `pkg/conflow/generator/` |
| Use iterators/streams | [Generators](./concepts/generators.md) | `pkg/blocks/iterator.go`, `pkg/blocks/exec.go` |
| Runtime directives | [Runtime directives](./language/runtime-directives.md) | `pkg/directives/` |
| JSON Schema validation | [JSON Schema](./integrations/json-schema.md) | `pkg/schema/` |
| OpenAPI export | [OpenAPI](./integrations/openapi.md) | `pkg/openapi/`, `examples/openapi/` |
| Debug dependency order | [Dependencies](./concepts/dependencies-and-order.md) | `pkg/conflow/dependency/resolver.go` |
| Find builtin API | [Reference](./reference/index.md) | `pkg/functions/registry.go`, `pkg/blocks/` |
| Match an example | [Example catalog](./examples/catalog.md) | `examples/` |

## Documentation tree

```text
docs/
  index.md
  product/
    index.md
    overview.md
    agent-guide.md          ← this file
    concepts/               ← blocks, dependencies, lifecycle, generators
    language/               ← syntax, types, @directives in .cf
    extending/              ← Go blocks, functions, codegen
    runtime/                ← embed, evaluate, scheduler, CLI
    integrations/           ← JSON Schema, OpenAPI
    reference/              ← builtins, schema annotations
    examples/               ← example catalog
```

Every directory has `index.md`. Follow `parent` links in YAML front matter to ascend.

## Terminology (stable)

| Term | Meaning |
|------|---------|
| Block | Typed node in workflow graph |
| Interpreter | Generated `XInterpreter` in `*.cf.go` |
| ParseContext | Registry + file set for parse/eval |
| main | Root block; `.cf` file is its body |
| Generator | Block that `PublishBlock` for children |
| Directive | Metadata block (`@retry` or `@required`) |

## Code generation rules

1. Edit `*.go` files with annotations — **never** hand-edit `*.cf.go`
2. Run `conflow generate` from the module root or pass a target path (relative or absolute)
3. Register `NewXInterpreter{}` on main's `ParseContextOverride`

## Common pitfalls

| Symptom | Likely cause |
|---------|--------------|
| Block not found in Conflow | Interpreter not registered on main |
| Circular dependency | Generator + dependent reference same block fields |
| Input param rejected | Not marked `@input` or wrong type at `Evaluate` |
| generate fails on path | Pass `.` or an explicit directory; ensure it contains annotated `.go` files |
| import module fails | Missing `module` block in registry / no `main` in module |

## Test locations

| Area | Tests |
|------|-------|
| Parsers | `pkg/parsers/*_test.go` |
| Schema | `pkg/schema/*_test.go` |
| Functions | `pkg/functions/**/*_test.go` |
| Directives | `pkg/directives/*_test.go` |
| Integration | `examples/*/main.go` (manual), `pkg/conflow/*_test.go` |

## Front matter convention

Each doc file includes YAML:

```yaml
title: ...
summary: ...        # use for quick context injection
parent: ...         # relative path to parent index
keywords: [...]
```

## License note

Source files use Mozilla Public License 2.0 headers.
