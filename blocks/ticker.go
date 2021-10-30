// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blocks

import (
	"context"
	"time"

	"github.com/opsidian/conflow/conflow"
	"github.com/opsidian/conflow/conflow/block"
)

// @block
type Ticker struct {
	// @id
	id conflow.ID
	// @required
	interval time.Duration
	// @generated
	tick *Tick
	// @dependency
	blockPublisher conflow.BlockPublisher
}

func (t *Ticker) ID() conflow.ID {
	return t.id
}

func (t *Ticker) Run(ctx context.Context) (conflow.Result, error) {
	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()

	for {
		select {
		case tickerTime := <-ticker.C:
			_, err := t.blockPublisher.PublishBlock(&Tick{
				id:   t.tick.id,
				time: tickerTime,
			}, nil)
			if err != nil {
				return nil, err
			}
		case <-ctx.Done():
			return nil, nil
		}
	}
}

func (t *Ticker) ParseContextOverride() conflow.ParseContextOverride {
	return conflow.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"tick": TickInterpreter{},
		},
	}
}

// @block
type Tick struct {
	// @id
	id conflow.ID
	// @read_only
	time time.Time
}

func (t *Tick) ID() conflow.ID {
	return t.id
}
