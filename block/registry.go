package block

import (
	"fmt"

	"github.com/opsidian/parsley/parsley"

	"github.com/opsidian/ocl/ocl"
)

// Registry is a list of BlockFactory creator objects
type Registry map[string]ocl.BlockFactoryCreator

// AddBlockFactoryCreator registers a new block factory creator
func (r Registry) AddBlockFactoryCreator(blockType string, creator ocl.BlockFactoryCreator) {
	_, exists := r[blockType]
	if exists {
		panic(fmt.Sprintf("%s block factory creator was already registered", blockType))
	}
	r[blockType] = creator
}

// BlockTypeExists returns true if the registry has a factory creator for the given type
func (r Registry) BlockTypeExists(blockType string) bool {
	_, exists := r[blockType]
	return exists
}

// CreateBlockFactory creates a block factory
func (r Registry) CreateBlockFactory(
	ctx interface{},
	typeNode parsley.Node,
	idNode parsley.Node,
	paramNodes map[string]parsley.Node,
	blockNodes []parsley.Node,
) (ocl.BlockFactory, parsley.Error) {
	blockType, err := typeNode.Value(ctx)
	if err != nil {
		return nil, err
	}

	creator, exists := r[blockType.(string)]
	if !exists {
		return nil, parsley.NewErrorf(typeNode.Pos(), "%q type is invalid or not allowed here", blockType)
	}
	return creator.CreateBlockFactory(typeNode, idNode, paramNodes, blockNodes)
}
