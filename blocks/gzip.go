// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blocks

import (
	"compress/gzip"
	"io"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/block"
)

// @block
type Gzip struct {
	// @id
	id basil.ID
	// @required
	in io.ReadCloser
	// @generated
	out *Stream
}

func (g *Gzip) ID() basil.ID {
	return g.id
}

func (g *Gzip) Main(ctx basil.BlockContext) error {
	var pipeWriter io.WriteCloser
	g.out.Stream, pipeWriter = io.Pipe()
	defer g.out.Stream.Close()
	defer pipeWriter.Close()

	published, err := ctx.PublishBlock(g.out, func() error {
		gzipWriter := gzip.NewWriter(pipeWriter)
		_, err := io.Copy(gzipWriter, g.in)
		_ = gzipWriter.Close()
		_ = pipeWriter.Close()
		return err
	})
	if !published {
		_, _ = io.Copy(io.Discard, g.in)
	}

	return err
}

func (g *Gzip) ParseContextOverride() basil.ParseContextOverride {
	return basil.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"out": StreamInterpreter{},
		},
	}
}
