package fixtures

import "github.com/opsidian/basil/basil"

type BlockInterface interface {
	basil.Block
	Foo() string
}

//go:generate basil generate BlockWithBlockInterface
type BlockWithBlockInterface struct {
	IDField        basil.ID             `basil:"id"`
	BlockFactories []basil.BlockFactory `basil:"factory"`
	Blocks         []BlockInterface     `basil:"block"`
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
