// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"github.com/opsidian/basil/blocks"
	"github.com/opsidian/basil/functions"
	"github.com/opsidian/basil/parsers"

	"github.com/opsidian/basil/examples/common"
	"github.com/opsidian/basil/util"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/block"
)

//go:generate basil generate
type Main struct {
	id basil.ID `basil:"id"`
}

func (m *Main) ID() basil.ID {
	return m.id
}

func (m *Main) ParseContextOverride() basil.ParseContextOverride {
	return basil.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"iterator": common.IteratorInterpreter{},
			"print":    blocks.PrintInterpreter{},
			"println":  blocks.PrintlnInterpreter{},
			"do":       MainInterpreter{},
		},
		FunctionTransformerRegistry: functions.DefaultRegistry(),
	}
}

func main() {
	ctx, cancel := util.CreateDefaultContext()
	defer cancel()

	parseCtx := common.NewParseContext()

	p := parsers.NewMain("main", MainInterpreter{})

	if err := p.ParseFile(
		parseCtx,
		"main.basil",
	); err != nil {
		panic(err)
	}

	common.Main(ctx, parseCtx)
}
