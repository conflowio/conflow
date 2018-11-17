package fixtures

import "github.com/opsidian/basil/basil"

//go:generate basil generate
type BlockSimple struct {
	IDField basil.ID    `basil:"id"`
	Value   interface{} `basil:"value"`
}

func (b *BlockSimple) ID() basil.ID {
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
