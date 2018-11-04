package fixtures

import "github.com/opsidian/basil/basil"

//go:generate basil generate BlockWithFactory
type BlockWithFactory struct {
	IDField        basil.ID              `basil:"id"`
	BlockFactories []*BlockSimpleFactory `basil:"factory,stage=-"`
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
