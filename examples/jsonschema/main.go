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
	"github.com/conflowio/conflow/conflow/schema"
	"github.com/conflowio/conflow/examples/common"
	"github.com/conflowio/conflow/functions"
	"github.com/conflowio/conflow/parsers"
	"github.com/conflowio/conflow/util"
)

var mainRegistry = block.InterpreterRegistry{
	"print":   blocks.PrintInterpreter{},
	"println": blocks.PrintlnInterpreter{},
}

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

	personSchemaJSON, err := os.ReadFile("person.json")
	if err != nil {
		return fmt.Errorf("failed to read person.json: %w", err)
	}
	personSchema, err := schema.UnmarshalJSON(personSchemaJSON)
	if err != nil {
		return fmt.Errorf("failed to parse person.json: %w", err)
	}

	if _, ok := personSchema.(schema.ObjectKind); !ok {
		return fmt.Errorf("person.json must define an object type, got %s", personSchema.Type())
	}

	mainRegistry["person"] = block.NewInterpreter(personSchema.(schema.ObjectKind))

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
