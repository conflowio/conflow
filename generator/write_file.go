package generator

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

func writeFile(name string, filename string, path string, content []byte) error {
	formatted, err := format.Source(content)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, formatted, 0644)
	if err != nil {
		return errors.Wrapf(err, "failed to write %s to %s", name, getRelativePath(path))
	}

	goimportsCmd := exec.Command("goimports", filename)
	out, err := goimportsCmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "failed to run goimports on %s", getRelativePath(path))
	}
	err = ioutil.WriteFile(path, out, 0644)
	if err != nil {
		return errors.Wrapf(err, "failed to write %s to %s", name, getRelativePath(path))
	}

	fmt.Printf("Wrote `%s` to `%s`\n", name, getRelativePath(path))

	return nil
}

func getRelativePath(path string) string {
	_, caller, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(caller)
	return strings.Replace(path, basePath+"/", "", 1)
}
