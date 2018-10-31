package fixtures

//go:generate ocl generate BlockSimple
type BlockSimple struct {
	IDField string `ocl:"id"`
}

func (b *BlockSimple) ID() string {
	return b.IDField
}

func (b *BlockSimple) Type() string {
	return "block_simple"
}

func (b *BlockSimple) Context(ctx interface{}) interface{} {
	return ctx
}

func (b *BlockSimple) Foo() string {
	return "bar"
}
