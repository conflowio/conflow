package ocl

import "github.com/opsidian/parsley/parsley"

// Block is an interface for a block object
type Block interface {
	Identifiable
	TypeAware
}

// BlockFactory defines an interface about how to create blocks
type BlockFactory interface {
	CreateBlock(
		ctx interface{},
		typeNode parsley.Node,
		idNode parsley.Node,
		paramNodes map[string]parsley.Node,
		blockNodes []parsley.Node,
	) (Block, parsley.Error)
}

// BlockFactoryFunc defines a helper to implement the BlockFactory interface with functions
type BlockFactoryFunc func(
	ctx interface{},
	typeNode parsley.Node,
	idNode parsley.Node,
	paramNodes map[string]parsley.Node,
	blockNodes []parsley.Node,
) (Block, parsley.Error)

// CreateBlock calls the function
func (f BlockFactoryFunc) CreateBlock(
	ctx interface{},
	typeNode parsley.Node,
	idNode parsley.Node,
	paramNodes map[string]parsley.Node,
	blockNodes []parsley.Node,
) (Block, parsley.Error) {
	return f(ctx, typeNode, idNode, paramNodes, blockNodes)
}

// BlockRegistry is an interface for a block registry
//go:generate counterfeiter . BlockRegistry
type BlockRegistry interface {
	BlockFactory
	BlockFactoryExists(blockType string) bool
	RegisterBlockFactory(blockType string, factory BlockFactory)
	BlockIDExists(string) bool
	GenerateBlockID() string
}

// BlockRegistryAware defines a function to access a block registry
type BlockRegistryAware interface {
	GetBlockRegistry() BlockRegistry
}
