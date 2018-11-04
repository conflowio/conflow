package basil

import "github.com/opsidian/parsley/parsley"

// BlockFactory defines an interface about how to create and evaluate blocks
type BlockFactory interface {
	HasForeignID() bool
	HasShortFormat() bool
	CreateBlock(parentCtx interface{}) (block Block, ctx interface{}, err parsley.Error)
	EvalBlock(ctx interface{}, stage string, res Block) parsley.Error
}

// BlockFactoryCreator defines a function to create a BlockFactory
type BlockFactoryCreator interface {
	CreateBlockFactory(
		typeNode parsley.Node,
		idNode parsley.Node,
		paramNodes map[ID]parsley.Node,
		blockNodes []parsley.Node,
	) (BlockFactory, parsley.Error)
}

// BlockFactoryCreatorFunc is a function which implements the BlockFactoryCreator interface
type BlockFactoryCreatorFunc func(
	typeNode parsley.Node,
	idNode parsley.Node,
	paramNodes map[ID]parsley.Node,
	blockNodes []parsley.Node,
) (BlockFactory, parsley.Error)

// CreateBlockFactory creates a block factory
func (f BlockFactoryCreatorFunc) CreateBlockFactory(
	typeNode parsley.Node,
	idNode parsley.Node,
	paramNodes map[ID]parsley.Node,
	blockNodes []parsley.Node,
) (BlockFactory, parsley.Error) {
	return f(typeNode, idNode, paramNodes, blockNodes)
}
