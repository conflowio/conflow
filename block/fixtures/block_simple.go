package fixtures

import "github.com/opsidian/basil/variable"

//go:generate basil generate BlockSimple
type BlockSimple struct {
	IDField variable.ID `basil:"id"`
	Value   interface{} `basil:"value"`
}

func (b *BlockSimple) ID() variable.ID {
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
