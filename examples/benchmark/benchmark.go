// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"time"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/block"
)

// @block
type Benchmark struct {
	// @id
	id basil.ID
	// @required
	duration time.Duration
	elapsed  time.Duration
	// @read_only
	counter int64
	// @generated
	run *BenchmarkRun
	// @dependency
	blockPublisher basil.BlockPublisher
}

func (b *Benchmark) ID() basil.ID {
	return b.id
}

func (b *Benchmark) Run(ctx context.Context) (basil.Result, error) {
	timer := time.NewTimer(b.duration)
	defer timer.Stop()

	started := time.Now()

	for {
		select {
		case <-timer.C:
			b.elapsed = time.Now().Sub(started)
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

func (b *Benchmark) ParseContextOverride() basil.ParseContextOverride {
	return basil.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"run": BenchmarkRunInterpreter{},
		},
	}
}

// @block
type BenchmarkRun struct {
	// @id
	id basil.ID
	// @read_only
	cnt int64
}

func (b *BenchmarkRun) ID() basil.ID {
	return b.id
}
