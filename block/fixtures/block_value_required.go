package fixtures

import "github.com/opsidian/basil/basil"

//go:generate basil generate BlockValueRequired
type BlockValueRequired struct {
	IDField basil.ID    `basil:"id"`
	Value   interface{} `basil:"value,required"`
}

func (b *BlockValueRequired) ID() basil.ID {
	return b.IDField
}

func (b *BlockValueRequired) Type() string {
	return "block_value_required"
}

func (b *BlockValueRequired) Context(ctx interface{}) interface{} {
	return ctx
}
