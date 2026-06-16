---
title: Syntax overview
summary: Conflow surface syntax for blocks, parameters, comments, modules, and schema shorthand.
parent: language/index.md
keywords: [syntax, blocks, modules, comments]
---

# Syntax overview

## Comments

```conflow
// line comment
```

## User-defined parameters

Assignment with `:=` defines a parameter on the **current block** (usually `main`):

```conflow
program := "test.sh"
colors := ["green", "red", "blue"]
```

## Block invocation

Full form:

```conflow
block_id block_type {
  param_name = expression
  child_id child_type {
    ...
  }
}
```

Value-parameter shorthand (when block type has `@value`):

```conflow
println "Hello World"
print "Result: " + value
```

## Parameter references

```conflow
main.program          // parameter on main
test.stdout           // output on block test
colors[i1.value]      // indexing
```

## Block labels on expressions

Some blocks bind a generated value into a local parameter:

```conflow
sleep1 sleep {
    i1 := i1.value
    duration = 50ms
}
```

## Modules and imports

```conflow
sum import "./sum"

one_plus_two sum {
  a = 1
  b = 2
}
```

`import` is a parse-stage task block (`pkg/blocks/import.go`) that loads a directory's `main` module and exposes it as a block type named by the block ID (`sum`).

See `examples/modules/`.

## Multi-file programs

`ParseFile` / `ParseDir` can load multiple `.cf` files into one parse context (`examples/multifile`).

## Schema shorthand (OpenAPI / JSON Schema)

Typed schema blocks use colon syntax:

```conflow
schema:object "Pet" {
    property:string "name"
    required = ["name"]
}

schema:array {
    items:ref "#/components/schemas/Pet"
}

parameter {
    schema:integer {
        format = "int32"
    }
}
```

Used in OpenAPI examples (`examples/openapi/petstore.cf`).

## Multiline strings

Triple-quoted strings for long text (e.g. OpenAPI descriptions):

```conflow
description = """
  Long description
  across lines
"""
```

## Directives on blocks

Prefix with `@`:

```conflow
@doc "My schema was defined in the person.json file"
@triggers ["sleep2"]
@input {
  type:string
}
```

See [Runtime directives](./runtime-directives.md).

## Root block

There is no explicit `main { }` wrapper in source — the file **is** the main body:

```conflow
hello {
  to = "World"
}
```

## Notation reference

README points to Go's EBNF notation for future formal grammar. Parser tests in `pkg/parsers/*_test.go` are practical examples of accepted syntax.

## See also

- [Value types and expressions](./value-types.md)
- [Modules and imports](../reference/built-in-blocks.md#import)
- [OpenAPI integration](../integrations/openapi.md)
