---
title: Built-in functions
summary: Callable functions registered in pkg/functions/registry.go.
parent: reference/index.md
keywords: [built-in functions, str_, json_]
---

# Built-in functions

Registry: `pkg/functions/registry.go` (`functions.DefaultRegistry()`).

Functions are invoked in expressions: `len(arr)`, `str_format("x", y)`.

## Common

| Name | Go implementation | Package |
|------|-------------------|---------|
| `len` | `Len` | `pkg/functions/len.go` |
| `string` | `String` | `pkg/functions/string.go` |

## Array

| Name | Implementation |
|------|----------------|
| `arr_contains` | `array.Contains` |

## JSON

| Name | Implementation |
|------|----------------|
| `json_decode` | `json.Decode` |
| `json_encode` | `json.Encode` |

## Math

| Name | Implementation |
|------|----------------|
| `abs` | `math.Abs` |
| `ceil` | `math.Ceil` |
| `floor` | `math.Floor` |
| `max` | `math.Max` |
| `min` | `math.Min` |
| `round` | `math.Round` |
| `trunc` | `math.Trunc` |

## Strings

| Name | Implementation |
|------|----------------|
| `str_contains` | `strings.Contains` |
| `str_format` | `strings.Format` |
| `str_has_prefix` | `strings.HasPrefix` |
| `str_has_suffix` | `strings.HasSuffix` |
| `str_join` | `strings.Join` |
| `str_lower` | `strings.Lower` |
| `str_replace` | `strings.Replace` |
| `str_split` | `strings.Split` |
| `str_title` | `strings.Title` |
| `str_trim_prefix` | `strings.TrimPrefix` |
| `str_trim_space` | `strings.TrimSpace` |
| `str_trim_suffix` | `strings.TrimSuffix` |
| `str_upper` | `strings.Upper` |

## Time

| Name | Implementation |
|------|----------------|
| `time_now` | `time.Now` |

## Adding custom functions

See [Defining functions in Go](../extending/go-functions.md).

## See also

- [Value types and expressions](../language/value-types.md)
