package fixtures

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/block"
)

//go:generate basil generate BlockWithFactory
type BlockWithFactory struct {
	IDField    basil.ID          `basil:"id"`
	BlockNodes []basil.BlockNode `basil:"node,stage=-"`
}

func (b *BlockWithFactory) ID() basil.ID {
	return b.IDField
}

func (b *BlockWithFactory) Type() string {
	return "block_with_factory"
}

func (b *BlockWithFactory) Context(ctx interface{}) interface{} {
	return ctx
}

func (b *BlockWithFactory) Registry() block.Registry {
	return block.Registry{
		"block_simple": BlockSimpleInterpreter{},
	}
}
