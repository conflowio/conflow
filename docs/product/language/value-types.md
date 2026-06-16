---
title: Value types and expressions
summary: Literals, durations, collections, operators, function calls, and type checking.
parent: language/index.md
keywords: [types, expressions, literals, functions]
---

# Value types and expressions

## Literals

| Type | Examples |
|------|----------|
| String | `"hello"`, `"""multi"""` |
| Integer | `42`, `-1` |
| Number | `3.14` |
| Boolean | `true`, `false` |
| Duration | `50ms`, `1s` (where schema allows) |
| Null | Context-dependent |

## Collections

```conflow
["a", "b", "c"]
```

Maps and nested structures follow schema of the target parameter.

## Operators

Expressions support arithmetic, boolean logic, and concatenation where types allow (similar to Go):

```conflow
1 + 2
"a" + "b"
str_format("Hi %s", name)
len(colors)
```

Exact operator set is implemented in `pkg/parsers/expression.go` and related parsers.

## Function calls

Functions are invoked by name with parentheses:

```conflow
len(main.colors)
str_lower("ABC")
json_encode(obj)
```

Built-in registry: `pkg/functions/registry.go`. Custom functions via `// @function` in Go.

## Type checking

Types are **not** inferred freely in the Conflow file — they come from:

1. **Target parameter schema** on the block being configured
2. **Function signature** schema from generated function interpreters
3. **JSON Schema** constraints (`@minimum`, `@enum`, `@format`, …)

Invalid assignments fail at parse/resolve or at `Evaluate` input validation.

## Formats

String formats (email, date-time, hostname, ip, …) live in `pkg/schema/formats/` and are applied via `@format` on fields.

## JSON types in schema

Conflow's schema layer mirrors JSON Schema types:

- `string`, `integer`, `number`, `boolean`, `array`, `object`
- Combinators: `allOf`, `oneOf`, `ref`

## Common functions by category

| Category | Names |
|----------|-------|
| Common | `len`, `string` |
| Array | `arr_contains` |
| JSON | `json_decode`, `json_encode` |
| Math | `abs`, `ceil`, `floor`, `max`, `min`, `round`, `trunc` |
| Strings | `str_contains`, `str_format`, `str_has_prefix`, `str_has_suffix`, `str_join`, `str_lower`, `str_replace`, `str_split`, `str_title`, `str_trim_prefix`, `str_trim_space`, `str_trim_suffix`, `str_upper` |
| Time | `time_now` |

Full list: [Built-in functions](../reference/built-in-functions.md).

## See also

- [Schema annotations](../reference/schema-annotations.md)
- [Built-in functions](../reference/built-in-functions.md)
