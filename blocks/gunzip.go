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

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/block"
)

// @block
type Gunzip struct {
	// @id
	id basil.ID
	// @required
	in io.ReadCloser
	// @generated
	out *Stream
	// @dependency
	blockPublisher basil.BlockPublisher
}

func (g *Gunzip) ID() basil.ID {
	return g.id
}

func (g *Gunzip) Run(ctx context.Context) error {
	var err error
	g.out.Stream, err = gzip.NewReader(g.in)
	if err != nil {
		return err
	}
	defer g.out.Stream.Close()

	published, err := g.blockPublisher.PublishBlock(g.out, nil)
	if !published {
		_, _ = io.Copy(io.Discard, g.in)
	}

	if err != nil {
		return err
	}

	return nil
}

func (g *Gunzip) ParseContextOverride() basil.ParseContextOverride {
	return basil.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"out": StreamInterpreter{},
		},
	}
}
