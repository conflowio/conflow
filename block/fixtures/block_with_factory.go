package fixtures

//go:generate ocl generate BlockWithFactory
type BlockWithFactory struct {
	IDField        string                `ocl:"id"`
	BlockFactories []*BlockSimpleFactory `ocl:"factory"`
}

func (b *BlockWithFactory) ID() string {
	return b.IDField
}

func (b *BlockWithFactory) Type() string {
	return "block_with_factory"
}

func (b *BlockWithFactory) Context(ctx interface{}) interface{} {
	return ctx
}
