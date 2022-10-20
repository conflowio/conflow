// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package openapi

import (
	"github.com/spf13/cobra"

	"github.com/conflowio/conflow/cmd/conflow/openapi/generate"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "openapi",
		Short: "Work with OpenAPI definitions",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(generate.Command())

	return cmd
}
