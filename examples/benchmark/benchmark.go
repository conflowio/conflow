// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"time"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/conflow/block"
	"github.com/conflowio/conflow/pkg/conflow/types"
)

// @block "task"
type Benchmark struct {
	// @id
	id conflow.ID
	// @required
	duration types.Duration
	elapsed  types.Duration
	// @read_only
	counter int64
	// @generated
	run *BenchmarkRun
	// @dependency
	blockPublisher conflow.BlockPublisher
}

func (b *Benchmark) ID() conflow.ID {
	return b.id
}

func (b *Benchmark) Run(ctx context.Context) (conflow.Result, error) {
	timer := time.NewTimer(time.Duration(b.duration))
	defer timer.Stop()

	started := time.Now()

	for {
		select {
		case <-timer.C:
			b.elapsed = types.Duration(time.Now().Sub(started))
			return nil, nil
		default:
			b.counter++
			_, _ = b.blockPublisher.PublishBlock(&BenchmarkRun{
				id:  b.run.id,
				cnt: b.counter,
			}, nil)
		}
	}
}

func (b *Benchmark) ParseContextOverride() conflow.ParseContextOverride {
	return conflow.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"run": BenchmarkRunInterpreter{},
		},
	}
}

// @block "configuration"
type BenchmarkRun struct {
	// @id
	id conflow.ID
	// @read_only
	cnt int64
}

func (b *BenchmarkRun) ID() conflow.ID {
	return b.id
}
