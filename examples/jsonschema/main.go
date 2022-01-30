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
	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/conflow/block"
	"github.com/conflowio/conflow/src/functions"
	"github.com/conflowio/conflow/src/parsers"
	"github.com/conflowio/conflow/src/schema"
	"github.com/conflowio/conflow/src/util"
)

var mainRegistry = block.InterpreterRegistry{
	"print":   blocks.PrintInterpreter{},
	"println": blocks.PrintlnInterpreter{},
}

// @block "main"
type Main struct {
	// @id
	id conflow.ID
}

func (m *Main) ID() conflow.ID {
	return m.id
}

func (m *Main) ParseContextOverride() conflow.ParseContextOverride {
	return conflow.ParseContextOverride{
		BlockTransformerRegistry:    mainRegistry,
		FunctionTransformerRegistry: functions.DefaultRegistry(),
	}
}

func main() {
	if err := mainWithError(); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}

func mainWithError() error {
	ctx, cancel := util.CreateDefaultContext()
	defer cancel()

	schemaID := "file://person.json"
	schema.RegisterResolver(schemaID, schema.ResolverFunc(schema.ResolveFromFile))

	personInterpreter, err := block.NewInterpreter(schemaID)
	if err != nil {
		return err
	}
	mainRegistry["person"] = personInterpreter

	parseCtx := common.NewParseContext()

	p := parsers.NewMain("main", MainInterpreter{})

	if err := p.ParseFile(
		parseCtx,
		"main.cf",
	); err != nil {
		return err
	}

	common.Main(ctx, parseCtx, nil)

	return nil
}
