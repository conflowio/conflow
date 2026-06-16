---
title: CLI reference
summary: conflow binary commands — version, generate, openapi.
parent: runtime/index.md
keywords: [cli, cobra, commands]
---

# CLI reference

Binary entry: `cmd/conflow/main.go`

```bash
conflow <command> [subcommand] [args]
```

## `conflow version`

Prints `conflow.Version` (`pkg/conflow/version.go`).

## `conflow generate`

Generates `*.cf.go` interpreters from `// @block` and `// @function` annotations.

```bash
conflow generate [path]
conflow generate --local github.com/myorg/myproject
```

Details: [Code generation workflow](../extending/codegen-workflow.md).

## `conflow openapi`

```bash
conflow openapi generate go
conflow openapi generate json
conflow openapi generate yaml
```

Generates artifacts from OpenAPI definitions authored in Conflow. See [OpenAPI integration](../integrations/openapi.md).

## Logging

When stdin is a terminal, CLI uses zerolog console writer with timestamps.

## Context

Root command attaches logger to context and handles `SIGINT` for cancellation.

## Building the CLI

```bash
go build -o conflow ./cmd/conflow
```

Examples assume `conflow` is on `PATH` or use `go run` / Make targets.

## See also

- [Code generation workflow](../extending/codegen-workflow.md)
- [OpenAPI integration](../integrations/openapi.md)
