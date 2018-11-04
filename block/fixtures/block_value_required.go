package fixtures

//go:generate basil generate BlockValueRequired
type BlockValueRequired struct {
	IDField string      `basil:"id"`
	Value   interface{} `basil:"value,required"`
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
