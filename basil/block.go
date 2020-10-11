// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

import (
	"github.com/opsidian/parsley/parsley"
)

// Block field tag constants
const (
	BlockTagBlock      ID = "block"
	BlockTagDeprecated    = "deprecated"
	BlockTagGenerated     = "generated"
	BlockTagID            = "id"
	BlockTagIgnore        = "ignore"
	BlockTagName          = "name"
	BlockTagOutput        = "output"
	BlockTagReference     = "reference"
	BlockTagRequired      = "required"
	BlockTagStage         = "stage"
	BlockTagValue         = "value"
	BlockTagDefault       = "default"
)

// BlockTags contains the valid block tags with descriptions
var BlockTags = map[ID]string{
	BlockTagBlock:      "marks an array field which should store child blocks",
	BlockTagDeprecated: "marks the field as deprecated (for documentation purposes)",
	BlockTagGenerated:  "marks the block as generated",
	BlockTagID:         "marks the id field in the block",
	BlockTagIgnore:     "the field is ignored when processing the block",
	BlockTagName:       "overrides the parameter name, otherwise the field name will be converted to under_score",
	BlockTagOutput:     "marks the field as output",
	BlockTagReference:  "marks the field that it must reference an existing identifier",
	BlockTagRequired:   "marks the field as required (must be set but can be empty)",
	BlockTagStage:      "sets the evaluation stage for the field",
	BlockTagValue:      "sets the field as the value field to be used for the short block format",
	BlockTagDefault:    "sets the default value for the field",
}

// Block is an interface for a block object
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Block
type Block interface {
	Identifiable
}

// BlockDescriptor describes a block
type BlockDescriptor struct {
	EvalStage   EvalStage
	IsRequired  bool
	IsGenerated bool
	IsMany      bool
}

// BlockContainer is a simple wrapper around a block object
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . BlockContainer
type BlockContainer interface {
	Container
	Param(ID) interface{}
	SetChild(Container)
	SetError(parsley.Error)
	PublishBlock(Block, func() error) (bool, error)
	EvalStageAware
}

// BlockInitialiser defines an Init() function which runs before the main evaluation stage
// If the skipped return value is true then the block won't be evaluated
type BlockInitialiser interface {
	Init(blockCtx BlockContext) (skipped bool, err error)
}

// BlockRunner defines a Main() function which runs the main business logic
type BlockRunner interface {
	Main(blockCtx BlockContext) error
}

// BlockCloser defines a Close function which runs after the main evaluation stage
type BlockCloser interface {
	Close(blockCtx BlockContext) error
}

type EvalStageAware interface {
	EvalStage() EvalStage
}

// BlockNode is the AST node for a block
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . BlockNode
type BlockNode interface {
	Node
	Children() []Node
	BlockType() ID
	ParamType(ID) (string, bool)
	Interpreter() BlockInterpreter
	SetDescriptor(BlockDescriptor)
}

// BlockNodeRegistry is an interface for looking up named blocks
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . BlockNodeRegistry
type BlockNodeRegistry interface {
	BlockNode(ID) (BlockNode, bool)
	AddBlockNode(BlockNode) error
}

// BlockTransformerRegistryAware is an interface to get a block node transformer registry
type BlockTransformerRegistryAware interface {
	BlockTransformerRegistry() parsley.NodeTransformerRegistry
}

// BlockInterpreter defines an interpreter for blocks
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . BlockInterpreter
type BlockInterpreter interface {
	CreateBlock(ID) Block
	SetParam(b Block, name ID, value interface{}) error
	SetBlock(b Block, name ID, value interface{}) error
	Param(b Block, name ID) interface{}
	Params() map[ID]ParameterDescriptor
	Blocks() map[ID]BlockDescriptor
	ValueParamName() ID
	HasForeignID() bool
	ParseContext(*ParseContext) *ParseContext
	EvalStageAware
}
