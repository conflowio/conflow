package main

import (
	"time"

	"github.com/opsidian/basil/basil/block"

	"github.com/opsidian/basil/basil"
)

//go:generate basil generate
type Benchmark struct {
	id       basil.ID      `basil:"id"`
	duration time.Duration `basil:"required"`
	elapsed  time.Duration
	counter  int64                   `basil:"output"`
	run      chan basil.BlockMessage `basil:"block,output"`
}

func (b *Benchmark) Main(ctx basil.BlockContext) error {
	timer := time.NewTimer(b.duration)
	defer timer.Stop()

	started := time.Now()

	for {
		select {
		case <-timer.C:
			b.elapsed = time.Now().Sub(started)
			return nil
		default:
			b.counter++
			msg := basil.NewBlockMessage(&BenchmarkRun{cnt: b.counter})
			b.run <- msg
			<-msg.Done()
		}
	}
}

func (b *Benchmark) ParseContextOverride() basil.ParseContextOverride {
	return basil.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"run": BenchmarkRunInterpreter{},
		},
	}
}

//go:generate basil generate
type BenchmarkRun struct {
	id  basil.ID `basil:"id"`
	cnt int64    `basil:"output"`
}
