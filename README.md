# Basil - DSL generator and parser

## Concepts

### All block ids are global.

This helps to move blocks around easily without changing references. Also shortens variable references, e.g. `build_project.version`.

### All variables must be referenced with the block id they are defined in.

This makes it clear where a variable is coming from. Especially helpful when a block is long as it's might be hard to see what level are you currently in. This means no `self.foo` or `this.foo` references.

Top level variables are defined in a module called `main` and should be referenced as `main.foo` 



