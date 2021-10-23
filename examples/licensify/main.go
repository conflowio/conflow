// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	"os"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/block"
	"github.com/opsidian/basil/blocks"
	"github.com/opsidian/basil/examples/common"
	"github.com/opsidian/basil/parsers"
	"github.com/opsidian/basil/util"
)

// @block
type Main struct {
	// @id
	id basil.ID
}

func (m *Main) ID() basil.ID {
	return m.id
}

func (m *Main) ParseContextOverride() basil.ParseContextOverride {
	return basil.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"glob":      GlobInterpreter{},
			"licensify": LicensifyInterpreter{},
			"print":     blocks.PrintInterpreter{},
			"println":   blocks.PrintlnInterpreter{},
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
		"main.basil",
	); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	common.Main(ctx, parseCtx, nil)
}
