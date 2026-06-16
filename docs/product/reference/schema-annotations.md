---
title: Schema annotations
summary: Complete list of schema and field directives from pkg/schema/directives.
parent: reference/index.md
keywords: [schema, annotations, required, format]
---

# Schema annotations

Schema annotations are Go comment directives on struct fields (and some directive blocks). They populate `schema.Schema` metadata used for validation and code generation.

Registry: `pkg/schema/directives/directives.go`

## Block and function metadata

| Directive | Purpose |
|-----------|---------|
| `@block` | Mark struct as block; set `type`, `eval_stage` |
| `@function` | Mark Go function as Conflow callable |

## Field identity and exposure

| Directive | Purpose |
|-----------|---------|
| `@id` | Block ID field |
| `@name` | Override JSON/parameter name |
| `@ignore` | Exclude from Conflow surface |
| `@value` | Value parameter for short syntax |
| `@key` | Map key binding |

## Input/output semantics

| Directive | Purpose |
|-----------|---------|
| `@required` | Required input |
| `@read_only` | Output parameter |
| `@write_only` | Write-only input |
| `@generated` | Generated child block field |
| `@dependency` | Runtime-injected dependency |

## Evaluation

| Directive | Purpose |
|-----------|---------|
| `@eval_stage` | Stage when field is evaluated (`init`, `main`, …) |
| `@default` | Default value if omitted |

## JSON Schema constraints

| Directive | Purpose |
|-----------|---------|
| `@const` | Constant value |
| `@enum` | Enumeration |
| `@minimum`, `@maximum` | Numeric bounds |
| `@exclusive_minimum`, `@exclusive_maximum` | Exclusive bounds |
| `@multiple_of` | Numeric multiple |
| `@min_length`, `@max_length` | String length |
| `@pattern` | Regex pattern |
| `@format` | String format (see `pkg/schema/formats/`) |
| `@min_items`, `@max_items` | Array length |
| `@unique_items` | Unique array elements |
| `@min_properties`, `@max_properties` | Object property count |
| `@dependent_required` | Conditional required fields |
| `@one_of` | Union type |

## Documentation metadata

| Directive | Purpose |
|-----------|---------|
| `@title` | Short title |
| `@deprecated` | Deprecation flag |
| `@examples` | Example values |
| `@result_type` | Result typing hint |

## String formats (via @format)

Implemented in `pkg/schema/formats/`: includes `date`, `date-time`, `time`, `email`, `hostname`, `ipv4`, `ipv6`, `uri`, `uuid`, `byte`, and others.

## Runtime vs schema directives

| Layer | Location | Examples |
|-------|----------|----------|
| Schema (Go comments on fields) | `pkg/schema/directives/` | `@required`, `@format` |
| Runtime (Conflow `@` in source) | `pkg/directives/` | `@retry`, `@input` |

## See also

- [Annotations reference](../extending/annotations.md)
- [JSON Schema integration](../integrations/json-schema.md)
- [Runtime directives](../language/runtime-directives.md)
