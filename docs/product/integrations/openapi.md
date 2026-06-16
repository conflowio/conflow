---
title: OpenAPI integration
summary: Authoring OpenAPI 3 definitions in Conflow and generating Go, JSON, or YAML.
parent: integrations/index.md
keywords: [openapi, api, petstore]
---

# OpenAPI integration

## Overview

OpenAPI 3 API definitions can be written as **configuration blocks** in Conflow, then exported via the CLI.

Packages:

- `pkg/openapi/` — block types for paths, operations, parameters, responses, components
- `pkg/openapi/generator/` — code generation templates
- `cmd/conflow/openapi/` — CLI wiring

## Authoring example

Full pet store API: `examples/openapi/petstore.cf`

Structure:

```conflow
openapi = "3.0.0"

info {
  version = "1.0.0"
  title = "Swagger Petstore"
  ...
}

server {
  url = "http://petstore.swagger.io/api"
}

path "/pets" {
  get { ... }
  post { ... }
}

schema:object "NewPet" { ... }
schema:all_of "Pet" { ... }
```

Uses `schema:ref`, `content "application/json"`, parameter blocks, response codes.

## Code generation

```bash
conflow openapi generate go    # Go server stubs / types
conflow openapi generate json  # OpenAPI JSON document
conflow openapi generate yaml  # OpenAPI YAML document
```

Commands: `cmd/conflow/openapi/generate/`

Example Makefile: `examples/openapi/Makefile`

## Block type

OpenAPI structs use `@block "configuration"` — structured trees evaluated at parse/configuration stage.

Example: `pkg/openapi/operation.go`, `pkg/openapi/contact.go`.

## Relationship to JSON Schema

Component schemas in OpenAPI use the same `pkg/schema` types and `schema:*` Conflow syntax.

## See also

- [JSON Schema integration](./json-schema.md)
- [CLI](./../runtime/cli.md)
- [Syntax overview](../language/syntax-overview.md)
