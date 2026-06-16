---
title: Embedding Conflow
summary: Minimal pattern for parsing .cf files and evaluating main in a Go application.
parent: runtime/index.md
keywords: [embedding, ParseContext, Evaluate]
---

# Embedding Conflow

## Minimal host application

Pattern from `examples/helloworld/main.go` and `examples/common/main.go`:

```go
parseCtx := common.NewParseContext()

p := parsers.NewMain("main", MainInterpreter{})

if err := p.ParseFile(parseCtx, "main.cf"); err != nil {
    log.Fatal(err)
}

common.Main(ctx, parseCtx, inputParams)
```

`common.NewParseContext()` (`examples/common/main.go`):

```go
idRegistry := conflow.NewIDRegistry(8, 16)
return conflow.NewParseContext(
    parsley.NewFileSet(),
    idRegistry,
    directives.DefaultRegistry(),
)
```

## Evaluate API

`pkg/conflow/eval.go`:

```go
conflow.Evaluate(
    parseCtx *ParseContext,
    context context.Context,
    userContext interface{},
    logger Logger,
    scheduler JobScheduler,
    id ID,                    // usually "main"
    inputParams map[ID]interface{},
) (interface{}, error)
```

Behavior:

- Validates `@input` parameters against main block schema
- Builds `EvalContext` and runs Parsley evaluation on the main node
- Returns main block value or error with path-transformed messages

## Registering custom blocks

On the **main** block struct:

```go
func (m *Main) ParseContextOverride() conflow.ParseContextOverride {
    return conflow.ParseContextOverride{
        BlockTransformerRegistry: block.InterpreterRegistry{
            "myblock": MyBlockInterpreter{},
        },
    }
}
```

Built-in blocks from `pkg/blocks/` can be mixed in (print, exec, import, …).

## User context

`userContext` is available in `EvalContext.UserContext` for injecting services (DB handles, config) into blocks via `@dependency` fields populated by custom wiring (if implemented in interpreters).

## Logging

Implement `conflow.Logger` or use `pkg/loggers/zerolog`. Set `CONFLOW_LOG` env var in examples.

## Cancellation

Pass a `context.Context` cancelled on SIGINT (`cmd/conflow/main.go` pattern) or application shutdown.

## Profiling

`examples/common/main.go` starts `http.ListenAndServe("localhost:6060", nil)` for `net/http/pprof`.

## Parse variants

| API | Use |
|-----|-----|
| `ParseFile` | Single `.cf` file |
| `ParseDir` | Directory of Conflow sources (modules, multifile) |

Parsers: `pkg/parsers/`.

## See also

- [Evaluation pipeline](./evaluation-pipeline.md)
- [Job scheduler](./job-scheduler.md)
- [Runtime directives](../language/runtime-directives.md)
