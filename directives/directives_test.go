// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives_test

import (
	"context"
	"fmt"
	"runtime"

	"github.com/opsidian/basil/basil/block"
	"github.com/opsidian/basil/functions"
	"github.com/opsidian/basil/parsers"

	"github.com/opsidian/basil/blocks"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/job"
	"github.com/opsidian/basil/directives"
	"github.com/opsidian/basil/loggers/zerolog"
)

func eval(input string) {
	idRegistry := basil.NewIDRegistry(8, 16)
	parseCtx := basil.NewParseContext(idRegistry, directives.DefaultRegistry())
	logger := zerolog.NewDisabledLogger()
	scheduler := job.NewScheduler(logger, runtime.NumCPU()*2, 100)
	scheduler.Start()
	defer scheduler.Stop()

	mainInterpreter := blocks.MainInterpreter{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"println": blocks.PrintlnInterpreter{},
			"range":   blocks.RangeInterpreter{},
			"sleep":   blocks.SleepInterpreter{},
		},
		FunctionTransformerRegistry: functions.DefaultRegistry(),
	}

	p := parsers.NewMain("main", mainInterpreter)
	if err := p.ParseText(parseCtx, input); err != nil {
		fmt.Printf("Example errored: %s\n", err.Error())
		return
	}

	if _, err := basil.Evaluate(parseCtx, context.Background(), nil, logger, scheduler, "main", nil); err != nil {
		fmt.Printf("Example errored: %s\n", err.Error())
		return
	}
}
