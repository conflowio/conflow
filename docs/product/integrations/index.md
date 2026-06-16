---
title: Integrations
summary: Index for JSON Schema and OpenAPI integrations in Conflow.
parent: index.md
keywords: [integrations, json schema, openapi]
---

# Integrations

Conflow's schema layer aligns with **JSON Schema**, enabling validation and external schema import. **OpenAPI 3** definitions can be authored in Conflow and exported to Go, JSON, or YAML.

## Topics

| Document | Summary |
|----------|---------|
| [JSON Schema](./json-schema.md) | Loading schemas, validation in workflows |
| [OpenAPI](./openapi.md) | Authoring APIs, code generation |

## Packages

| Package | Role |
|---------|------|
| `pkg/schema/` | Type system, validation, formats |
| `pkg/schema/directives/` | Schema annotation interpreters |
| `pkg/openapi/` | OpenAPI configuration blocks |
| `pkg/openapi/generator/` | Export templates (Go, JSON, YAML) |

## See also

- [Schema annotations](../reference/schema-annotations.md)
- [Value types](../language/value-types.md)
