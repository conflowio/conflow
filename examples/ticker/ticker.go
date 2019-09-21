// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

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
	count    int64
	ticks    int64 `basil:"output"`
	tick     *Tick `basil:"generated"`
}

func (t *Ticker) ID() basil.ID {
	return t.id
}

func (t *Ticker) Main(ctx basil.BlockContext) error {
	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()

	for {
		select {
		case tickerTime := <-ticker.C:
			_ = ctx.PublishBlock(&Tick{
				id:   t.tick.id,
				time: tickerTime,
			})

			t.ticks++
			if t.count > 0 && t.ticks >= t.count {
				return nil
			}
		case <-ctx.Done():
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

func (t *Tick) ID() basil.ID {
	return t.id
}
