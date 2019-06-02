package basil

import (
	"github.com/opsidian/parsley/parsley"
)

// Block field tag constants
const (
	BlockTagBlock      = "block"
	BlockTagDeprecated = "deprecated"
	BlockTagID         = "id"
	BlockTagIgnore     = "ignore"
	BlockTagName       = "name"
	BlockTagNode       = "node"
	BlockTagOut        = "out"
	BlockTagReference  = "reference"
	BlockTagRequired   = "required"
	BlockTagStage      = "stage"
	BlockTagValue      = "value"
)

// BlockTags contains the valid block tags with descriptions
var BlockTags = map[string]string{
	BlockTagBlock:      "marks an array field which should store child blocks",
	BlockTagDeprecated: "marks the field as deprecated (for documentation purposes)",
	BlockTagID:         "marks the id field in the block",
	BlockTagIgnore:     "the field is ignored when processing the block",
	BlockTagName:       "overrides the parameter name, otherwise the field name will be converted to under_score",
	BlockTagNode:       "marks an array field which should store child block AST nodes",
	BlockTagOut:        "marks the field as output",
	BlockTagReference:  "marks the field that it must reference an existing identifier",
	BlockTagRequired:   "marks the field as required (must be set but can be empty)",
	BlockTagStage:      "sets the evaluation stage for the field",
	BlockTagValue:      "sets the field as the value field to be used for the short block format",
}

// Block is an interface for a block object
//go:generate counterfeiter . Block
type Block interface {
}

// BlockContainer is a simple wrapper around a block object
//go:generate counterfeiter . BlockContainer
type BlockContainer interface {
	Container
	Node() BlockNode
	Param(ID) interface{}
}

// BlockInitialiser defines an Init() function which runs before the main evaluation stage
// If the boolean return value is false then the block won't be evaluated
type BlockInitialiser interface {
	Init(blockCtx BlockContext) (bool, error)
}

// BlockRunner defines a Main() function which runs the main business logic
type BlockRunner interface {
	Main(blockCtx BlockContext) error
}

// BlockCloser defines a Close function which runs after the main evaluation stage
type BlockCloser interface {
	Close(blockCtx BlockContext) error
}

// BlockNode is the AST node for a block
//go:generate counterfeiter . BlockNode
type BlockNode interface {
	Node
	Children() []Node
	BlockType() ID
	ParamType(ID) (string, bool)
	Interpreter() BlockInterpreter
}

// BlockNodeRegistry is an interface for looking up named blocks
//go:generate counterfeiter . BlockNodeRegistry
type BlockNodeRegistry interface {
	BlockNode(ID) (BlockNode, bool)
	AddBlockNode(BlockNode) error
}

// BlockTransformerRegistryAware is an interface to get a block node transformer registry
type BlockTransformerRegistryAware interface {
	BlockTransformerRegistry() parsley.NodeTransformerRegistry
}

// BlockInterpreter defines an interpreter for blocks
//go:generate counterfeiter . BlockInterpreter
type BlockInterpreter interface {
	CreateBlock(ID) Block
	SetParam(b Block, name ID, value interface{}) error
	SetBlock(b Block, name ID, value interface{}) error
	Param(b Block, name ID) interface{}
	Params() map[ID]ParameterDescriptor
	ValueParamName() ID
	HasForeignID() bool
	ParseContext(*ParseContext) *ParseContext
}
