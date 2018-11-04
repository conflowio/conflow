package fixtures

import (
	"github.com/opsidian/basil/basil"
)

type BlockFactoryInterface interface {
	basil.BlockFactory
	Foo() string
}

//go:generate basil generate BlockWithFactoryInterface
type BlockWithFactoryInterface struct {
	IDField        string                  `basil:"id"`
	BlockFactories []BlockFactoryInterface `basil:"factory"`
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
