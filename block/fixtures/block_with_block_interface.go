package fixtures

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/block"
)

type BlockInterface interface {
	basil.Block
	Foo() string
}

//go:generate basil generate
type BlockWithBlockInterface struct {
	IDField basil.ID         `basil:"id"`
	Blocks  []BlockInterface `basil:"block=block_simple"`
}

func (b *BlockWithBlockInterface) ID() basil.ID {
	return b.IDField
}

func (b *BlockWithBlockInterface) ParseContext(ctx *basil.ParseContext) *basil.ParseContext {
	return ctx.New(basil.ParseContextConfig{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"block_simple": BlockSimpleInterpreter{},
		},
	})
}
