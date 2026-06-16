---
title: JSON Schema integration
summary: Using JSON Schema files and schema blocks for validated structured configuration.
parent: integrations/index.md
keywords: [json schema, validation]
---

# JSON Schema integration

## Schema package

`pkg/schema/` implements a JSON Schema–compatible type system:

- Types: `string`, `integer`, `number`, `boolean`, `array`, `object`, `ref`, `allOf`, `oneOf`, …
- Validation: `ValidateValue`, `ValidateSchema`
- Formats: `pkg/schema/formats/` (date, email, hostname, uuid, …)

Identifiers must match `NameRegExp` (`[a-z][a-z0-9]*(?:_[a-z0-9]+)*`).

## Loading external JSON Schema

`examples/jsonschema/`:

- External file `person.json` defines person/spouse/pet types
- Conflow program instantiates validated objects:

```conflow
@doc "My schema was defined in the person.json file"
you person {
    name = "You"
    you_spouse spouse { name = "Charlie" }
    pet1 pet { name = "Hunter" }
}
```

## Schema in Conflow source

Inline schema blocks use typed shorthand:

```conflow
schema:object "Name" {
    property:string "field"
    required = ["field"]
}
```

Schema directives registry: `pkg/schema/directives/directives.go`.

## Block schema from Go

Generated block interpreters expose `Schema()` from struct field annotations — same validation path as JSON Schema constraints.

## Validation timing

- Parameter assignment during parse/evaluate
- `Evaluate()` input parameter validation for `@input`
- `schema.Schema.Validate(ctx)` for standalone schema trees

## See also

- [Schema annotations](../reference/schema-annotations.md)
- [OpenAPI](./openapi.md) — OpenAPI schemas reuse the same model
