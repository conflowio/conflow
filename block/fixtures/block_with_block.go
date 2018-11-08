package fixtures

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/block"
)

//go:generate basil generate BlockWithBlock
type BlockWithBlock struct {
	IDField    basil.ID          `basil:"id"`
	BlockNodes []basil.BlockNode `basil:"node"`
	Blocks     []*BlockSimple    `basil:"block"`
}

func (b *BlockWithBlock) ID() basil.ID {
	return b.IDField
}

func (b *BlockWithBlock) Type() string {
	return "block_with_block"
}

func (b *BlockWithBlock) Context(ctx interface{}) interface{} {
	return ctx
}

func (b *BlockWithBlock) BlockRegistry() block.Registry {
	return block.Registry{
		"block_simple": BlockSimpleInterpreter{},
	}
}
