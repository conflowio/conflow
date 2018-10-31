package fixtures

import (
	"github.com/opsidian/ocl/ocl"
)

type BlockFactoryInterface interface {
	ocl.BlockFactory
	Foo() string
}

//go:generate ocl generate BlockWithFactoryInterface
type BlockWithFactoryInterface struct {
	IDField        string                  `ocl:"id"`
	BlockFactories []BlockFactoryInterface `ocl:"factory"`
}

func (b *BlockWithFactoryInterface) ID() string {
	return b.IDField
}

func (b *BlockWithFactoryInterface) Type() string {
	return "block_with_factory_interface"
}

func (b *BlockWithFactoryInterface) Context(ctx interface{}) interface{} {
	return ctx
}
