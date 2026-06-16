// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generate

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"

	"github.com/conflowio/conflow/pkg/conflow/generator"
)

func Command() *cobra.Command {
	var localPrefixes []string

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generates Conflow files",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("couldn't determine current working directory: %w", err)
			}

			if len(args) > 0 {
				if path.IsAbs(args[0]) {
					target = args[0]
				} else {
					target = path.Join(target, args[0])
				}
			}
			target = path.Clean(target)

			return generator.Generate(target, localPrefixes)
		},
	}

	cmd.Flags().StringArrayVar(&localPrefixes, "local", nil, "put imports beginning with this string after 3rd-party packages")

	return cmd
}
