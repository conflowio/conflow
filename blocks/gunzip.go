package blocks

import (
	"compress/gzip"
	"io"
	"io/ioutil"

	"github.com/opsidian/basil/basil/block"

	"github.com/opsidian/basil/basil"
)

//go:generate basil generate
type Gunzip struct {
	id  basil.ID      `basil:"id"`
	in  io.ReadCloser `basil:"required"`
	out *Stream       `basil:"generated"`
}

func (g *Gunzip) ID() basil.ID {
	return g.id
}

func (g *Gunzip) Main(blockCtx basil.BlockContext) error {
	var err error
	g.out.Stream, err = gzip.NewReader(g.in)
	if err != nil {
		return err
	}
	defer g.out.Stream.Close()

	published, err := blockCtx.PublishBlock(g.out, nil)
	if !published {
		_, _ = io.Copy(ioutil.Discard, g.in)
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
