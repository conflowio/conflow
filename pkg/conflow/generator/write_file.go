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

	generatortemplate "github.com/conflowio/conflow/pkg/conflow/generator/template"
	"github.com/conflowio/conflow/pkg/internal/utils"
)

func WriteGeneratedFile(dir, name string, content []byte, headerParams generatortemplate.HeaderParams) error {
	header, err := generatortemplate.GenerateHeader(headerParams)
	if err != nil {
		return err
	}

	content = append(header, content...)

	conflowFile := utils.ToSnakeCase(name) + ".cf.go"
	targetPath := path.Join(dir, conflowFile)

	info, direrr := os.Stat(dir)
	if os.IsExist(direrr) && !info.IsDir() {
		return fmt.Errorf("%s exists, but not a directory", dir)
	}

	if !os.IsExist(direrr) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create %s directory: %w", dir, err)
		}
	}

	if err := os.WriteFile(targetPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write %s to %s: %w", name, getRelativePath(targetPath), err)
	}

	formatted, err := format.Source(content)
	if err != nil {
		return err
	}

	if err := os.WriteFile(targetPath, formatted, 0644); err != nil {
		return fmt.Errorf("failed to write %s to %s: %w", name, getRelativePath(targetPath), err)
	}

	fmt.Printf("Wrote `%s` to `%s`\n", name, getRelativePath(targetPath))

	return nil
}

func getRelativePath(path string) string {
	_, caller, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(caller)
	return strings.Replace(path, basePath+"/", "", 1)
}
