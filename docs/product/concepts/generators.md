---
title: Generators
summary: Generator blocks, generated children, PublishBlock, and iterator/stream patterns.
parent: concepts/index.md
keywords: [generators, PublishBlock, iterator, dynamic blocks]
---

# Generators

**Generator** blocks emit **generated** child block instances at runtime. This enables iterators, tickers, subprocess stdout/stderr streams, and reactive re-evaluation patterns.

## Concepts

| Term | Meaning |
|------|---------|
| Generator block | `type = "generator"`; runs `Run()` and calls `PublishBlock` |
| Generated block type | Child type marked `@generated` on the generator struct |
| Generated instance | Concrete child published at runtime (e.g. each `it` with a new `value`) |

In Conflow source, you **declare** the generated child slot in the generator body:

```conflow
iterator {
    count = 3
    i1 it
}
```

Each published `it` instance can trigger dependents (e.g. `println { value = i1.value }`) once per publication.

## PublishBlock

`BlockPublisher` (`pkg/conflow/block.go`):

```go
PublishBlock(block Block, onScheduled func() error) (published bool, err error)
```

Behavior:

- Returns immediately with `published=false` if nothing depends on the published block.
- Otherwise blocks until dependent blocks complete processing that publication.
- `onScheduled` callback runs after the published block is scheduled.

## Iterator generator (built-in)

`pkg/blocks/iterator.go`:

```go
// @block "generator"
type Iterator struct {
    // @id
    id conflow.ID
    // @required
    count int64
    // @generated
    it *It
    // @dependency
    blockPublisher conflow.BlockPublisher
}
```

`Run` loops `count` times, publishing `It{ value: i }`.

## Exec stdout/stderr

`pkg/blocks/exec.go` publishes `Stream` blocks for stdout and stderr while the subprocess runs — dependents can consume streams in parallel with process execution.

## Ticker

`pkg/blocks/ticker.go` — time-based generator emitting `tick` blocks (`examples/ticker`).

## Dependency graph impact

Generators add **start** nodes in the dependency resolver so dependents on generated IDs wait for publication start without cyclic dependencies on the generator's own fields.

## Design use cases

| Use case | Pattern |
|----------|---------|
| For-each loop | `iterator` + generated `it` |
| Periodic work | `ticker` + `tick` |
| Stream processing | `exec` + `line_scanner` on stream |
| File watcher / reload | Custom generator republishing config blocks |

## Examples

- `examples/iterator` — nested iterators over arrays
- `examples/exec` — streaming command output
- `examples/ticker` — timed events
- `examples/licensify` — file walker generator

## See also

- [Block types](./block-types.md)
- [Built-in blocks: iterator, exec](../reference/built-in-blocks.md)
- [Dependencies and evaluation order](./dependencies-and-order.md)
