package fixtures

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/block"
)

//go:generate basil generate
type BlockWithNode struct {
	IDField    basil.ID          `basil:"id"`
	BlockNodes []basil.BlockNode `basil:"node,stage=-"`
}

func (b *BlockWithNode) ID() basil.ID {
	return b.IDField
}

func (b *BlockWithNode) Type() string {
	return "block_with_node"
}

func (b *BlockWithNode) Context(ctx interface{}) interface{} {
	return ctx
}

func (b *BlockWithNode) BlockRegistry() block.Registry {
	return block.Registry{
		"block_simple": BlockSimpleInterpreter{},
	}
}
