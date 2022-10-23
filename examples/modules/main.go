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
	"github.com/conflowio/conflow/pkg/blocks"
	"github.com/conflowio/conflow/pkg/conflow/block"
	"github.com/conflowio/conflow/pkg/functions"
	"github.com/conflowio/conflow/pkg/parsers"
	"github.com/conflowio/conflow/pkg/util"
)

func main() {
	ctx, cancel := util.CreateDefaultContext()
	defer cancel()

	parseCtx := common.NewParseContext()

	mainInterpreter := blocks.MainInterpreter{}
	mainInterpreter.BlockTransformerRegistry = block.InterpreterRegistry{
		"module":   mainInterpreter,
		"import":   blocks.ImportInterpreter{},
		"iterator": common.IteratorInterpreter{},
		"print":    blocks.PrintInterpreter{},
		"println":  blocks.PrintlnInterpreter{},
	}
	mainInterpreter.FunctionTransformerRegistry = functions.DefaultRegistry()

	p := parsers.NewMain("main", mainInterpreter)

	if err := p.ParseFile(
		parseCtx,
		"main.cf",
	); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	common.Main(ctx, parseCtx, nil)
}
