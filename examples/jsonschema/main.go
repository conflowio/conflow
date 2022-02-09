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
	"github.com/conflowio/conflow/src/schema"
	"github.com/conflowio/conflow/src/util"
)

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

	parseCtx := common.NewParseContext()

	p := parsers.NewRoot("root", blocks.RootInterpreter{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"print":   blocks.PrintInterpreter{},
			"println": blocks.PrintlnInterpreter{},
			"person":  personInterpreter,
		},
		FunctionTransformerRegistry: functions.DefaultRegistry(),
	})

	if err := p.ParseFile(
		parseCtx,
		"main.cf",
	); err != nil {
		return err
	}

	common.Main(ctx, parseCtx, nil)

	return nil
}
