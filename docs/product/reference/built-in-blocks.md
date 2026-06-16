---
title: Built-in blocks
summary: Standard blocks in pkg/blocks for workflows, I/O, timing, and modules.
parent: reference/index.md
keywords: [built-in blocks, exec, sleep, import]
---

# Built-in blocks

Package: `pkg/blocks/`

Register interpreters on your main block's `ParseContextOverride` as needed. Examples often import `blocks.PrintInterpreter`, `blocks.ExecInterpreter`, etc.

## Task blocks

| Block | Type | Purpose | Source |
|-------|------|---------|--------|
| `exec` | task | Run subprocess; publishes stdout/stderr streams | `exec.go` |
| `sleep` | task | Delay execution | `sleep.go` |
| `print` | task | Write to stdout (value param) | `print.go` |
| `println` | task | Write line to stdout | `println.go` |
| `fail` | task | Fail with error (testing/retry demos) | `fail.go` |
| `gzip` | task | Gzip compress | `gzip.go` |
| `gunzip` | task | Gzip decompress | `gunzip.go` |
| `line_scanner` | task | Scan stream into lines | `line_scanner.go` |
| `import` | task (parse stage) | Load Conflow module from path | `import.go` |

## Generator blocks

| Block | Generated | Purpose | Source |
|-------|-----------|---------|--------|
| `iterator` | `it` | Counting iterator | `iterator.go`, `it.cf.go` |
| `ticker` | `tick` | Time-based events | `ticker.go`, `tick.cf.go` |

`exec` also publishes `stdout` / `stderr` stream configuration blocks (`stream.go`).

## Configuration blocks

| Block | Purpose | Source |
|-------|---------|--------|
| `basic` | Minimal configuration holder | `basic.go` |
| `stream` | IO stream wrapper | `stream.go` |
| `line` | Line from scanner | `line.cf.go` |

## Module support

| Component | Role |
|-----------|------|
| `blocks.NewModuleInterpreter` | Wraps parsed module `main` as importable block | `module.go` |
| `import` block | Parse-time loader | `import.go` |

### Import usage

```conflow
sum import "./sum"

result sum {
  a = 1
  b = 2
}
```

Module directory must expose a `main` block (`examples/modules/`).

## Exec example

`examples/exec/main.cf`:

```conflow
hello exec {
    cmd = "sh"
    params = ["-c", "..."]
    hello_stdout stdout
    hello_stderr stderr
}

println { value = hello_stdout.stream }
```

## Iterator example

```conflow
iterator {
    count = 3
    i1 it
}

println { value = i1.value }
```

## See also

- [Generators](../concepts/generators.md)
- [Example catalog](../examples/catalog.md)
