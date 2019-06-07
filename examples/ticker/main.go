package main

import (
	"github.com/opsidian/basil/function"

	"github.com/opsidian/basil/examples/common"
	"github.com/opsidian/basil/parser"
	"github.com/opsidian/basil/util"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/block"
)

//go:generate basil generate
type Main struct {
	id basil.ID `basil:"id"`
}

func (m *Main) ParseContextOverride() basil.ParseContextOverride {
	return basil.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"ticker":  TickerInterpreter{},
			"print":   common.PrintInterpreter{},
			"println": common.PrintlnInterpreter{},
		},
		FunctionTransformerRegistry: function.Registry(),
	}
}

func main() {
	ctx, cancel := util.CreateDefaultContext()
	defer cancel()

	parseCtx := basil.NewParseContext(basil.NewIDRegistry(8, 16))

	p := parser.NewMain("main", MainInterpreter{})

	if err := p.ParseFile(
		parseCtx,
		"main.basil",
	); err != nil {
		panic(err)
	}

	common.Main(ctx, parseCtx)
}
