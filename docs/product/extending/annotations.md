---
title: Annotations reference
summary: Comment annotation syntax for @block and struct field directives in Go source.
parent: extending/index.md
keywords: [annotations, @block, @required, directives]
---

# Annotations reference

Annotations are **Go comment directives** read by `conflow generate`. They are not Go struct tags.

## Block declaration

```go
// @block "task"
type MyBlock struct { ... }
```

Multi-line block metadata (type, eval_stage):

```go
//	@block {
//	  type = "directive"
//	  eval_stage = "init"
//	}
```

## Function declaration

```go
// @function
func MyFunc(x int64) int64 { ... }
```

## Common field annotations

| Annotation | Applies to | Meaning |
|------------|------------|---------|
| `@id` | Field | Block ID storage |
| `@required` | Input | Must be set in Conflow |
| `@read_only` | Field | Output parameter |
| `@write_only` | Field | Input-only in API schema sense |
| `@generated` | Field | Generated child block |
| `@dependency` | Field | Runtime injection |
| `@ignore` | Field | Hidden from Conflow |
| `@value` | Field | Value parameter for short syntax |
| `@default` | Field | Default if omitted |
| `@name` | Field | Override Conflow parameter name |
| `@eval_stage` | Field | When parameter is evaluated |
| `@key` | Field | Map key association |

## JSON Schema style (on fields)

| Annotation | Example |
|------------|---------|
| `@minimum`, `@maximum` | Numeric bounds |
| `@minLength`, `@maxLength` | String length |
| `@minItems`, `@maxItems` | Array size |
| `@pattern` | Regex |
| `@format` | `email`, `date-time`, … |
| `@enum` | Allowed values |
| `@const` | Fixed value |
| `@one_of` | Union types |

Full schema directive list: [Schema annotations](../reference/schema-annotations.md).

## Runtime directives in Conflow (not Go comments)

In `.cf` files: `@retry`, `@input`, `@doc`, … — see [Runtime directives](../language/runtime-directives.md).

## Parser location

- Block/function scan: `pkg/conflow/generator/generate.go`
- Metadata parsing: `pkg/conflow/generator/parser/metadata.go`
- Schema directives: `pkg/schema/directives/`

## See also

- [Schema annotations](../reference/schema-annotations.md)
- [Go blocks](./go-blocks.md)
