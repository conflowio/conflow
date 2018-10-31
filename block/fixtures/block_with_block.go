package fixtures

import (
	"github.com/opsidian/ocl/block"
	"github.com/opsidian/ocl/ocl"
)

//go:generate ocl generate BlockWithBlock
type BlockWithBlock struct {
	IDField        string                `ocl:"id"`
	BlockFactories []*BlockSimpleFactory `ocl:"factory"`
	Blocks         []*BlockSimple        `ocl:"block"`
}

func (b *BlockWithBlock) ID() string {
	return b.IDField
}

func (b *BlockWithBlock) Type() string {
	return "block_with_block"
}

func (b *BlockWithBlock) Context(ctx interface{}) interface{} {
	return ocl.NewContext(ctx, ocl.ContextConfig{
		BlockRegistry: block.Registry{
			"block_simple": ocl.BlockFactoryCreatorFunc(NewBlockSimpleFactory),
		},
	})
}
