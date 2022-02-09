// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	"os"

	"github.com/conflowio/conflow/examples/common"
	"github.com/conflowio/conflow/src/blocks"
	"github.com/conflowio/conflow/src/conflow/block"
	"github.com/conflowio/conflow/src/functions"
	"github.com/conflowio/conflow/src/parsers"
	"github.com/conflowio/conflow/src/util"
)

func main() {
	ctx, cancel := util.CreateDefaultContext()
	defer cancel()

	parseCtx := common.NewParseContext()

	p := parsers.NewRoot("root", blocks.RootInterpreter{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"fail":    FailInterpreter{},
			"print":   blocks.PrintInterpreter{},
			"println": blocks.PrintlnInterpreter{},
		},
		FunctionTransformerRegistry: functions.DefaultRegistry(),
	})

	if err := p.ParseFile(
		parseCtx,
		"main.cf",
	); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	common.Main(ctx, parseCtx, nil)
}
