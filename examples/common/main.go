// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package common

import (
	"context"
	"os"
	"runtime"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/job"
	"github.com/opsidian/basil/logger"
	"github.com/rs/zerolog"
)

func Main(ctx context.Context, parseCtx *basil.ParseContext) {
	level := zerolog.InfoLevel
	if os.Getenv("DEBUG") == "1" {
		level = zerolog.DebugLevel
	}

	zl := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05.000"}).With().
		Timestamp().
		Logger().
		Level(level)

	scheduler := job.NewScheduler(runtime.NumCPU()*2, 100)
	scheduler.Start()
	defer scheduler.Stop()

	if _, err := basil.Evaluate(
		parseCtx,
		ctx,
		nil,
		logger.NewZeroLogLogger(zl),
		scheduler,
		"main",
	); err != nil {
		panic(err)
	}
}
