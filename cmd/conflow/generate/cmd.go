// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generate

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"

	"github.com/conflowio/conflow/pkg/conflow/generator"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generates Conflow files",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("couldn't determine current working directory: %w", err)
			}

			goPath := os.Getenv("GOPATH")
			if goPath == "" || !path.IsAbs(goPath) {
				return errors.New("GOPATH is not defined or invalid")
			}
			srcPath := path.Join(goPath, "src")

			if len(args) > 0 {
				if path.IsAbs(args[0]) {
					target = args[0]
				} else {
					target = path.Join(target, args[0])
				}
			}
			target = path.Clean(target)

			if !strings.HasPrefix(target, srcPath) {
				return fmt.Errorf("path must be in %s", srcPath)
			}

			return generator.Generate(target)
		},
	}
	return cmd
}
