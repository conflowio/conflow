// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package openapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/conflowio/conflow/examples/common"
	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/conflow/job"
	"github.com/conflowio/conflow/src/loggers/zerolog"
	"github.com/conflowio/conflow/src/openapi"
	"github.com/conflowio/conflow/src/parsers"
)

func generateCommand() *cobra.Command {
	var format string
	var compact bool
	var stdin bool

	cmd := &cobra.Command{
		Use:     "generate",
		Short:   "Generates OpenAPI definition files",
		Example: "  To read from stdin, use '-' as argument",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if format != "json" && format != "yaml" {
				return fmt.Errorf("invalid value for 'format': must be json or yaml")
			}

			if stdin && len(args) > 0 {
				return fmt.Errorf("no arguments are allowed when --stdin is used")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			parseCtx := common.NewParseContext()
			p := parsers.NewMain("main", openapi.OpenAPIInterpreter{})

			if stdin {
				var in []byte
				var readErr error
				finished := make(chan struct{})
				go func() {
					in, readErr = io.ReadAll(os.Stdin)
					finished <- struct{}{}
				}()

				select {
				case <-finished:
					if readErr != nil {
						return readErr
					}

					if err := p.ParseText(parseCtx, string(in)); err != nil {
						return err
					}
				case <-cmd.Context().Done():
					return errors.New("interrupt received")
				}
			} else {
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

				fileInfo, err := os.Lstat(target)
				if err != nil {
					return err
				}
				if fileInfo.IsDir() {
					if err := p.ParseDir(parseCtx, target); err != nil {
						return err
					}
				} else {
					if err := p.ParseFile(parseCtx, target); err != nil {
						return err
					}
				}
			}

			logger := zerolog.NewDisabledLogger()
			scheduler := job.NewScheduler(logger, runtime.NumCPU()*2, 100)
			scheduler.Start()
			defer scheduler.Stop()

			res, err := conflow.Evaluate(
				parseCtx,
				cmd.Context(),
				nil,
				logger,
				scheduler,
				"main",
				nil,
			)
			if err != nil {
				return err
			}

			switch format {
			case "json":
				encoder := json.NewEncoder(cmd.OutOrStdout())
				if !compact {
					encoder.SetIndent("", "  ")
				}
				return encoder.Encode(res)
			case "yaml":
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
			default:
				panic(fmt.Errorf("unexpected format: %s", format))
			}
		},
	}

	cmd.Flags().StringVar(&format, "format", "json", "Output format: json, or yaml")
	cmd.Flags().BoolVar(&compact, "compact", false, "Whether to compact output (applies to json only)")
	cmd.Flags().BoolVar(&stdin, "stdin", false, "Whether to read from stdin")

	return cmd
}
