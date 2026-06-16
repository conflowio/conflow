# Data isolation — Batch 2: Wire bind at runtime boundaries

**Design Doc:** [2026-06-16-data-isolation-design.md](./2026-06-16-data-isolation-design.md)

**Goal:** Call `bind.BindValue` at the three runtime cross-block boundaries (`setChild`, `@input` in `Evaluate`, variable reads) so mutable collections cannot alias across blocks.

**Scope:**

- Included: `container.go` (`setChild`), `eval.go` (`@input`), `variable/node.go` (defensive read bind), isolation tests at each wiring point
- Excluded: codegen `SetParam` template (Batch 4), literal eval → immutable types (Batch 3), `static_container.go`, product docs

---

### Task 1: Bind in `setChild` before parameter assignment

**Files:**

- Modify: `pkg/conflow/block/container.go:491-522`
- Test: `pkg/conflow/block/isolation_test.go` (package `block`, whitebox)

**Step 1: Write the failing test**

In `isolation_test.go`, build a minimal `Container` with a fake interpreter whose schema defines an array property. Create a `parameter.Container` holding `[]interface{}{"shared"}`. Call `setChild` and assert `SetParam` received a bound value (`*values.List`) that does not alias the upstream slice.

Also test the `extraParams` path (property not in schema) uses `schema.AnyValue()` for bind.

**Step 2: Run test to verify it fails**

Run: `go test ./pkg/conflow/block/... -v -run Isolation`

Expected: FAIL — `SetParam` receives the raw slice pointer (aliasing).

**Step 3: Implement bind in `setChild`**

Import `github.com/conflowio/conflow/pkg/conflow/bind`.

Before `SetParam` and before `extraParams` assignment:

```go
prop, inSchema := s.(*schema.Object).PropertyByParameterName(string(name))
var propSchema schema.Schema
if inSchema {
    propSchema = prop
} else {
    propSchema = schema.AnyValue()
}
bound, bindErr := bind.BindValue(propSchema, value)
if bindErr != nil {
    return parsley.NewError(r.Node().Pos(), bindErr)
}
value = bound
```

**Step 4: Run test to verify it passes**

Run: `go test ./pkg/conflow/block/... -v -run Isolation`

Expected: PASS

---

### Task 2: Bind `@input` params in `Evaluate` (after validation)

**Files:**

- Modify: `pkg/conflow/eval.go:44-55`
- Test: `pkg/conflow/eval_isolation_test.go`

**Step 1: Write the failing test**

Parse `foo testblock { field_array }` with a test interpreter whose schema marks `FieldArray` as `@input` (`annotations.UserDefined`). Call `conflow.Evaluate` with `map[ID]interface{}{"field_array": []interface{}{"a"}}`. Read the block's `field_array` param; assert it is a bound `*values.List[interface{}]` and mutating the original input slice does not change the stored value.

**Step 2: Run test to verify it fails**

Run: `go test ./pkg/conflow/... -v -run "Evaluate input isolation"`

Expected: FAIL — param aliases caller slice.

**Step 3: Implement bind after `ValidateValue`**

Current `Array.ValidateValue` / `Map.ValidateValue` only accept Go slices/maps (Batch 3 will add immutable types). Bind **after** validation so schema checks still pass:

```go
nv, err := property.ValidateValue(v)
if err != nil {
    return nil, fmt.Errorf("invalid input parameter %q: %w", k, err)
}
bound, err := bind.BindValue(property, nv)
if err != nil {
    return nil, fmt.Errorf("invalid input parameter %q: %w", k, err)
}
inputParams[k] = bound
```

**Step 4: Run test to verify it passes**

Run: `go test ./pkg/conflow/... -v -run "Evaluate input isolation"`

Expected: PASS

---

### Task 3: Defensive bind on variable cross-block read

**Files:**

- Modify: `pkg/conflow/variable/node.go:90-97`
- Test: `pkg/conflow/variable/isolation_test.go`

**Step 1: Write the failing test**

Ginkgo test: fake `BlockContainer.Param` returns `[]interface{}{"x"}`. Evaluate `variable.Node.Value()` with array schema set in `StaticCheck`. Assert result is bound (not identical to upstream slice).

**Step 2: Run test to verify it fails**

Run: `go test ./pkg/conflow/variable/... -v -run Isolation`

Expected: FAIL — returns raw slice from `Param`.

**Step 3: Implement defensive bind**

```go
raw := blockContainer.Param(n.paramNameNode.ID())
if n.schema == nil {
    return raw, nil
}
bound, err := bind.BindValue(n.schema, raw)
if err != nil {
    return nil, parsley.NewError(n.Pos(), err)
}
return bound, nil
```

**Step 4: Run test to verify it passes**

Run: `go test ./pkg/conflow/variable/... -v -run Isolation`

Expected: PASS

---

### Task 4: Regression suite

**Files:**

- (no new files)

**Step 1: Run conflow package tests**

Run: `go test ./pkg/conflow/... -v`

Expected: PASS (existing variable parser test still passes for string params).

**Step 2: Run full test suite**

Run: `go test ./...`

Expected: PASS

---

## Review

Before starting the tasks above, record the baseline SHA: `git rev-parse HEAD`

After completing all tasks:

1. Launch a review subagent:
   - **Quick:** use the `review-bugbot` skill on uncommitted changes, or
   - **Thorough:** use the Task tool with the `code-review` skill, passing this plan file path, baseline SHA, and design doc reference
2. Use the `receiving-code-review` skill to evaluate and address the feedback
3. Verify all fixes pass: `go test ./pkg/conflow/...` and `go test ./...`
4. Commits may happen only when the user explicitly requests them, and only after steps 1–3 are complete

## Batch 3 blockers / notes

- `Array.ValidateValue` / `Map.ValidateValue` still require `[]interface{}` / `map[string]interface{}` — `@input` bind must stay **after** validation until Batch 3.
- Literal eval still produces mutable slices/maps — bound copies at boundaries are correct but slower until Batch 3 emits `*values.List` / `*values.Map`.
- Generated `SetParam` in block interpreters still assigns directly (Batch 4); `setChild` bind is the primary guard for evaluated parameter wiring.
- `static_container.go` also calls `SetParam` without bind — out of Batch 2 scope; consider in a follow-up if static eval paths need isolation.
