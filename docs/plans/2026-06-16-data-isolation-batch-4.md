# Data isolation — Batch 4: Codegen SetParam bind + regen

**Design Doc:** [2026-06-16-data-isolation-design.md](./2026-06-16-data-isolation-design.md)

**Goal:** Generated block interpreters call `bind.BindValue` in `SetParam` before assigning to struct fields, with `values.AsInterfaceSlice` / `AsStringInterfaceMap` for `[]interface{}` and `map[string]interface{}` fields. Regenerate all affected `.cf.go` files.

**Scope:**

- Included: `pkg/conflow/block/generator/interpreter_template.go`, `pkg/conflow/block/generator/funcs.go`, `pkg/schema/array.go`, `pkg/schema/map.go`, `pkg/schema/any.go` (`AssignValue` post-bind assignment), regen via `make generate`, generator/isolation test
- Excluded: function codegen `AssignValue` (Batch 3 policy keeps `ValidateValue` returning slices), schema block interpreters (`pkg/schema/interpreters/*.cf.go`), `static_container.go` bind wiring, product docs (Batch 5)

**Batch 2–3 context:** Runtime bind at `setChild`, `@input`, variable reads. Literals emit `*values.List` / `*values.Map`. `AssignValue` must accept bound values when converting to Go slice/map field types.

---

### Task 1: Add `bindAndAssignValue` template helper

**Files:**

- Modify: `pkg/conflow/block/generator/funcs.go`

**Step 1: Implement helper**

Generate per-parameter SetParam body:

```go
propSchema, _ := i.Schema().(*schema.Object).PropertyByParameterName("field_array")
bound, err := bind.BindValue(propSchema, value)
if err != nil {
    return err
}
// existing AssignValue using "bound" as source
```

Uses runtime schema from `i.Schema()` — no duplicated inline schema literals.

**Step 2: Wire template**

In `interpreter_template.go` `SetParam` cases, replace `assignValue` with `bindAndAssignValue`.

---

### Task 2: Implement `bindAndAssignValue` in block generator

**Files:**

- Modify: `pkg/conflow/block/generator/funcs.go`

**Step 1: Per-parameter SetParam body**

For each input property:

1. Resolve schema: `i.Schema().(*schema.Object).PropertyByParameterName(paramName)`
2. `bound, err := bind.BindValue(propSchema, value)`
3. Assign by schema type:
   - **Array / map:** `values.AsInterfaceSlice` / `AsStringInterfaceMap`, then assign (typed items: loop with existing item `AssignValue`)
   - **Any:** slice / map / scalar fallback (same as `@value` field)
   - **Scalars:** existing `schema.AssignValue` on `bound`

**Note:** Do not change `schema.AssignValue` — function codegen embeds it in `var x = …` expressions and cannot use multi-statement bind/normalize blocks.

---

### Task 3: Regenerate interpreters

**Step 1: Build and generate**

```bash
make generate
```

Regenerates all `*.cf.go` under `--local github.com/conflowio/conflow`.

**Step 2: Spot-check**

`pkg/test/block.cf.go` `SetParam` imports `bind`, calls `bind.BindValue` per input param, then assigns via updated `AssignValue` helpers.

---

### Task 4: SetParam isolation test

**Files:**

- Create: `pkg/conflow/block/generator/setparam_isolation_test.go`

**Step 1: Test generated interpreter path**

Use `test.BlockInterpreter` (via `_ "github.com/conflowio/conflow/pkg/test"`). Call `SetParam` with mutable `[]interface{}{"shared"}`. Assert `FieldArray` is a distinct slice; mutating upstream does not change stored field.

Also cover `field_map` with `map[string]interface{}`.

**Step 2: Run**

```bash
go test ./pkg/conflow/block/generator/... -v -run SetParam
go test ./... -count=1
```

Expected: PASS

---

## Review

Before starting the tasks above, record the baseline SHA: `git rev-parse HEAD`

After completing all tasks:

1. Launch a review subagent on uncommitted changes
2. Verify: `go test ./... -count=1`
3. Do not commit unless the user requests

## Batch 5 preview

- Product docs: value semantics in `docs/product/concepts/blocks-and-parameters.md`
- Examples walkthrough
- Optional `CONFLOW_BIND_DEBUG` logging
- Consider `static_container.go` bind wiring
- Optional: function codegen `AssignValue` + reorder `@input` validate/bind
