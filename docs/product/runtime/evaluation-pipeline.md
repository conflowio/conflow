---
title: Evaluation pipeline
summary: Stages from parsing through dependency resolution to block Run and Close.
parent: runtime/index.md
keywords: [evaluation, pipeline, parse, resolve]
---

# Evaluation pipeline

## End-to-end flow

```text
1. Build ParseContext (ID registry, directive registry)
2. Register MainInterpreter + overrides
3. Parser.ParseFile / ParseDir  →  AST + block nodes
4. (Parse stage) import blocks load submodules
5. Dependency resolution per workflow region
6. For each block instance, per eval stage:
   a. Evaluate parameters whose dependencies are ready
   b. Run stage hooks (directives, Init)
   c. Schedule Run / Close on job scheduler
7. Generators PublishBlock → new child instances → repeat
8. Return result or aggregate errors
```

## Parse phase

- Conflow text → Parsley AST nodes (`pkg/parsers/`)
- Each block node holds `BlockInterpreter` reference and schema
- `import` blocks extend registry at parse stage (`pkg/blocks/import.go`)

## Resolve phase

Directives and metadata applied before main execution:

- `@triggers`, `@input` typing
- Schema cross-references

## Stage-driven evaluation

See [Lifecycle and evaluation stages](../concepts/lifecycle-and-stages.md).

`EvalContext` tracks current stage and schedules parameter evaluation accordingly.

## Dependency resolution

Before/during evaluation, `pkg/conflow/dependency/resolver.go` orders nodes and validates:

- All references resolve
- No illegal cycles (including generator cases)

## Block execution

1. `CreateBlock` on interpreter — allocate Go struct
2. `SetParam` / `SetBlock` — populate from evaluated AST
3. `Init` — optional skip
4. `Run` — business logic
5. `Close` — optional cleanup

## Error handling

- Parse errors — parser with source positions
- `TransformPathErrors` — maps evaluation errors to Conflow paths (`pkg/conflow/eval.go`)
- Block errors propagate to dependents per dependency rules

## Pub/sub

`EvalContext` includes `PubSub` for coordinating generator publications and dependent wakeups (`pkg/conflow/pubsub.go`).

## See also

- [Dependencies and evaluation order](../concepts/dependencies-and-order.md)
- [Job scheduler](./job-scheduler.md)
