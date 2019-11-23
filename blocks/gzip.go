package blocks

import (
	"compress/gzip"
	"io"
	"io/ioutil"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/block"
)

//go:generate basil generate
type Gzip struct {
	id  basil.ID      `basil:"id"`
	in  io.ReadCloser `basil:"required"`
	out *Stream       `basil:"generated"`
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
		_, _ = io.Copy(ioutil.Discard, g.in)
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
