package fixtures

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/block"
	"github.com/opsidian/basil/variable"
)

//go:generate basil generate
type BlockWithBlock struct {
	IDField    variable.ID       `basil:"id"`
	BlockNodes []basil.BlockNode `basil:"node"`
	Blocks     []*BlockSimple    `basil:"block"`
}

func (b *BlockWithBlock) ID() variable.ID {
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
