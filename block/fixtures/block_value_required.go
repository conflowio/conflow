package fixtures

import "github.com/opsidian/basil/variable"

//go:generate basil generate
type BlockValueRequired struct {
	IDField variable.ID `basil:"id"`
	Value   interface{} `basil:"value,required"`
}

func (b *BlockValueRequired) ID() variable.ID {
	return b.IDField
}

func (b *BlockValueRequired) Type() string {
	return "block_value_required"
}

func (b *BlockValueRequired) Context(ctx interface{}) interface{} {
	return ctx
}
