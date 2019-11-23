// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blocks

import (
	"time"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/block"
)

//go:generate basil generate
type Ticker struct {
	id       basil.ID      `basil:"id"`
	interval time.Duration `basil:"required"`
	tick     *Tick         `basil:"generated"`
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
			_, err := ctx.PublishBlock(&Tick{
				id:   t.tick.id,
				time: tickerTime,
			}, nil)
			if err != nil {
				return err
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
