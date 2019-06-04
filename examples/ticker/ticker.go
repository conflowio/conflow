package main

import (
	"time"

	"github.com/opsidian/basil/basil/block"

	"github.com/opsidian/basil/basil"
)

//go:generate basil generate
type Ticker struct {
	id       basil.ID      `basil:"id"`
	interval time.Duration `basil:"required"`
	ticks    int64         `basil:"output"`
	tick     chan *Tick    `basil:"block,output"`
	count    int64
}

func (t *Ticker) Main(ctx basil.BlockContext) error {
	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()

	for {
		select {
		case time := <-ticker.C:
			// We do a non-blocking send here, the tick will sent again at the next interval
			select {
			case t.tick <- &Tick{time: time}:
			default:
			}
			t.ticks++
			if t.count > 0 && t.ticks >= t.count {
				return nil
			}
		case <-ctx.Context().Done():
			return nil
		}
	}
}

func (t *Ticker) ParseContextOverride() basil.ParseContextOverride {
	return basil.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"tick": TickInterpreter{},
		},
	}
}

//go:generate basil generate
type Tick struct {
	id   basil.ID  `basil:"id"`
	time time.Time `basil:"output"`
}
