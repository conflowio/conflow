---
title: Runtime directives
summary: Conflow source directives that affect runtime behavior, inputs, and documentation.
parent: language/index.md
keywords: [directives, retry, timeout, triggers, input]
---

# Runtime directives

Runtime directives are **blocks of type `directive`** declared in Conflow source with `@name` syntax. They configure evaluation rather than implementing domain logic.

Implementation: `pkg/directives/`, registered via `directives.DefaultRegistry()`.

## Directive catalog

| Directive | Stage | Effect |
|-----------|-------|--------|
| `@doc` | ignore | Documentation string on block/parameter |
| `@input` | resolve | Declares runtime input parameter for `Evaluate()` |
| `@retry` | init | Retry failed block up to `limit` times (`-1` = unlimited) |
| `@timeout` | init | Maximum duration for block execution |
| `@skip` | init | Skip block execution |
| `@triggers` | resolve | Run block only when listed block IDs complete |
| `@run` | — | Control run behavior (see `pkg/directives/run.go`) |
| `@todo` | — | Mark unfinished |
| `@bug` | — | Mark known bug |
| `@deprecated` | — | Deprecation marker |
| `@output` | — | Output metadata |
| `@string`, `@integer`, … | — | Inline type directives for `@input` |

Schema-only directives (`@required`, `@format`, …) are documented in [Schema annotations](../reference/schema-annotations.md).

## Examples

### Runtime input

`examples/inputs/main.cf`:

```conflow
@input {
  type:string
}
name := "You"

println {
  value = "Hello " + main.name + "!"
}
```

Passed to `conflow.Evaluate(..., inputParams map[ID]interface{})`.

### Retry

`examples/retry/main.cf` — `@retry` on a failing block with eventual success.

### Triggers

`examples/triggers/main.cf`:

```conflow
@triggers ["sleep2"]
println {
    value = str_format("%d %d", sleep1.i1, sleep2.i2)
}
```

`println` runs when `sleep2` completes (not only when all upstream blocks complete).

### Documentation

`examples/jsonschema/main.cf`:

```conflow
@doc "My schema was defined in the person.json file"
```

## Applying directives in Go

Runtime directives implement:

```go
ApplyToRuntimeConfig(*conflow.RuntimeConfig)
```

Parameter directives implement:

```go
ApplyToParameterConfig(*conflow.ParameterConfig)
```

`RuntimeConfig` fields: `Skip`, `Timeout`, `Triggers`, `RetryConfig` (`pkg/conflow/block_directive.go`).

## Registry

Default runtime directive registry (`pkg/directives/directives.go`):

`array`, `boolean`, `bug`, `deprecated`, `doc`, `input`, `integer`, `map`, `number`, `output`, `retry`, `run`, `skip`, `string`, `timeout`, `todo`, `triggers`

## See also

- [Lifecycle and evaluation stages](../concepts/lifecycle-and-stages.md)
- [Embedding and Evaluate API](../runtime/embedding.md)
