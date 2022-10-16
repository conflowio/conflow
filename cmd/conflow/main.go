// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/conflowio/conflow/cmd/conflow/generate"
)

func rootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "conflow <command> <subcommand> [args]",
		Short:         "Conflow - Configuration and workflow language",
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(generate.Command())

	return cmd
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	defer func() {
		signal.Reset(os.Interrupt)
	}()

	go func() {
		<-interrupt
		cancel()
	}()

	var logger zerolog.Logger
	if term.IsTerminal(int(os.Stdin.Fd())) {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05.000"}).With().
			Timestamp().
			Logger()
	} else {
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	rootCmd := rootCommand()
	cobra.CheckErr(rootCmd.ExecuteContext(logger.WithContext(ctx)))
}
