package fixtures

import "github.com/opsidian/basil/basil"

//go:generate basil generate BlockWithReference
type BlockWithReference struct {
	IDField basil.ID `basil:"id,reference"`
}

func (b *BlockWithReference) ID() basil.ID {
	return b.IDField
}

func (b *BlockWithReference) Type() string {
	return "block_value_required"
}

func (b *BlockWithReference) Context(ctx interface{}) interface{} {
	return ctx
}
