// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/opsidian/basil/basil/generator"
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
	packageName := os.Getenv("GOPACKAGE")
	filename := os.Getenv("GOFILE")
	line, _ := strconv.Atoi(os.Getenv("GOLINE"))
	err := generator.Generate(dir, packageName, filename, line)
	if err != nil {
		fail(fmt.Sprintf("failed to run generate on %s: %s", filename, err.Error()))
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
