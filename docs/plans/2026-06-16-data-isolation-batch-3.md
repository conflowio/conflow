# Data isolation — Batch 3: Immutable literal eval and schema validation

**Design Doc:** [2026-06-16-data-isolation-design.md](./2026-06-16-data-isolation-design.md)

**Goal:** Conflow `[...]` and `map{...}` literals evaluate to `*values.List` / `*values.Map`, and schema `ValidateValue` accepts those types so downstream bind and validation paths work without converting back to mutable Go collections first.

**Scope:**

- Included: `pkg/conflow/array.go`, `pkg/conflow/map.go`, `pkg/schema/array.go`, `pkg/schema/map.go`, unit tests, test helper normalization for parser eval assertions, `pkg/functions/len.go` (hand-written function receiving validated values)
- Also required for green suite: `pkg/schema/any.go`, `pkg/schema/one_of.go`, `pkg/schema/schema.go` (`GetSchemaForValue`), `pkg/values/convert.go` (generic list/map detection), `pkg/functions/json/encode.go`
- Excluded: codegen `SetParam` template (Batch 4), function codegen `AssignValue` template, product docs, `static_container.go` bind wiring

**Batch 2 context:** Bind is wired at `setChild`, `@input`, and variable reads. `ValidateValue` still required `[]interface{}` / `map[string]interface{}`; `@input` validates then binds. Compatibility helpers live in `pkg/values/convert.go`.

**ValidateValue return policy:** Accept immutable input but **return** `[]interface{}` / `map[string]interface{}` after validation so generated function interpreters (which type-assert slices until a future codegen batch) keep working. Bind converts to immutable at boundaries.

---

### Task 1: Literal eval emits immutable collections

**Files:**

- Modify: `pkg/conflow/array.go:64-78`
- Modify: `pkg/conflow/map.go:66-80`
- Create: `pkg/conflow/literal_isolation_test.go`

**Step 1: Write the failing test**

Ginkgo tests in `literal_isolation_test.go`:

- Evaluate `ArrayNode` with two integer child nodes → result is `*values.List[interface{}]`, not `[]interface{}`
- Evaluate empty `ArrayNode` → empty frozen list (Len 0)
- Evaluate `MapNode` with one entry → `*values.Map[string, interface{}]`
- Evaluate empty `MapNode` → empty frozen map
- Nested literal: inner list from array literal is also `*values.List[interface{}]`
- Mutating `Elems()` copy does not change frozen list

**Step 2: Run test to verify it fails**

Run: `go test ./pkg/conflow/... -v -run Literal`

Expected: FAIL — result is `[]interface{}` or `map[string]interface{}`.

**Step 3: Implement**

```go
// array.go Value()
if len(a.items) == 0 {
    return values.NewListBuilder[interface{}]().Freeze(), nil
}
builder := values.NewListBuilder[interface{}]()
for _, item := range a.items {
    value, err := parsley.EvaluateNode(userCtx, item)
    if err != nil {
        return nil, err
    }
    builder.Append(value)
}
return builder.Freeze(), nil
```

```go
// map.go Value() — same pattern with MapBuilder
```

**Step 4: Run test to verify it passes**

Run: `go test ./pkg/conflow/... -v -run Literal`

Expected: PASS

---

### Task 2: Schema `ValidateValue` accepts immutable types

**Files:**

- Modify: `pkg/schema/array.go:256-357` (`ValidateValue`; optionally `CompareValues`, `StringValue`)
- Modify: `pkg/schema/map.go:285-358`
- Modify: `pkg/schema/array_test.go`, `pkg/schema/map_test.go`

**Step 1: Write the failing tests**

Add `DescribeTable` entries:

**Array:**

- `values.ListOf(int64(1))` with integer items schema → no error
- Empty `values.NewListBuilder[interface{}]().Freeze()` → no error
- Nested list built from immutable literals validates under nested array schema

**Map:**

- `values.MapOf(map[string]interface{}{"a": int64(1)})` → no error
- Empty frozen map → no error

Assert returned value is still `[]interface{}` / `map[string]interface{}` (normalized output).

**Step 2: Run test to verify it fails**

Run: `go test ./pkg/schema/... -v -run "Validate accepts immutable"`

Expected: FAIL — `"must be array"` / `"must be map"`.

**Step 3: Implement**

At top of `ValidateValue`, normalize input via `values.AsInterfaceSlice` / `values.AsStringInterfaceMap`:

```go
v, err := values.AsInterfaceSlice(value)
if err != nil {
    return nil, errors.New("must be array")
}
// existing validation on v; return v (slice), not immutable handle
```

Same for map. Update `CompareValues` / `StringValue` to use the same helpers so const/enum checks work if called directly with immutable values.

**Step 4: Run tests**

Run: `go test ./pkg/schema/... -v -run "Validate accepts"`

Expected: PASS

---

### Task 3: Parser test helper and `Len` function compatibility

**Files:**

- Modify: `pkg/test/helper.go` (`ExpectParserToEvaluate`, `ExpectFunctionToEvaluate`)
- Modify: `pkg/functions/len.go`

**Step 1: Normalize eval comparison**

Add `normalizeEvalValue(v interface{}) interface{}` that recursively converts `*values.List[interface{}]` → `[]interface{}` and `*values.Map[string, interface{}]` → `map[string]interface{}` for `reflect.DeepEqual` comparison. Use in `ExpectParserToEvaluate` so existing parser tests keep expecting slice/map literals without per-file edits.

**Step 2: Update `Len`**

Handle `*values.List[interface{}]` and `*values.Map[string, interface{}]` in the `Len` switch (return `.Len()`).

**Step 3: Run parser and function tests**

Run: `go test ./pkg/parsers/... ./pkg/functions/... -count=1`

Expected: PASS

---

### Task 4: Full regression

**Files:** (none)

Run: `go test ./... -count=1`

Expected: PASS

---

## Review

Before starting the tasks above, record the baseline SHA: `git rev-parse HEAD`

After completing all tasks:

1. Launch a review subagent:
   - **Quick:** use the `review-bugbot` skill on uncommitted changes, or
   - **Thorough:** use the Task tool with the `code-review` skill, passing this plan file path, baseline SHA, and design doc reference
2. Use the `receiving-code-review` skill to evaluate and address the feedback
3. Verify all fixes pass: `go test ./... -count=1`
4. Commits may happen only when the user explicitly requests them, and only after steps 1–3 are complete

## Batch 4 notes

- Codegen `SetParam` template should call `bind.BindValue` and use `values.AsInterfaceSlice` / `AsStringInterfaceMap` when assigning to block struct fields typed as `[]interface{}` / `map[string]interface{}`.
- Function codegen `AssignValue` still type-asserts `[]interface{}` — either keep `ValidateValue` returning slices (current Batch 3 policy) or extend function codegen in a follow-up batch.
- After Batch 4, `@input` could validate immutable values without slice round-trip; optional reorder to bind-before-validate if validation accepts bound form only.
- `static_container.go` `SetParam` path still lacks bind — consider with Batch 4 or follow-up.
