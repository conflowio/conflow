// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"github.com/opsidian/basil/basil/block"
	"github.com/opsidian/basil/blocks"
	"github.com/opsidian/basil/examples/common"
	"github.com/opsidian/basil/functions"
	"github.com/opsidian/basil/parsers"
	"github.com/opsidian/basil/util"
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
		"main.basil",
	); err != nil {
		panic(err)
	}

	common.Main(ctx, parseCtx, nil)
}
