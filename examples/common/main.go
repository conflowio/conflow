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

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/job"
	"github.com/opsidian/basil/directives"
	"github.com/opsidian/basil/logger/zerolog"
	uzerolog "github.com/rs/zerolog"
)

func NewParseContext() *basil.ParseContext {
	idRegistry := basil.NewIDRegistry(8, 16)
	return basil.NewParseContext(idRegistry, directives.DefaultRegistry())
}

func Main(ctx context.Context, parseCtx *basil.ParseContext) {
	level := uzerolog.InfoLevel
	if envLevel := os.Getenv("BASIL_LOG"); envLevel != "" {
		var err error
		level, err = uzerolog.ParseLevel(envLevel)
		if err != nil {
			panic(fmt.Errorf("invalid log level %q", envLevel))
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
	); err != nil {
		panic(err)
	}
}
