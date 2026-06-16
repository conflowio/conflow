# Data isolation — Batch 1: Immutable values and bind core

**Design Doc:** [2026-06-16-data-isolation-design.md](./2026-06-16-data-isolation-design.md)

**Goal:** Add immutable `List`/`Map` types and a schema-driven `BindValue` function with deep-copy fallback — tested in isolation, no runtime wiring yet.

**Scope:**

- Included: `pkg/values`, `pkg/conflow/bind`, extended cloner helpers, unit tests
- Excluded: `SetParam` integration, codegen changes, literal eval changes, docs (later batches)

---

### Task 1: Immutable `List` and `ListBuilder`

**Files:**

- Create: `pkg/values/list.go`
- Create: `pkg/values/list_test.go`
- Create: `pkg/values/values_suite_test.go`

**Step 1: Write the skeleton**

```go
// pkg/values/list.go
package values

type List[T any] struct {
    elems []T
}

func NewList[T any](elems ...T) *List[T] { panic("not implemented") }
func (l *List[T]) Len() int { return 0 }
func (l *List[T]) At(i int) T { var zero T; return zero }
func (l *List[T]) Elems() []T { return nil }

type ListBuilder[T any] struct{}

func NewListBuilder[T any]() *ListBuilder[T] { return &ListBuilder[T]{} }
func (b *ListBuilder[T]) Append(v T) {}
func (b *ListBuilder[T]) Freeze() *List[T] { return nil }
```

```go
// pkg/values/values_suite_test.go
package values_test

import (
    "testing"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestValues(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Values Suite")
}
```

**Step 2: Write the failing tests**

In `list_test.go` (Ginkgo):

- `NewList` copies input slice (mutating original does not affect list)
- `List` has no exported field access to backing slice via public API
- `ListBuilder.Freeze` returns independent immutable list
- Two blocks scenario: share pointer from `Freeze()`, append to original builder slice does not affect frozen list

**Step 3: Run test to verify it fails**

Run: `go test ./pkg/values/... -v -run List`

Expected: FAIL with assertion / panic, not import errors.

**Step 4: Implement**

```go
func NewList[T any](elems ...T) *List[T] {
    cp := make([]T, len(elems))
    copy(cp, elems)
    return &List[T]{elems: cp}
}

func NewListFromSlice[T any](s []T) *List[T] {
    cp := make([]T, len(s))
    copy(cp, s)
    return &List[T]{elems: cp}
}

func (l *List[T]) Len() int { return len(l.elems) }

func (l *List[T]) At(i int) T {
    if i < 0 || i >= len(l.elems) {
        panic(fmt.Sprintf("values.List.At: index %d out of range [0,%d)", i, len(l.elems)))
    }
    return l.elems[i]
}

// Elems returns a copy for interop (json, stdlib). Callers may mutate the copy.
func (l *List[T]) Elems() []T {
    cp := make([]T, len(l.elems))
    copy(cp, l.elems)
    return cp
}

func NewListBuilder[T any]() *ListBuilder[T] {
    return &ListBuilder[T]{elems: make([]T, 0)}
}

func (b *ListBuilder[T]) Append(v T) {
    b.elems = append(b.elems, v)
}

func (b *ListBuilder[T]) Freeze() *List[T] {
    return NewListFromSlice(b.elems)
}
```

Add helpers: `ListOf[T](...)` alias for `NewList`.

**Step 5: Run tests**

Run: `go test ./pkg/values/... -v`

Expected: PASS

---

### Task 2: Immutable `Map` and `MapBuilder`

**Files:**

- Create: `pkg/values/map.go`
- Create: `pkg/values/map_test.go`

**Step 1: Skeleton**

Same pattern as List — unexported `m map[K]V`, `Get`, `Keys`, `Len`, `MapBuilder` with `Set` + `Freeze`.

**Important:** `NewMapFromGoMap` must **clone** the input map (learn from Parsley `IntMap` footgun).

**Step 2: Failing tests**

- Mutating original Go map after `NewMapFromGoMap` does not affect immutable map
- `Freeze` produces independent map
- `Get` missing key returns `ok == false`

**Step 3–5:** Implement, verify PASS.

Run: `go test ./pkg/values/... -v`

---

### Task 3: Schema-driven deep copy in `pkg/conflow/bind`

**Files:**

- Create: `pkg/conflow/bind/bind.go`
- Create: `pkg/conflow/bind/bind_test.go`
- Create: `pkg/conflow/bind/bind_suite_test.go`
- Modify: `pkg/util/cloner.go` (add `DeepCopyValue(schema.Schema, interface{})` wrapper if needed)

**Step 1: Skeleton**

```go
// pkg/conflow/bind/bind.go
package bind

import "github.com/conflowio/conflow/pkg/schema"

func BindValue(s schema.Schema, value interface{}) (interface{}, error) {
    return value, nil // stub
}
```

**Step 2: Failing tests**

Use table-driven Ginkgo tests:

| Input | Schema | Expect |
|-------|--------|--------|
| `[]interface{}{"a","b"}` | `array` of string | New slice, different pointer |
| `*values.List[string]` | array of string | Same pointer returned |
| `map[string]interface{}{"k": 1}` | map | New map, different pointer |
| `*values.Map[string,int64]` | map | Same pointer |
| `int64(42)` | integer | Same value |
| nested `map[string]interface{}` | object schema | Deep equal, no shared inner maps |

Add test that mutating bind result does not affect source (for slice/map deep copy case).

**Step 3: Run failing tests**

Run: `go test ./pkg/conflow/bind/... -v`

**Step 4: Implement `BindValue`**

Dispatch on `schema.Schema` type:

```go
func BindValue(s schema.Schema, value interface{}) (interface{}, error) {
    if value == nil {
        return nil, nil
    }
    switch s.Type() {
    case schema.TypeBoolean, schema.TypeInteger, schema.TypeNumber, schema.TypeString:
        return value, nil // scalars
    case schema.TypeArray:
        return bindArray(s.(*schema.Array), value)
    case schema.TypeMap:
        return bindMap(s.(*schema.Map), value)
    case schema.TypeObject:
        return bindObject(s.(*schema.Object), value)
    default:
        return nil, fmt.Errorf("bind: unsupported schema type %s", s.Type())
    }
}
```

`bindArray`:

- `*values.List[T]` → return as-is (type assert via reflection or switch on `[]interface{}` wrapper types)
- `[]interface{}` → validate items recursively, build **new** slice OR convert to `*values.List[interface{}]` via builder (prefer convert-once per design D3)
- Typed `[]string` etc. → deep copy slice, copy each element if needed

`bindMap`: same pattern with `*values.Map`.

`bindObject`: recurse properties via `schema.Object.Properties`; support `map[string]interface{}`; for structs use reflection or existing `ValidateValue` path then deep-copy map representation.

Use `util.CloneValue` for scalars, compose cloners for collections.

**Step 5: PASS**

Run: `go test ./pkg/conflow/bind/... -v`

---

### Task 4: Convert helpers (slice/map ↔ immutable)

**Files:**

- Create: `pkg/values/convert.go`
- Create: `pkg/values/convert_test.go`

**Step 1: Skeleton**

```go
func FromSlice[T any](s []T) *List[T]
func FromInterfaceSlice(s []interface{}) (*List[interface{}], error)
func FromGoMap[K comparable, V any](m map[K]V) *Map[K, V]
func FromStringInterfaceMap(m map[string]interface{}) (*Map[string, interface{}], error)
```

**Step 2–5:** Tests + implement. Used by `bind` when normalizing `[]interface{}`.

Run: `go test ./pkg/values/... ./pkg/conflow/bind/... -v`

---

### Task 5: Aliasing regression test (bind integration)

**Files:**

- Create: `pkg/conflow/bind/isolation_test.go`

**Step 1: Write test (no new prod code)**

Simulate cross-block handoff:

```go
upstream := []interface{}{"shared"}
bound1, _ := bind.BindValue(arraySchema, upstream)
bound2, _ := bind.BindValue(arraySchema, upstream)

// bound1 and bound2 must not alias upstream or each other (for mutable slice input)
// mutate upstream after bind — bound values unchanged
```

For `*values.List` input, bound1 and bound2 may share pointer (same immutable handle).

**Step 2–3:** Run, ensure passes with Task 3 implementation.

Run: `go test ./pkg/conflow/bind/... -run Isolation -v`

---

## Review

Before starting the tasks above, record the baseline SHA: `git rev-parse HEAD`

After completing all tasks:

1. Launch a review subagent:
   - **Quick:** use the `review-bugbot` skill on uncommitted changes, or
   - **Thorough:** use the Task tool with the `code-review` skill, passing this plan file path, baseline SHA, and design doc reference
2. Use the `receiving-code-review` skill to evaluate and address the feedback
3. Verify all fixes pass: `go test ./pkg/values/... ./pkg/conflow/bind/...`
4. Commits may happen only when the user explicitly requests them, and only after steps 1–3 are complete

## Next batch preview

**Batch 2** wires `bind.BindValue` into:

- `pkg/conflow/block/container.go` (`setChild` before `SetParam`)
- `pkg/conflow/eval.go` (`@input` params)
- `pkg/conflow/variable/node.go` (defensive bind on cross-block read)

Do not start Batch 2 until Batch 1 is reviewed and merged.
