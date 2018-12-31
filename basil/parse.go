package basil

import (
	"github.com/opsidian/parsley/parsley"
)

// ParseContext is the parsing context
type ParseContext struct {
	blockTransformerRegistry    parsley.NodeTransformerRegistry
	functionTransformerRegistry parsley.NodeTransformerRegistry
	idRegistry                  IDRegistry
	blockNodeRegistry           BlockNodeRegistry
}

// ParseContextAware defines an interface to return a child parsing content
type ParseContextAware interface {
	ParseContext(*ParseContext) *ParseContext
}

// ParseContextConfig stores override values for a child context
type ParseContextConfig struct {
	BlockTransformerRegistry    parsley.NodeTransformerRegistry
	FunctionTransformerRegistry parsley.NodeTransformerRegistry
}

// NewParseContext returns with a new parsing context
func NewParseContext(
	blockTransformerRegistry parsley.NodeTransformerRegistry,
	functionTransformerRegistry parsley.NodeTransformerRegistry,
	idRegistry IDRegistry,
	blockNodeRegistry BlockNodeRegistry,
) *ParseContext {
	return &ParseContext{
		blockTransformerRegistry:    blockTransformerRegistry,
		functionTransformerRegistry: functionTransformerRegistry,
		idRegistry:                  idRegistry,
		blockNodeRegistry:           blockNodeRegistry,
	}
}

// New creates a new child context
func (p *ParseContext) New(config ParseContextConfig) *ParseContext {
	ctx := &ParseContext{
		idRegistry:        p.idRegistry,
		blockNodeRegistry: p.blockNodeRegistry,
	}
	if config.BlockTransformerRegistry != nil {
		ctx.blockTransformerRegistry = config.BlockTransformerRegistry
	} else {
		ctx.blockTransformerRegistry = p.blockTransformerRegistry
	}
	if config.FunctionTransformerRegistry != nil {
		ctx.functionTransformerRegistry = config.FunctionTransformerRegistry
	} else {
		ctx.functionTransformerRegistry = p.functionTransformerRegistry
	}

	return ctx
}

// BlockTransformerRegistry returns with the block node transformer registry
func (p *ParseContext) BlockTransformerRegistry() parsley.NodeTransformerRegistry {
	return p.blockTransformerRegistry
}

// FunctionTransformerRegistry returns with the function node transformer registry
func (p *ParseContext) FunctionTransformerRegistry() parsley.NodeTransformerRegistry {
	return p.functionTransformerRegistry
}

// IDRegistry returns with the identifier registry
func (p *ParseContext) IDRegistry() IDRegistry {
	return p.idRegistry
}

// BlockNodeRegistry returns with the the block node registry
func (p *ParseContext) BlockNodeRegistry() BlockNodeRegistry {
	return p.blockNodeRegistry
}
