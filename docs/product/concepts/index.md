---
title: Core concepts
summary: Index of fundamental Conflow product concepts — blocks, parameters, dependencies, lifecycle, generators.
parent: index.md
keywords: [concepts, blocks, parameters, dependencies]
---

# Core concepts

Conflow programs are **graphs of typed blocks** evaluated in **stages** according to **dependencies**. Extension authors implement blocks in Go; end users compose them in Conflow text.

## Topics

| Document | Summary |
|----------|---------|
| [Blocks and parameters](./blocks-and-parameters.md) | Structure of blocks, parameter kinds, references |
| [Block types](./block-types.md) | `main`, `task`, `generator`, `configuration`, `directive` |
| [Dependencies and evaluation order](./dependencies-and-order.md) | How the graph is built and scheduled |
| [Lifecycle and evaluation stages](./lifecycle-and-stages.md) | `init`, `main`, `close`, lazy evaluation |
| [Generators](./generators.md) | Dynamic child blocks and `PublishBlock` |

## Mental model

```text
main (root)
 ├── user parameter: program := "test.sh"
 ├── task block: test exec { ... }
 └── task block: print { value = test.stdout }
```

- The **body of `main`** is what you write in a `.cf` file.
- **Named blocks** have IDs (`test`, `foo`); parameters are referenced as `test.stdout` or `main.program`.
- **Only one instance** of a named block exists at a time (generators create new instances via published children).

## Key implementation paths

| Concept | Package / file |
|---------|----------------|
| Block interfaces | `pkg/conflow/block.go` |
| Parameter model | `pkg/conflow/parameter.go` |
| Dependency resolver | `pkg/conflow/dependency/resolver.go` |
| Eval stages | `pkg/conflow/eval_context.go` |

## See also

- [Language syntax](../language/index.md)
- [Runtime evaluation](../runtime/index.md)
- [Defining blocks in Go](../extending/go-blocks.md)
