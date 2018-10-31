package fixtures

import "github.com/opsidian/ocl/ocl"

type BlockInterface interface {
	ocl.Block
	Foo() string
}

//go:generate ocl generate BlockWithBlockInterface
type BlockWithBlockInterface struct {
	IDField        string             `ocl:"id"`
	BlockFactories []ocl.BlockFactory `ocl:"factory"`
	Blocks         []BlockInterface   `ocl:"block"`
}

func (b *BlockWithBlockInterface) ID() string {
	return b.IDField
}

func (b *BlockWithBlockInterface) Type() string {
	return "block_with_block_interface"
}

func (b *BlockWithBlockInterface) Context(ctx interface{}) interface{} {
	return ctx
}
