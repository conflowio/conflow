package basil

import (
	"github.com/opsidian/parsley/parsley"
)

// Block is an interface for a block object
//go:generate counterfeiter . Block
type Block interface {
	Identifiable
}

// BlockContainer is a simple wrapper around a block object
//go:generate counterfeiter . BlockContainer
type BlockContainer interface {
	Block() Block
	Param(ID) interface{}
	SetParam(ctx *EvalContext, name ID, node parsley.Node) parsley.Error
}

// BlockContainerRegistry stores block container instances
type BlockContainerRegistry interface {
	BlockContainer(ID) (BlockContainer, bool)
	AddBlockContainer(BlockContainer) error
}

// BlockContainerRegistryAware defines an interface to get a block registry
type BlockContainerRegistryAware interface {
	BlockContainerRegistry() BlockContainerRegistry
}

// BlockInitialiser Init() runs after the pre-evaluation and before the main evaluation stage
// If the boolean return value is false then the block won't be evaluated
type BlockInitialiser interface {
	Init(ctx *EvalContext) (bool, error)
}

// BlockRunner Run() runs after the main evaluation stage
type BlockRunner interface {
	Run(ctx *EvalContext) error
}

// BlockFinisher Finish() runs after the post-evaluation stage
type BlockFinisher interface {
	Finish(ctx *EvalContext) error
}
