// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package generate

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/conflowio/conflow/examples/common"
	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/conflow/job"
	"github.com/conflowio/conflow/src/loggers/zerolog"
	"github.com/conflowio/conflow/src/openapi"
	"github.com/conflowio/conflow/src/parsers"
)

func evaluateStdin(ctx context.Context) (*openapi.OpenAPI, error) {
	parseCtx := common.NewParseContext()
	p := parsers.NewMain("main", openapi.OpenAPIInterpreter{})

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
			return nil, readErr
		}

		if err := p.ParseText(parseCtx, string(in)); err != nil {
			return nil, err
		}
	case <-ctx.Done():
		return nil, errors.New("interrupt received")
	}

	return evaluate(ctx, parseCtx)
}

func evaluatePath(ctx context.Context, args []string, recursive bool) (*openapi.OpenAPI, error) {
	parseCtx := common.NewParseContext()
	p := parsers.NewMain("main", openapi.OpenAPIInterpreter{})

	target, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("couldn't determine current working directory: %w", err)
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
		return nil, err
	}
	if fileInfo.IsDir() {
		if err := p.ParseDir(parseCtx, target, recursive); err != nil {
			return nil, err
		}
	} else {
		if err := p.ParseFile(parseCtx, target); err != nil {
			return nil, err
		}
	}

	return evaluate(ctx, parseCtx)
}

func evaluate(ctx context.Context, parseCtx *conflow.ParseContext) (*openapi.OpenAPI, error) {
	logger := zerolog.NewDisabledLogger()
	scheduler := job.NewScheduler(logger, runtime.NumCPU()*2, 100)
	scheduler.Start()
	defer scheduler.Stop()

	res, err := conflow.Evaluate(
		parseCtx,
		ctx,
		nil,
		logger,
		scheduler,
		"main",
		nil,
	)
	if err != nil {
		return nil, err
	}

	return res.(*openapi.OpenAPI), nil
}
