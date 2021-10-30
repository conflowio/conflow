// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package common

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"

	"github.com/opsidian/parsley/parsley"
	uzerolog "github.com/rs/zerolog"

	"github.com/opsidian/conflow/basil"
	"github.com/opsidian/conflow/basil/job"
	"github.com/opsidian/conflow/directives"
	"github.com/opsidian/conflow/loggers/zerolog"
)

func NewParseContext() *basil.ParseContext {
	idRegistry := basil.NewIDRegistry(8, 16)
	return basil.NewParseContext(parsley.NewFileSet(), idRegistry, directives.DefaultRegistry())
}

func Main(ctx context.Context, parseCtx *basil.ParseContext, inputParams map[basil.ID]interface{}) {
	level := uzerolog.InfoLevel
	if envLevel := os.Getenv("BASIL_LOG"); envLevel != "" {
		var err error
		level, err = uzerolog.ParseLevel(envLevel)
		if err != nil {
			fmt.Printf("Error: invalid log level %q\n", envLevel)
			os.Exit(1)
		}
	}

	logger := zerolog.NewConsoleLogger(level)
	scheduler := job.NewScheduler(logger, runtime.NumCPU()*2, 100)
	scheduler.Start()
	defer scheduler.Stop()

	go func() {
		err := http.ListenAndServe("localhost:6060", nil)
		logger.Debug().Err(err).Msg("http server stopped")
	}()

	if _, err := basil.Evaluate(
		parseCtx,
		ctx,
		nil,
		logger,
		scheduler,
		"main",
		inputParams,
	); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}
