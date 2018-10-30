package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/opsidian/ocl/block"
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
	out, err := block.GenerateFactory(dir, name, os.Getenv("GOPACKAGE"))

	if err != nil {
		fail(fmt.Sprintf("failed to generate %s: %s", name, err.Error()))
	}

	filename := regexp.MustCompile("[A-Z][a-z0-9_]*").ReplaceAllStringFunc(name, func(str string) string {
		return "_" + strings.ToLower(str)
	})
	filename = strings.TrimLeft(filename, "_") + "_factory.ocl.go"
	filePath := path.Join(cwd(), filename)

	err = ioutil.WriteFile(filePath, out, 0644)
	if err != nil {
		fail("failed to write %s: %s", filePath, err.Error())
	}

	goimportsCmd := exec.Command("goimports", filename)
	out, err = goimportsCmd.CombinedOutput()
	if err != nil {
		fail("failed to run goimports on %s: %s", filePath, err.Error())
	}
	err = ioutil.WriteFile(filePath, out, 0644)
	if err != nil {
		fail("failed to write %s: %s", filePath, err.Error())
	}

	fmt.Printf("Wrote `%sFactory` to `%s`\n", name, getRelativePath(filePath))
}

func getRelativePath(path string) string {
	_, caller, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(caller)
	return strings.Replace(path, basePath+"/", "", 1)
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
    ocl generate <struct name>
`
