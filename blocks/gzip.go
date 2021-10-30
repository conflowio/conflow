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

	"github.com/opsidian/conflow/conflow"
	"github.com/opsidian/conflow/conflow/block"
)

// @block
type Gzip struct {
	// @id
	id conflow.ID
	// @required
	in io.ReadCloser
	// @generated
	out *Stream
	// @dependency
	blockPublisher conflow.BlockPublisher
}

func (g *Gzip) ID() conflow.ID {
	return g.id
}

func (g *Gzip) Run(ctx context.Context) (conflow.Result, error) {
	var pipeWriter io.WriteCloser
	g.out.Stream, pipeWriter = io.Pipe()
	defer g.out.Stream.Close()
	defer pipeWriter.Close()

	published, err := g.blockPublisher.PublishBlock(g.out, func() error {
		gzipWriter := gzip.NewWriter(pipeWriter)
		_, err := io.Copy(gzipWriter, g.in)
		_ = gzipWriter.Close()
		_ = pipeWriter.Close()
		return err
	})
	if !published {
		_, _ = io.Copy(io.Discard, g.in)
	}

	return nil, err
}

func (g *Gzip) ParseContextOverride() conflow.ParseContextOverride {
	return conflow.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"out": StreamInterpreter{},
		},
	}
}
