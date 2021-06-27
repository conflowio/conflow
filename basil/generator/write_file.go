// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generator

import (
	"fmt"
	"go/format"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/xerrors"

	"github.com/opsidian/basil/basil/generator/parser"
)

func writeFile(dir, name string, content []byte) error {
	basilFile := parser.ToSnakeCase(name) + ".basil.go"
	filepath := path.Join(dir, basilFile)

	info, direrr := os.Stat(dir)
	if os.IsExist(direrr) && !info.IsDir() {
		return fmt.Errorf("%s exists, but not a directory", dir)
	}

	if !os.IsExist(direrr) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create %s directory: %w", dir, err)
		}
	}

	err := os.WriteFile(filepath, content, 0644)
	if err != nil {
		return xerrors.Errorf("failed to write %s to %s: %w", name, getRelativePath(filepath), err)
	}

	formatted, err := format.Source(content)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath, formatted, 0644)
	if err != nil {
		return xerrors.Errorf("failed to write %s to %s: %w", name, getRelativePath(filepath), err)
	}

	fmt.Printf("Wrote `%s` to `%s`\n", name, getRelativePath(filepath))

	return nil
}

func getRelativePath(path string) string {
	_, caller, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(caller)
	return strings.Replace(path, basePath+"/", "", 1)
}
