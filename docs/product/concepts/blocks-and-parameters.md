---
title: Blocks and parameters
summary: Block structure, parameter value forms, IDs, and references between blocks.
parent: concepts/index.md
keywords: [blocks, parameters, references, value parameter]
---

# Blocks and parameters

## Blocks

A **block** is a typed node in the workflow graph. In source text it appears as:

```conflow
block_id block_type {
  parameter = value
  child_id child_type { ... }
}
```

Or with a **value parameter** shorthand when the block type marks one field as `@value`:

```conflow
print "Hello " + name
```

### Block identity

- Blocks may have a **globally unique ID** within a parse context (`test`, `hello_stdout`).
- Parameters on other blocks are referenced as **`block_id.parameter_name`** (e.g. `main.program`, `test.stdout`).
- The root block is always **`main`**; user-defined parameters on main use `main.name` when referenced from nested blocks.

### Block body

The block body is **optional**. Generator output blocks (e.g. `it` from `iterator`) often have no body.

## Parameter value forms

A parameter value can be:

| Form | Example |
|------|---------|
| Literal | `"hello"`, `42`, `true`, `["a","b"]` |
| User-defined parameter | `program := "test.sh"` (defines on current block) |
| Parameter reference | `main.program`, `bar.p1` |
| Function call | `len(colors)`, `str_format("Hi %s", name)` |
| Expression | `1 + 2`, `"a" + "b"`, boolean logic |
| Array/map element | `colors[i1.value]` |
| Typed schema shorthand | `schema:object { ... }`, `items:string` |

## Parameter kinds (Go field annotations)

| Annotation | Meaning |
|------------|---------|
| (none) | Input parameter, set from Conflow |
| `@read_only` | Output; set by block logic, readable in Conflow |
| `@generated` | Child block type emitted by generator |
| `@dependency` | Injected by runtime (e.g. `BlockPublisher`, `io.Writer`) |
| `@ignore` | Not exposed to Conflow |
| `@id` | Block instance ID field |
| `@value` | Enables short block syntax |
| `@input` (directive) | Exposed as runtime input to `Evaluate()` |

Schema validation annotations (`@required`, `@minimum`, `@format`, …) apply to input parameters. See [Schema annotations](../reference/schema-annotations.md).

## Global block registry

Blocks are registered in a **global parse context**. Any named block's parameters can be referenced from anywhere in the program (not only children). This enables wiring across the graph without nesting.

## Child blocks vs parameters

- **Child blocks** are nested block declarations inside a parent block body.
- They participate in the dependency graph as nodes.
- **Generated** children are declared in the generator's body but **instances** are created at runtime by `PublishBlock`.

## Example (from README)

```conflow
program := "test.sh"

test exec {
    program = main.program
}

print "Result was: " + test.stdout
```

Here `program` is a user parameter on `main`, `test` is an `exec` task block, and `print` uses the value-parameter form.

## Implementation notes

- `BlockInterpreter` (`pkg/conflow/block.go`) — `SetParam`, `SetBlock`, `Schema()`, `CreateBlock`
- User-defined parameters are tagged in schema via `annotations.UserDefined`
- `pkg/conflow/parameter/` — containers and transforms

## See also

- [Block types](./block-types.md)
- [Dependencies and evaluation order](./dependencies-and-order.md)
- [Built-in blocks](../reference/built-in-blocks.md)
