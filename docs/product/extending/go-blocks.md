---
title: Defining blocks in Go
summary: How to implement custom task, generator, and main blocks with Conflow interfaces.
parent: extending/index.md
keywords: [go blocks, Run, Init, BlockPublisher]
---

# Defining blocks in Go

## Basic task block

```go
// @block "task"
type Hello struct {
    // @id
    id conflow.ID
    // @required
  to string
    // @read_only
    greeting string
}

func (h *Hello) ID() conflow.ID {
    return h.id
}

func (h *Hello) Run(ctx context.Context) (conflow.Result, error) {
    h.greeting = "Hello " + h.to + "!"
    return nil, nil
}
```

After `conflow generate`, use `HelloInterpreter{}` in the block registry.

## Required methods

| Method | Purpose |
|--------|---------|
| `ID() conflow.ID` | Instance identifier from Conflow |

For **tasks**, implement `Run`. For **generators**, `Run` publishes children.

## Optional interfaces

| Interface | Method |
|-----------|--------|
| `BlockInitialiser` | `Init(ctx) (skipped bool, err error)` |
| `BlockCloser` | `Close(ctx) error` |
| `BlockProvider` | `BlockInterpreters(parseCtx) (map[ID]BlockInterpreter, error)` |
| `ParseContextOverride` | Customize registries for child block types |

`ParseContextOverride` example — `exec` registers `stdout`/`stderr` stream types (`pkg/blocks/exec.go`).

## Dependency injection fields

Mark runtime-injected fields with `@dependency`:

```go
// @dependency
blockPublisher conflow.BlockPublisher
```

Common dependencies:

- `conflow.BlockPublisher` — generators
- `io.Writer` — output sinks
- Custom services passed via `EvalContext.UserContext`

## Generated child fields

```go
// @generated
it *It
```

Child type `It` is typically another `@block` struct.

## Main block pattern

```go
// @block "main"
type Main struct {
    // @id
    id conflow.ID
}

func (m *Main) ParseContextOverride() conflow.ParseContextOverride {
    return conflow.ParseContextOverride{
        BlockTransformerRegistry: block.InterpreterRegistry{
            "hello": HelloInterpreter{},
            "print": blocks.PrintInterpreter{},
        },
    }
}
```

Reference: `examples/helloworld/main.go`.

## Parse-stage blocks

`import` uses `eval_stage = "parse"` and `BlockInterpreters` to load modules at parse time (`pkg/blocks/import.go`).

## Block type metadata

```go
// @block "generator"
// @block "configuration"
//	@block {
//	  type = "directive"
//	  eval_stage = "init"
//	}
```

## See also

- [Annotations reference](./annotations.md)
- [Generators](../concepts/generators.md)
- [Code generation workflow](./codegen-workflow.md)
