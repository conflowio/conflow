package fixtures

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/block"
)

type BlockInterface interface {
	basil.Block
	Foo() string
}

//go:generate basil generate
type BlockWithBlockInterface struct {
	IDField     basil.ID         `basil:"id"`
	BlockSimple []BlockInterface `basil:"block"`
}

func (b *BlockWithBlockInterface) ParseContextOverride(ctx *basil.ParseContext) basil.ParseContextOverride {
	return basil.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"block_simple": BlockSimpleInterpreter{},
		},
	}
}
