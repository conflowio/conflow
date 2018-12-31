package fixtures

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/block"
)

//go:generate basil generate
type BlockWithBlock struct {
	IDField basil.ID       `basil:"id"`
	Blocks  []*BlockSimple `basil:"block=block_simple"`
}

func (b *BlockWithBlock) ID() basil.ID {
	return b.IDField
}

func (b *BlockWithBlock) ParseContext(ctx *basil.ParseContext) *basil.ParseContext {
	return ctx.New(basil.ParseContextConfig{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"block_simple": BlockSimpleInterpreter{},
		},
	})
}
