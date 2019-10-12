// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package common

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/opsidian/basil/logger/zerolog"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/job"
	uzerolog "github.com/rs/zerolog"
)

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
