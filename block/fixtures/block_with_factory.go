package fixtures

//go:generate basil generate BlockWithFactory
type BlockWithFactory struct {
	IDField        string                `basil:"id"`
	BlockFactories []*BlockSimpleFactory `basil:"factory,stage=-"`
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
