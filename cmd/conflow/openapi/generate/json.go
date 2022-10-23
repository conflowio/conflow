// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generate

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/conflowio/conflow/pkg/openapi"
)

func jsonCommand() *cobra.Command {
	var compact bool
	var stdin bool
	var recursive bool

	cmd := &cobra.Command{
		Use:     "json",
		Short:   "Generates OpenAPI definition in JSON format",
		Example: "  To read from stdin, use '-' as argument",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if stdin && len(args) > 0 {
				return fmt.Errorf("no arguments are allowed when --stdin is used")
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

			encoder := json.NewEncoder(cmd.OutOrStdout())
			if !compact {
				encoder.SetIndent("", "  ")
			}
			return encoder.Encode(res)
		},
	}

	cmd.Flags().BoolVar(&compact, "compact", false, "Compact output")
	cmd.Flags().BoolVar(&stdin, "stdin", false, "Read from stdin")
	cmd.Flags().BoolVarP(&recursive, "recursive", "r", true, "Read files recursively when a directory name is given")

	return cmd
}
