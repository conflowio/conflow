package fixtures

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/block"
)

//go:generate basil generate
type BlockWithBlock struct {
	IDField basil.ID       `basil:"id"`
	Blocks  []*BlockSimple `basil:"block=block_simple"`
}

func (b *BlockWithBlock) ParseContextOverride() basil.ParseContextOverride {
	return basil.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"block_simple": BlockSimpleInterpreter{},
		},
	}
}
