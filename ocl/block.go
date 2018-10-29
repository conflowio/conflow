package ocl

import (
	"github.com/opsidian/parsley/parsley"
)

// Block is an interface for a block object
type Block interface {
	Identifiable
	TypeAware
}

// BlockFactory defines an interface about how to create and evaluate blocks
type BlockFactory interface {
	CreateBlock(ctx interface{}) Block
	EvalBlock(ctx interface{}, stage string, res interface{}) parsley.Error
}

// BlockPreInterpreter as an interface for running a pre-evaluation hook
type BlockPreInterpreter interface {
	PreEval(ctx interface{}, stage string) parsley.Error
}

// BlockPostInterpreter as an interface for running a post-evaluation hook
type BlockPostInterpreter interface {
	PostEval(ctx interface{}, stage string) parsley.Error
}

type BlockFactoryCreator interface {
	CreateBlockFactory(
		typeNode parsley.Node,
		idNode parsley.Node,
		paramNodes map[string]parsley.Node,
		blockNodes []parsley.Node,
	) (BlockFactory, parsley.Error)
}

type BlockFactoryCreatorFunc func(
	typeNode parsley.Node,
	idNode parsley.Node,
	paramNodes map[string]parsley.Node,
	blockNodes []parsley.Node,
) (BlockFactory, parsley.Error)

func (f BlockFactoryCreatorFunc) CreateBlockFactory(
	typeNode parsley.Node,
	idNode parsley.Node,
	paramNodes map[string]parsley.Node,
	blockNodes []parsley.Node,
) (BlockFactory, parsley.Error) {
	return f(typeNode, idNode, paramNodes, blockNodes)
}

// BlockRegistry is a list of BlockFactory creator objects
type BlockRegistry map[string]BlockFactoryCreator

func (b BlockRegistry) TypeExists(blockType string) bool {
	_, exists := b[blockType]
	return exists
}

// BlockRegistryAware defines a function to access a block registry
type BlockRegistryAware interface {
	GetBlockRegistry() BlockRegistry
}
