// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conflow

import (
	"context"

	"github.com/opsidian/parsley/parsley"

	"github.com/opsidian/conflow/basil/schema"
)

// Block is an interface for a block object
type Block interface {
}

// BlockContainer is a simple wrapper around a block object
//counterfeiter:generate . BlockContainer
type BlockContainer interface {
	Container
	Param(ID) interface{}
	SetChild(Container)
	SetError(parsley.Error)
	EvalStage() EvalStage
}

// BlockInitialiser defines an Init() function which runs before the main evaluation stage
// If the skipped return value is true then the block won't be evaluated
type BlockInitialiser interface {
	Init(context.Context) (skipped bool, err error)
}

// BlockRunner defines a Run() function which runs the main business logic
type BlockRunner interface {
	Run(ctx context.Context) (Result, error)
}

// BlockCloser defines a Close function which runs after the main evaluation stage
type BlockCloser interface {
	Close(ctx context.Context) error
}

// BlockPublisher defines an interface which generator blocks should use to publish generated blocks
// The PublishBlock function will either:
//  * return immediately with published=false if the block is not a dependency of other blocks
//  * otherwise it will block until all other blocks depending on the published block will finish running
//
// If the onScheduled function is not nil, it will be called after the published block was scheduled
type BlockPublisher interface {
	PublishBlock(block Block, onScheduled func() error) (published bool, err error)
}

// BlockNode is the AST node for a block
//counterfeiter:generate . BlockNode
type BlockNode interface {
	Node
	Children() []Node
	ParameterName() ID
	BlockType() ID
	Interpreter() BlockInterpreter
	SetSchema(schema.Schema)
	GetPropertySchema(ID) (schema.Schema, bool)
}

// BlockNodeRegistry is an interface for looking up named blocks
//counterfeiter:generate . BlockNodeRegistry
type BlockNodeRegistry interface {
	BlockNode(ID) (BlockNode, bool)
	AddBlockNode(BlockNode) error
}

// BlockTransformerRegistryAware is an interface to get a block node transformer registry
type BlockTransformerRegistryAware interface {
	BlockTransformerRegistry() parsley.NodeTransformerRegistry
}

// BlockInterpreter defines an interpreter for blocks
//counterfeiter:generate . BlockInterpreter
type BlockInterpreter interface {
	Schema() schema.Schema
	CreateBlock(ID, *BlockContext) Block
	SetParam(b Block, name ID, value interface{}) error
	SetBlock(b Block, name ID, value interface{}) error
	Param(b Block, name ID) interface{}
	ValueParamName() ID
	ParseContext(*ParseContext) *ParseContext
}

// BlockProvider is an interface for an object which provides additional block types
type BlockProvider interface {
	BlockInterpreters(*ParseContext) (map[ID]BlockInterpreter, error)
}
