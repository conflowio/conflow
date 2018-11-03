package fixtures

//go:generate ocl generate BlockValueRequired
type BlockValueRequired struct {
	IDField string      `ocl:"id"`
	Value   interface{} `ocl:"value,required"`
}

func (b *BlockValueRequired) ID() string {
	return b.IDField
}

func (b *BlockValueRequired) Type() string {
	return "block_value_required"
}

func (b *BlockValueRequired) Context(ctx interface{}) interface{} {
	return ctx
}
