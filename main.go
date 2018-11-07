package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/opsidian/basil/block"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fail("command must be set\n%s", usage)
		return
	}

	dir := cwd()
	cmd := args[0]

	switch cmd {
	case "generate":
		generate(dir, args[1:])
	default:
		fail("unknown command\n%s", usage)
	}
}

func generate(dir string, args []string) {
	if len(args) == 0 {
		fail("struct name is missing\n%s", usage)
	}

	name := args[0]
	err := block.GenerateInterpreter(dir, name, os.Getenv("GOPACKAGE"))
	if err != nil {
		fail(fmt.Sprintf("failed to generate %s: %s", name, err.Error()))
	}
}

func cwd() string {
	dir, err := os.Getwd()
	if err != nil {
		fail("couldn't determine current working directory")
	}
	return dir
}

func fail(s string, args ...interface{}) {
	fmt.Printf("Error: "+s+"\n", args...)
	os.Exit(1)
}

var usage = `
USAGE
    basil generate <struct name>
`
