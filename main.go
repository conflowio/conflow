package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
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

	cmd := args[0]
	switch cmd {
	case "generate":
		generate(args[1:])
	default:
		fail("unknown command\n%s", usage)
	}
}

func generate(args []string) {
	if len(args) == 0 {
		fail("struct name is missing\n%s", usage)
	}

	name := args[0]
	out, err := block.GenerateFactory(block.FactoryTemplateParams{
		Package: os.Getenv("GOPACKAGE"),
		Type:    "foo",
		Name:    name,
		Stages:  []string{"pre"},
		Params: map[string][]block.Param{
			"pre": []block.Param{
				block.Param{
					Name:      "f1",
					FieldName: "field1",
					Type:      "string",
					Required:  true,
				},
			},
			"default": []block.Param{
				block.Param{
					Name:      "f2",
					FieldName: "field2",
					Type:      "int64",
				},
			},
		},
	})

	if err != nil {
		fail(fmt.Sprintf("failed to generate %s: %s", name, err.Error()))
	}

	filename := regexp.MustCompile("[A-Z][a-z0-9_]*").ReplaceAllStringFunc(name, func(str string) string {
		return "_" + strings.ToLower(str)
	})
	filename = strings.TrimLeft(filename, "_") + "_factory.go"
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

	gofmtCmd := exec.Command("gofmt", filename)
	out, err = gofmtCmd.CombinedOutput()
	if err != nil {
		fail("failed to run gofmt on %s: %s", filePath, err.Error())
	}
	err = ioutil.WriteFile(filePath, out, 0644)
	if err != nil {
		fail("failed to write %s: %s", filePath, err.Error())
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
    ocl generate <struct name>
`
