package fixtures

import "github.com/opsidian/basil/variable"

//go:generate basil generate
type BlockWithReference struct {
	IDField variable.ID `basil:"id,reference"`
}

func (b *BlockWithReference) ID() variable.ID {
	return b.IDField
}

func (b *BlockWithReference) Type() string {
	return "block_value_required"
}

func (b *BlockWithReference) Context(ctx interface{}) interface{} {
	return ctx
}
