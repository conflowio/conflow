package common

import (
	"fmt"

	"github.com/opsidian/basil/basil"
)

// Println will write a string followed by a new line to the standard output
//go:generate basil generate
type Println struct {
	id    basil.ID    `basil:"id"`
	value interface{} `basil:"value,required"`
}

func (p *Println) Main(ctx basil.BlockContext) error {
	fmt.Println(p.value)
	return nil
}
