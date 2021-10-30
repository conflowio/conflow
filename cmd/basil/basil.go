// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/opsidian/conflow/basil/generator"
)

var usage = `
USAGE
    basil generate <path>
`

func main() {
	log.SetFlags(0)
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		log.Fatalf("command must be set\n%s", usage)
		return
	}

	dir := cwd()
	cmd := args[0]

	switch cmd {
	case "generate":
		generate(dir, args[1:])
	default:
		log.Fatalf("unknown command\n%s", usage)
	}
}

func generate(dir string, args []string) {
	if len(args) > 0 {
		if path.IsAbs(args[0]) {
			dir = args[0]
		} else {
			dir = path.Join(dir, args[0])
		}
	}
	err := generator.Generate(dir)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error: %s", err.Error()))
	}
}

func cwd() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("couldn't determine current working directory")
		return ""
	}
	return dir
}
