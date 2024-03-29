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

	"github.com/conflowio/parsley/parsley"
	uzerolog "github.com/rs/zerolog"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/conflow/job"
	"github.com/conflowio/conflow/pkg/directives"
	"github.com/conflowio/conflow/pkg/loggers/zerolog"
)

func NewParseContext() *conflow.ParseContext {
	idRegistry := conflow.NewIDRegistry(8, 16)
	return conflow.NewParseContext(parsley.NewFileSet(), idRegistry, directives.DefaultRegistry())
}

func Main(ctx context.Context, parseCtx *conflow.ParseContext, inputParams map[conflow.ID]interface{}) interface{} {
	level := uzerolog.InfoLevel
	if envLevel := os.Getenv("CONFLOW_LOG"); envLevel != "" {
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

	res, err := conflow.Evaluate(
		parseCtx,
		ctx,
		nil,
		logger,
		scheduler,
		"main",
		inputParams,
	)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	return res
}
