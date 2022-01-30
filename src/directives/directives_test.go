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

	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/text"

	"github.com/conflowio/conflow/src/blocks"
	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/conflow/block"
	"github.com/conflowio/conflow/src/conflow/job"
	"github.com/conflowio/conflow/src/directives"
	"github.com/conflowio/conflow/src/functions"
	"github.com/conflowio/conflow/src/loggers/zerolog"
	"github.com/conflowio/conflow/src/parsers"
)

func eval(input string) {
	f := text.NewFile("test", []byte(input))
	fs := parsley.NewFileSet(f)
	idRegistry := conflow.NewIDRegistry(8, 16)
	parseCtx := conflow.NewParseContext(fs, idRegistry, directives.DefaultRegistry())
	logger := zerolog.NewDisabledLogger()
	scheduler := job.NewScheduler(logger, runtime.NumCPU()*2, 100)
	scheduler.Start()
	defer scheduler.Stop()

	mainInterpreter := blocks.MainInterpreter{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"println":  blocks.PrintlnInterpreter{},
			"iterator": blocks.IteratorInterpreter{},
			"sleep":    blocks.SleepInterpreter{},
		},
		FunctionTransformerRegistry: functions.DefaultRegistry(),
	}

	p := parsers.NewMain("main", mainInterpreter)
	if err := p.ParseText(parseCtx, input); err != nil {
		fmt.Printf("Example errored: %s\n", err.Error())
		return
	}

	if _, err := conflow.Evaluate(parseCtx, context.Background(), nil, logger, scheduler, "main", nil); err != nil {
		fmt.Printf("Example errored: %s\n", err.Error())
		return
	}
}
