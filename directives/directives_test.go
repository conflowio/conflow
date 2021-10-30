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

	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"

	"github.com/opsidian/conflow/basil/block"
	"github.com/opsidian/conflow/functions"
	"github.com/opsidian/conflow/parsers"

	"github.com/opsidian/conflow/blocks"

	"github.com/opsidian/conflow/basil"
	"github.com/opsidian/conflow/basil/job"
	"github.com/opsidian/conflow/directives"
	"github.com/opsidian/conflow/loggers/zerolog"
)

func eval(input string) {
	f := text.NewFile("test", []byte(input))
	fs := parsley.NewFileSet(f)
	idRegistry := basil.NewIDRegistry(8, 16)
	parseCtx := basil.NewParseContext(fs, idRegistry, directives.DefaultRegistry())
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
