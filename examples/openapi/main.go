// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/conflowio/conflow/examples/common"
	"github.com/conflowio/conflow/src/openapi"
	"github.com/conflowio/conflow/src/parsers"
	"github.com/conflowio/conflow/src/util"
)

func main() {
	ctx, cancel := util.CreateDefaultContext()
	defer cancel()

	parseCtx := common.NewParseContext()

	p := parsers.NewMain("main", openapi.OpenAPIInterpreter{})

	if err := p.ParseDir(
		parseCtx,
		".",
	); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	res := common.Main(ctx, parseCtx, nil)

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(res); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}
