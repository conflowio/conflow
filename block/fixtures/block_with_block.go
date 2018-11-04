package fixtures

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/block"
)

//go:generate basil generate BlockWithBlock
type BlockWithBlock struct {
	IDField        string                `basil:"id"`
	BlockFactories []*BlockSimpleFactory `basil:"factory"`
	Blocks         []*BlockSimple        `basil:"block"`
}

func (b *BlockWithBlock) ID() string {
	return b.IDField
}

func (b *BlockWithBlock) Type() string {
	return "block_with_block"
}

func (b *BlockWithBlock) Context(ctx interface{}) interface{} {
	return basil.NewContext(ctx, basil.ContextConfig{
		BlockRegistry: block.Registry{
			"block_simple": basil.BlockFactoryCreatorFunc(NewBlockSimpleFactory),
		},
	})
}
