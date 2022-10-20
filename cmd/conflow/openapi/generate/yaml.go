// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generate

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/conflowio/conflow/src/openapi"
)

func yamlCommand() *cobra.Command {
	var stdin bool
	var recursive bool

	cmd := &cobra.Command{
		Use:     "yaml",
		Short:   "Generates OpenAPI definition in YAML format",
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

			jsonData, err := json.Marshal(res)
			if err != nil {
				return err
			}

			m := map[string]interface{}{}
			if err := json.Unmarshal(jsonData, &m); err != nil {
				return err
			}

			encoder := yaml.NewEncoder(cmd.OutOrStdout())
			defer func() {
				if err := encoder.Close(); err != nil {
					_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Error: %s", err.Error())
					os.Exit(1)
				}
			}()
			return encoder.Encode(m)
		},
	}

	cmd.Flags().BoolVar(&stdin, "stdin", false, "Read from stdin")
	cmd.Flags().BoolVarP(&recursive, "recursive", "r", true, "Read files recursively when a directory name is given")

	return cmd
}
