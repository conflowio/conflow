// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"

	"github.com/opsidian/conflow/basil/block"
	"github.com/opsidian/conflow/blocks"
	"github.com/opsidian/conflow/conflow"
	"github.com/opsidian/conflow/examples/common"
	"github.com/opsidian/conflow/functions"
	"github.com/opsidian/conflow/parsers"
	"github.com/opsidian/conflow/util"
)

func main() {
	ctx, cancel := util.CreateDefaultContext()
	defer cancel()

	inputParams, err := evalParams()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	parseCtx := common.NewParseContext()

	p := parsers.NewMain("main", blocks.MainInterpreter{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"print":   blocks.PrintInterpreter{},
			"println": blocks.PrintlnInterpreter{},
		},
		FunctionTransformerRegistry: functions.DefaultRegistry(),
	})

	if err := p.ParseFile(
		parseCtx,
		"main.basil",
	); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	common.Main(ctx, parseCtx, inputParams)
}

func evalParams() (map[conflow.ID]interface{}, error) {
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	params := flags.String("params", "", "List of key-value pairs, separated by comma")
	if err := flags.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	f := text.NewFile("params", []byte(*params))
	fs := parsley.NewFileSet(f)

	value, err := parsley.Evaluate(parsley.NewContext(fs, text.NewReader(f)), combinator.Sentence(parsers.KeyValuePairs()))
	if err != nil {
		return nil, err
	}

	return value.(map[conflow.ID]interface{}), nil
}
