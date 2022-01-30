// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blocks

import (
	"compress/gzip"
	"context"
	"io"

	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/conflow/block"
)

// @block "task"
type Gunzip struct {
	// @id
	id conflow.ID
	// @required
	in io.ReadCloser
	// @generated
	out *Stream
	// @dependency
	blockPublisher conflow.BlockPublisher
}

func (g *Gunzip) ID() conflow.ID {
	return g.id
}

func (g *Gunzip) Run(ctx context.Context) (conflow.Result, error) {
	var err error
	g.out.Stream, err = gzip.NewReader(g.in)
	if err != nil {
		return nil, err
	}
	defer g.out.Stream.Close()

	published, err := g.blockPublisher.PublishBlock(g.out, nil)
	if !published {
		_, _ = io.Copy(io.Discard, g.in)
	}

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (g *Gunzip) ParseContextOverride() conflow.ParseContextOverride {
	return conflow.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"out": StreamInterpreter{},
		},
	}
}
