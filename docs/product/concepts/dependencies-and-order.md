---
title: Dependencies and evaluation order
summary: How Conflow builds the dependency graph, resolves order, detects cycles, and schedules parallel work.
parent: concepts/index.md
keywords: [dependencies, evaluation order, parallel, cycles]
---

# Dependencies and evaluation order

## Dependency rules

1. **Block dependencies** — If any parameter of block B references a parameter or output of block A, B depends on A.
2. **Parameter dependencies** — Parameters within a block depend on other parameters they reference.
3. **Stage ordering** — Within a block, earlier evaluation stages complete before later ones for dependent parameters.

Example:

```conflow
baz block {
    p2 = bar.p1
}

bar block {
    p1 = bar.u1
    u1 := "user defined"
}
```

`baz` runs after `bar`; `bar.p1` runs after `bar.u1`.

## Parallel execution

Conflow generates a **parallel programming model**: independent blocks can run concurrently when the job scheduler has capacity. Blocks are **not** evaluated in strict source order unless dependencies force it.

When a block has all dependencies available for its current stage, a **block instance** is created and scheduled.

## Single instance rule

Only **one block instance** exists per named block at a time. **Generators** create new work by publishing **generated** child blocks (new instances of `it`, stream readers, etc.), which can trigger new dependent task instances.

## Dependency resolution

`pkg/conflow/dependency/resolver.go`:

- Builds a graph from all nodes in the workflow.
- Uses **Tarjan's strongly connected components** to detect cycles.
- Handles **generator** nodes specially: splits start vs finish nodes so dependents on generated children do not create false cycles with the generator's own fields.

Errors:

- Missing dependency (referenced block/parameter does not exist)
- Circular dependency

## Job scheduler

`pkg/conflow/job/`:

- `JobScheduler.ScheduleJob` queues block work.
- Examples use `runtime.NumCPU()*2` workers (`examples/common/main.go`).
- Lightweight jobs may run synchronously (`Job.Lightweight()`).

## Triggers (conditional scheduling)

`@triggers` restricts which upstream block completions cause a block to run. See [Runtime directives](../language/runtime-directives.md) and `examples/triggers`.

## PublishBlock synchronization

When a generator publishes a block, `PublishBlock` may block until dependent blocks finish consuming that publication — ensuring backpressure for iterators and streams.

## See also

- [Lifecycle and evaluation stages](./lifecycle-and-stages.md)
- [Generators](./generators.md)
- [Job scheduler](../runtime/job-scheduler.md)
