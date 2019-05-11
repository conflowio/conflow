package fixtures

import "github.com/opsidian/basil/basil"

//go:generate basil generate
type BlockSimple struct {
	IDField basil.ID    `basil:"id"`
	Value   interface{} `basil:"value"`
}

func (b *BlockSimple) Foo() string {
	return "bar"
}
