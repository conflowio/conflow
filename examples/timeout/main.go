// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	"os"

	"github.com/conflowio/conflow/blocks"
	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/block"
	"github.com/conflowio/conflow/examples/common"
	"github.com/conflowio/conflow/parsers"
	"github.com/conflowio/conflow/util"
)

// @block
type Main struct {
	// @id
	id conflow.ID
}

func (m *Main) ID() conflow.ID {
	return m.id
}

func (m *Main) ParseContextOverride() conflow.ParseContextOverride {
	return conflow.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"sleep":   blocks.SleepInterpreter{},
			"print":   blocks.PrintInterpreter{},
			"println": blocks.PrintlnInterpreter{},
		},
	}
}

func main() {
	ctx, cancel := util.CreateDefaultContext()
	defer cancel()

	parseCtx := common.NewParseContext()

	p := parsers.NewMain("main", MainInterpreter{})

	if err := p.ParseFile(
		parseCtx,
		"main.cf",
	); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	common.Main(ctx, parseCtx, nil)
}
