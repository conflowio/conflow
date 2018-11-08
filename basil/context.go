// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

import (
	"github.com/opsidian/parsley/parsley"
)

// ContextAware returns a custom context
type ContextAware interface {
	Context(parentCtx interface{}) interface{}
}

// ContextConfig contains the parameters for a config
type ContextConfig struct {
	VariableProvider VariableProvider
	IDRegistry       IDRegistry
	BlockRegistry    parsley.NodeTransformerRegistry
	FunctionRegistry parsley.NodeTransformerRegistry
}

// Context is the evaluation context
type Context struct {
	variableProvider VariableProvider
	idRegistry       IDRegistry
	blockRegistry    parsley.NodeTransformerRegistry
	functionRegistry parsley.NodeTransformerRegistry
}

// NewContext creates a new context
func NewContext(
	parentCtx interface{},
	config ContextConfig,
) *Context {
	ctx := &Context{}

	if config.VariableProvider != nil {
		ctx.variableProvider = config.VariableProvider
	} else if provider, ok := parentCtx.(VariableProviderAware); ok {
		ctx.variableProvider = provider.VariableProvider()
	} else {
		panic("Variable provider must be set")
	}

	if config.IDRegistry != nil {
		ctx.idRegistry = config.IDRegistry
	} else if provider, ok := parentCtx.(IDRegistryAware); ok {
		ctx.idRegistry = provider.IDRegistry()
	} else {
		panic("ID registry must be set")
	}

	if config.BlockRegistry != nil {
		ctx.blockRegistry = config.BlockRegistry
	} else if provider, ok := parentCtx.(BlockRegistryAware); ok {
		ctx.blockRegistry = provider.BlockRegistry()
	} else {
		panic("Block registry must be set")
	}

	if config.FunctionRegistry != nil {
		ctx.functionRegistry = config.FunctionRegistry
	} else if provider, ok := parentCtx.(FunctionRegistryAware); ok {
		ctx.functionRegistry = provider.FunctionRegistry()
	} else {
		panic("Function registry must be set")
	}

	return ctx
}

// VariableProvider returns with the variable provider
func (c *Context) VariableProvider() VariableProvider {
	return c.variableProvider
}

// IDRegistry returns with the identifier registry
func (c *Context) IDRegistry() IDRegistry {
	return c.idRegistry
}

// BlockRegistry returns with the block node transformer registry
func (c *Context) BlockRegistry() parsley.NodeTransformerRegistry {
	return c.blockRegistry
}

// FunctionRegistry returns with the function node transformer registry
func (c *Context) FunctionRegistry() parsley.NodeTransformerRegistry {
	return c.functionRegistry
}
