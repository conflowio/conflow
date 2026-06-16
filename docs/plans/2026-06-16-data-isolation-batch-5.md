# Data isolation — Batch 5: Docs, examples, debug logging

**Design Doc:** [2026-06-16-data-isolation-design.md](./2026-06-16-data-isolation-design.md)

**Goal:** Document value isolation semantics for users, add illustrative examples (product docs), and optional `CONFLOW_BIND_DEBUG` logging at bind boundaries. Final batch — completes the data isolation initiative.

**Scope:**

- Included: `docs/product/concepts/blocks-and-parameters.md` (value semantics section), `pkg/conflow/bind/debug.go`, debug calls in `bind.go`, this plan file
- Excluded: new runnable example directory (doc-only Conflow snippets suffice), `static_container.go` bind wiring, function codegen `AssignValue`, commits

**Batch 1–4 context:** `pkg/values` immutable collections, `bind.BindValue`, runtime wiring (`setChild`, `@input`, variable reads), literal eval → `*values.List`/`*values.Map`, generated `SetParam` with bind.

---

### Task 1: `CONFLOW_BIND_DEBUG` logging

**Files:**

- Create: `pkg/conflow/bind/debug.go`
- Modify: `pkg/conflow/bind/bind.go`

**Step 1: Debug helper**

Follow `CONFLOW_LOG` env pattern from `examples/common/main.go` — check env once at init:

```go
var bindDebug = os.Getenv("CONFLOW_BIND_DEBUG") != ""

func debugBind(schemaType schema.Type, value interface{}) {
    if !bindDebug {
        return
    }
    fmt.Fprintf(os.Stderr, "conflow bind: schema=%s value=%s\n", schemaType, valueKind(value))
}
```

`valueKind` returns short labels: `nil`, `values.List`, `values.Map`, or `reflect.TypeOf(v).String()`.

**Step 2: Log at `BindValue` entry**

Call `debugBind(s.Type(), value)` at the start of `BindValue` (after nil check).

**Step 3: Verify**

Run: `go test ./pkg/conflow/bind/... -count=1`

Expected: PASS (debug off by default; no test pollution).

---

### Task 2: Product docs — value semantics

**Files:**

- Modify: `docs/product/concepts/blocks-and-parameters.md`

**Step 1: Add section after "Parameter value forms"**

New section **"Value semantics and isolation"** covering:

1. **Bind at boundaries** — values crossing block boundaries (`SetParam`, parameter references, `@input`) pass through `bind.BindValue`.
2. **Scalars** — copied by value (`bool`, `int64`, `float64`, `string`, …).
3. **Immutable collections** — Conflow literals `[...]` and `map{...}` produce frozen `*values.List` / `*values.Map`; bind shares the pointer (O(1), safe because immutable).
4. **Mutable Go collections** — `[]interface{}`, `map[string]interface{}`, typed slices/maps from Go `@input` or block outputs get deep-copied (or converted to immutable on first bind).
5. **Inside `Run()`** — block authors may mutate local copies; isolation applies at bind time, not inside block logic.
6. **Debug** — set `CONFLOW_BIND_DEBUG=1` to log bind entry (schema type + value kind) to stderr.

Include a Conflow example: shared literal list referenced by two blocks; explain they receive independent bound views if upstream was mutable, or shared immutable handle for literals.

Match existing doc tone: short paragraphs, tables where helpful, link to implementation notes.

**Step 2: Update Implementation notes**

Add `pkg/conflow/bind/` and `pkg/values/` to the implementation notes list.

---

### Task 3: Doc examples (no new example directory)

**Files:**

- Modify: `docs/product/concepts/blocks-and-parameters.md` (same file as Task 2)

**Step 1: Conflow walkthrough**

Add example block wiring:

```conflow
colors := ["red", "green", "blue"]

consumer_a print colors[0]
consumer_b print colors[1]
```

Explain: literal `colors` is an immutable list; parallel consumers bind safely. Contrast with Go `@input` passing a mutable slice (deep copy at bind).

Reference isolation tests: `pkg/conflow/bind/isolation_test.go`, `pkg/conflow/block/isolation_test.go`.

No new `examples/isolation/` — behavior is engine-internal and already covered by unit/integration tests; product doc is the user-facing surface.

---

## Review

Before starting the tasks above, record the baseline SHA: `git rev-parse HEAD`

After completing all tasks:

1. Verify: `go test ./... -count=1`
2. Spot-check docs render (front matter intact, section links valid)
3. Do not commit unless the user requests

## Initiative complete

After Batch 5 merges, all design deliverables are done:

| Batch | Status |
|-------|--------|
| 1 — `pkg/values` + bind core | Done |
| 2 — Runtime wiring | Done |
| 3 — Literal eval + schema validation | Done |
| 4 — Codegen SetParam + regen | Done |
| 5 — Docs + debug | This batch |
