package main

import (
	"path/filepath"
	"runtime"

	"github.com/opsidian/basil/basil/job"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/block"
	"github.com/opsidian/basil/examples/common"
	"github.com/opsidian/basil/parser"
	"github.com/opsidian/basil/util"
)

//go:generate basil generate
type Main struct {
	id basil.ID `basil:"id"`
}

func (m *Main) ParseContextOverride() basil.ParseContextOverride {
	return basil.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"print": common.PrintInterpreter{},
		},
	}
}

func main() {
	ctx, cancel := util.CreateDefaultContext()
	defer cancel()

	parseCtx := basil.NewParseContext(basil.NewIDRegistry(8, 16))

	p := parser.NewMain("main", MainInterpreter{})

	paths, _ := filepath.Glob("*.basil")

	if err := p.ParseFiles(
		parseCtx,
		paths...,
	); err != nil {
		panic(err)
	}

	scheduler := job.NewScheduler(runtime.NumCPU()*2, 100)
	scheduler.Start()
	defer scheduler.Stop()

	if _, err := basil.Evaluate(
		parseCtx,
		basil.NewBlockContext(ctx, nil),
		scheduler,
		"main",
	); err != nil {
		panic(err)
	}
}
