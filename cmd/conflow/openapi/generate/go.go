// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generate

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/conflowio/conflow/pkg/openapi"
	"github.com/conflowio/conflow/pkg/openapi/generator"
	"github.com/conflowio/conflow/pkg/util"
)

func goCommand() *cobra.Command {
	var stdin bool
	var recursive bool
	var router string
	var packageName string
	var outputDir string

	cmd := &cobra.Command{
		Use:   "go",
		Short: "Generates Go implementation files",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if stdin && len(args) > 0 {
				return fmt.Errorf("no arguments are allowed when --stdin is used")
			}

			supportedRouters := []string{"echo"}
			if !util.StringSliceContains(supportedRouters, router) {
				return fmt.Errorf("unsupported router: %s, valid values: %v", router, supportedRouters)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var res *openapi.OpenAPI
			if stdin {
				var err error
				if res, err = evaluateStdin(cmd.Context()); err != nil {
					return err
				}
			} else {
				var err error
				if res, err = evaluatePath(cmd.Context(), args, recursive); err != nil {
					return err
				}
			}

			return generator.Generate(res, router, packageName, outputDir)
		},
	}

	cmd.Flags().BoolVar(&stdin, "stdin", false, "Read from stdin")
	cmd.Flags().BoolVarP(&recursive, "recursive", "r", true, "Read files recursively when a directory name is given")
	cmd.Flags().StringVar(&router, "router", "echo", "Go router library to use")
	cmd.Flags().StringVar(&packageName, "package", "", "Full package name")
	cmd.Flags().StringVar(&outputDir, "output-dir", "", "Path where the files will be generated")

	_ = cobra.MarkFlagRequired(cmd.Flags(), "package")
	_ = cobra.MarkFlagRequired(cmd.Flags(), "output-dir")

	return cmd
}
