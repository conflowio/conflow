package fixtures

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/block"
)

type BlockInterface interface {
	basil.Block
	Foo() string
}

//go:generate basil generate BlockWithBlockInterface
type BlockWithBlockInterface struct {
	IDField    basil.ID          `basil:"id"`
	BlockNodes []basil.BlockNode `basil:"node"`
	Blocks     []BlockInterface  `basil:"block"`
}

func (b *BlockWithBlockInterface) ID() basil.ID {
	return b.IDField
}

func (b *BlockWithBlockInterface) Type() string {
	return "block_with_block_interface"
}

func (b *BlockWithBlockInterface) Context(ctx interface{}) interface{} {
	return ctx
}

func (b *BlockWithBlockInterface) Registry() block.Registry {
	return block.Registry{
		"block_simple": BlockSimpleInterpreter{},
	}
}
