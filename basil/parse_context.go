// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

import (
	"fmt"

	"github.com/opsidian/parsley/parsley"
)

// ParseContext is the parsing context
type ParseContext struct {
	blockTransformerRegistry    parsley.NodeTransformerRegistry
	functionTransformerRegistry parsley.NodeTransformerRegistry
	idRegistry                  IDRegistry
	blockNodes                  map[ID]BlockNode
	fileSet                     *parsley.FileSet
}

// ParseContextOverride stores override values for a parse context
type ParseContextOverride struct {
	BlockTransformerRegistry    parsley.NodeTransformerRegistry
	FunctionTransformerRegistry parsley.NodeTransformerRegistry
}

// ParseContextOverrider defines an interface to be able to override a parse config
type ParseContextOverrider interface {
	ParseContextOverride() ParseContextOverride
}

// NewParseContext returns with a new parsing context
func NewParseContext(idRegistry IDRegistry) *ParseContext {
	return &ParseContext{
		idRegistry: idRegistry,
		blockNodes: make(map[ID]BlockNode, 32),
		fileSet:    parsley.NewFileSet(),
	}
}

// New creates a new child context
func (p *ParseContext) New(config ParseContextOverride) *ParseContext {
	ctx := &ParseContext{
		idRegistry: p.idRegistry,
		blockNodes: p.blockNodes,
		fileSet:    p.fileSet,
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

// BlockNode returns with the given block node if it exists
func (p *ParseContext) BlockNode(id ID) (BlockNode, bool) {
	node, ok := p.blockNodes[id]
	return node, ok
}

// AddBlockNode adds a new block node
// It returns with an error if a block with the same id was already registered
func (p *ParseContext) AddBlockNode(node BlockNode) error {
	if _, exists := p.blockNodes[node.ID()]; exists {
		return fmt.Errorf("%q is already defined, please use a globally unique identifier", node.ID())
	}

	p.blockNodes[node.ID()] = node

	return nil
}

// IDExists returns true if the identifier already exists
func (p *ParseContext) IDExists(id ID) bool {
	return p.idRegistry.IDExists(id)
}

// GenerateID generates a new id
func (p *ParseContext) GenerateID() ID {
	return p.idRegistry.GenerateID()
}

// RegisterID registers a new id and returns an error if it was already registered
func (p *ParseContext) RegisterID(id ID) error {
	return p.idRegistry.RegisterID(id)
}

// FileSet returns with the file set
func (p *ParseContext) FileSet() *parsley.FileSet {
	return p.fileSet
}
