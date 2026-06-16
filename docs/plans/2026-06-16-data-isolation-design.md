# Data isolation — design document

**Status:** Draft  
**Date:** 2026-06-16

## Problem

Conflow blocks run concurrently when dependencies allow. Parameter values cross block boundaries via references (`other_block.output`), `@input`, and generator publications. Today values are assigned as plain Go values (`interface{}`, `[]interface{}`, `map[string]interface{}`) with no copy policy — mutable slices and maps can alias across blocks.

## Goal

Enforce **data isolation** at cross-block boundaries:

1. Block A cannot mutate data visible to block B.
2. Sharing large collections is **O(1)** when safe (immutable handles).
3. Block authors can use normal Go inside `Run()`; isolation applies at **bind** time (inputs, outputs, references).
4. Non-immutable types still work via **deep copy** (correct but slower).

## Non-goals (this initiative)

- Untrusted plugin sandboxing (WASM / subprocess) — out of scope.
- Changing Conflow language syntax.
- Hard-rejecting `[]T` / `map[K]V` on block fields.
- Distributed workflow durability (Temporal-style).

## Design decisions

### D1: Bind at boundaries, not inside blocks

| Location | Policy |
|----------|--------|
| Inside `Run()` / `Init()` / `Close()` | Author's responsibility; any Go types |
| `SetParam` / `SetBlock` (incoming) | Always `BindValue(schema, value)` |
| `Container.Param` / variable references (outgoing read) | Return bound snapshot (or immutable handle copy) |
| `@read_only` publish to graph | Normalize once when value becomes visible |
| `@input` in `Evaluate()` | Bind before validation/assignment |

### D2: Two-tier value policy

| Value kind | Cross-block bind |
|------------|------------------|
| Scalars (`bool`, `int64`, `float64`, `string`, `time.Time`, `time.Duration`) | Copy by value (unchanged) |
| `*values.List[T]`, `*values.Map[K,V]` (frozen) | Copy pointer (O(1)); shared backing is safe |
| `[]T`, `map[string]T`, `[]interface{}`, `map[string]interface{}` | Deep copy (or convert-then-share — see D3) |
| Schema `object` (struct / nested map) | Schema-driven deep copy |
| Child `Block` references | Engine-owned; not data-copied |
| `io.Reader` / streams | Capability type; separate from value isolation |

### D3: Literal evaluation produces immutable collections

When Conflow evaluates `[...]` and `{...}` literals, produce `*values.List` / `*values.Map` instead of `[]interface{}` / `map[string]interface{}`. This makes the fast path the default for graph-native data.

Existing code that expects `[]interface{}` will be migrated incrementally (Batch 3+).

### D4: Author ergonomics

- Immutable types are the **efficient** path, not a hard requirement.
- Provide builders (`ListBuilder`, `MapBuilder`) and helpers (`ListOf`, `MapOf`, `FromSlice`).
- Codegen may later map schema `array`/`map` fields to immutable types; not required in batch 1.

### D5: Single bind entry point

All isolation logic lives in `pkg/conflow/bind`:

```go
func BindValue(s schema.Schema, value interface{}) (interface{}, error)
func BindValueForParam(s schema.Schema, value interface{}, mode BindMode) (interface{}, error)
```

Generated `SetParam` code calls `bind.BindValue` instead of direct assignment. No scattered copy logic.

### D6: `Clone()` interface for custom struct outputs

Types implementing `Clone() T` (via `util.SelfCloner`) get deep copy without reflection. Schema objects recurse through properties.

## Architecture

```text
Conflow literal / expression eval
        │
        ▼
   *values.List / *values.Map / scalars / objects
        │
        ▼
 variable reference (block.param) ──► Container.Param ──► bind (optional re-bind on read)
        │
        ▼
 SetParam / SetBlock ──► bind.BindValue(schema, value) ──► block struct field
        │
        ▼
   Run()  (private mutation OK)
        │
        ▼
 @read_only output ──► bind on publish (future: normalize outputs)
```

## Isolation points in codebase

| File | Change |
|------|--------|
| `pkg/values/` | New immutable `List`, `Map`, builders |
| `pkg/conflow/bind/` | Schema-driven `BindValue` |
| `pkg/util/cloner.go` | Extend for bind registry |
| `pkg/conflow/array.go`, `map.go` | Emit immutable types |
| `pkg/conflow/block/container.go` | Bind in `setChild` before `SetParam` |
| `pkg/conflow/variable/node.go` | Bind on cross-block read (defense in depth) |
| `pkg/conflow/eval.go` | Bind `@input` params |
| `pkg/conflow/block/generator/interpreter_template.go` | Generated `SetParam` uses bind |
| `pkg/schema/array.go`, `map.go` | `ValidateValue` accepts immutable types |
| `docs/product/concepts/blocks-and-parameters.md` | Document value semantics |

## Batch roadmap

| Batch | Deliverable | Mergeable alone? |
|-------|-------------|------------------|
| **1** | `pkg/values` + `pkg/conflow/bind` + tests | Yes |
| **2** | Wire bind into `setChild`, `Evaluate` inputs, variable reads | Yes (runtime guarded) |
| **3** | Literal eval → immutable types; schema `ValidateValue` updates | Yes |
| **4** | Codegen `SetParam` template + `conflow generate` regen | Yes |
| **5** | Examples, product docs, optional `CONFLOW_BIND_DEBUG` logging | Yes |

Each batch has its own plan file in `docs/plans/`.

## Verification strategy

- Unit tests: bind policy per schema type, aliasing detection (pointer inequality after bind).
- Integration test: two parallel blocks receive same upstream list; one mutates local copy; other unchanged.
- Regression: existing `go test ./...` green after each batch.

## Open questions (resolve during implementation)

1. **Re-bind on read:** Bind only at `SetParam`, or also when `variable.Node` reads `Container.Param`? Recommendation: bind at `SetParam` (primary); optional re-bind on read for `@read_only` fields returned by reference from struct pointers.
2. **Convert vs deep-copy for `[]interface{}`:** On bind, convert to `*values.List` once (then pointer-copy downstream) vs deep-copy slice each time. Recommendation: convert once at first bind.
3. **Go version:** Project is Go 1.20; avoid `iter.Seq` until go.mod is bumped.

## References

- Prior art: `vendor/github.com/conflowio/parsley/data/intset.go`, `intmap.go`
- Existing: `pkg/util/cloner.go`, `pkg/util/cloner_test.go`
- Product: `docs/product/concepts/blocks-and-parameters.md`
