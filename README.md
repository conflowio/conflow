# Conflow - Configuration and workflow language

:exclamation: This project is a technology preview and still under development. If you would like to have a chat or contribute, please open an issue or drop an email to hello@[project name].io

## Introduction

Conflow is a strongly typed configuration language focusing on simplicity and usability.
It makes conscious choices, so you can be more productive and make less mistakes.

If you are tired using YAML, hopefully you’ll love Conflow.

Conflow is able to **generate, parse and evaluate your own domain specific language (DSL)**. Your DSL can be used purely for configuration or you can define **complex workflows** with custom business logic. It's **written in Go** and similar to the Go language's syntax where applicable.

It generates a **parallel programming language** where a piece of code will be evaluated and run when all its dependencies are available. This capability makes it especially suitable for creating **workflow-as-code** languages.

It generates a **strongly typed language**, and runs **static checking** to catch most type errors before evaluation.

As a language creator you only need to define simple Go structs or functions and Conflow will generate the necessary Go code for parsing and evaluation. **No reflection is used runtime.**

## Simple demonstration

Let's say you want to write a DSL where you can say hello to the world.

This example assumes you've built or downloaded the `conflow` binary and it's available on your PATH.

This example is kept simple for demonstrational purposes, but a fully working version can be found in [examples/helloworld](examples/helloworld).

At first you will write your Hello struct in `hello.go`:

```go
package hello

import (
	"context"
	"fmt"
	"io"

	"github.com/conflowio/conflow/pkg/conflow"
)

// @block "task"
type Hello struct {
	// @id
	id conflow.ID
	// @required
	// @minLength 1
	to string
	// @dependency
	stdout io.Writer
}

func (h *Hello) ID() conflow.ID {
	return h.id
}

func (h *Hello) Run(ctx context.Context) (conflow.Result, error) {
	if _, err := fmt.Fprintf(h.stdout, "Hello %s!\n", h.to); err != nil {
		return nil, err
	}
	return nil, nil
}
```

This will define a block type called `hello` with a required parameter called `to` which will print "Hello [to]" to the stdout.

After running `go generate` a new file called `hello.cf.go` will be generated in the same directory, containing a struct called `HelloInterpreter`.

Then you can write code like this:

```conflow
hello {
  to = "World"
}
```

See [examples/helloworld](examples/helloworld) for the rest of the code.

## Main concepts

### Blocks and parameters

 * The language consists of blocks and parameters
 * The code you write is the body of the root level block called _main_.
 * A block can have predefined input parameters, user defined parameters, output parameters and blocks
 * A parameter's value can be:
    * a literal value (string, integer, bool, array, map, etc.)
    * an element of an array or map
    * a function call
    * a parameter reference
    * a complex (arithmetic, boolean, etc.) expression of all the above
 * A block can have a globally unique identifier and parameters can only be referenced as `<block id>.<parameter name>`, e.g. `foo.bar`
 * A block's body is optional
 * Blocks are registered in a global context, so a block can reference any named block's parameters

```conflow

// This is the body of the top level block called "main"

program := "test.sh" // This is a user defined parameter in the main block

test exec { // This is an "exec" type block with the id "test"
    program = main.program // Parameter referencing the "program" parameter in the top level block
}

// If the block supports it you can use the short block format if one parameter is marked as a value parameter
print "Result was: " + test.stdout

```

### Block lifecycle and business logic

 * When a block has all its dependencies available a block instance will be created
 * Only one block instance can exist of the same block at a given time (generated dependencies can cause new block instances to be created)
 * A block instance has a multi-stage init-main-close lifecycle
 * All stages can have custom business logic and are defined as a method on the block's struct (called init, main and close)
 * Inside a block parameters and child blocks are lazily evaluated (only when the matching stage has been started)
 * The init call returns with a boolean parameter to signal whether the block should run or be skipped (conditional runs)

```go
// @block "task"
type SampleBlock struct {
    // @id
	id      conflow.ID
	// @eval_stage "init"
	skipped bool
}

// conflow.BlockInitialiser interface
func (s *SampleBlock) Init(ctx context.Context) (bool, error) {
	return s.skipped, nil
}

// conflow.BlockRunner interface
func (s *SampleBlock) Run(ctx context.Context) (conflow.Result, error) {
	return nil, nil
}

// conflow.BlockCloser interface
func (s *SampleBlock) Close(ctx context.Context) error {
	return nil
}
```

### Block generators

* A block can generate multiple output blocks of the same type.
* Any output blocks need to be defined and have an identifier
* This concept is useful for e.g. implementing a simple iterator, a ticker, a queue reader, etc.
* You can write blocks which emit these output blocks if a change happens (e.g. on a file change re-read some configuration)
* Generated blocks can be sent using the block's context and the call blocks until all dependent blocks registered the new dependency or have been run.
* An error will be returned if any dependent blocks had an error during run.

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

func (it *Iterator) Run(ctx context.Context) (conflow.Result, error) {
	for i := int64(0); i < it.count; i++ {
		_, err := it.blockPublisher.PublishBlock(&It{
			id:    it.it.id,
			value: i,
		}, nil)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}
```

```conflow
iterator {
    count = 3
    i1 it // This will be the output block (no body)
}

println { // A block instance will be created for every i1.value value
  value = i1.value
}
```

See [examples/iterator](examples/iterator), [examples/ticker](examples/ticker) or [examples/licensify](examples/licensify) for working examples.

### Dependencies and evaluation order

 * Blocks can depend on other blocks if any of their parameters depend on an other block
 * Parameters can depend on other block parameters or parameters from the same block
 * A parameter or child block will be evaluated if its the matching evaluation stage and all dependencies were evaluated previously

```conflow
baz block { // "baz" will be evaluated after "bar"
    p2 = bar.p1
}

bar block {
    p1 = bar.u1 // This parameter will be evaluated after u1
    u1 := "user defined"
}
```

### Code generation

* A block can be defined as a Go struct
* Custom functions can be defined as simple Go functions
* You have to add the `@block` and `@function` directives to the comment header of your block structs and functions and run "conflow generate" whenever you change your implementations
* Conflow will generate a file next to your implementation with a `.cf.go` extension containing an Interpreter struct.

See [examples](examples) for various block definitions.

Simple function example:
```go
// @block
func Lower(s string) string {
	return strings.ToLower(s)
}
```

See [function](function) for the built-in functions.

## Language specification

### Notation

The syntax is specified using [Extended Backus-Naur Form (EBNF)](https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_form). A quick reference can be found [here](https://golang.org/ref/spec#Notation).

[TODO]
