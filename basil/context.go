// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

// ContextAware returns a custom context
type ContextAware interface {
	Context(parentCtx interface{}) interface{}
}

// ContextConfig contains the parameters for a config
type ContextConfig struct {
	VariableProvider VariableProvider
	FunctionRegistry FunctionRegistry
	IDRegistry       IDRegistry
}

// Context is the evaluation context
type Context struct {
	variableProvider VariableProvider
	functionRegistry FunctionRegistry
	idRegistry       IDRegistry
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
		ctx.variableProvider = provider.GetVariableProvider()
	} else {
		panic("Variable provider must be set")
	}

	if config.FunctionRegistry != nil {
		ctx.functionRegistry = config.FunctionRegistry
	} else if provider, ok := parentCtx.(FunctionRegistryAware); ok {
		ctx.functionRegistry = provider.GetFunctionRegistry()
	} else {
		panic("Function registry must be set")
	}

	if config.IDRegistry != nil {
		ctx.idRegistry = config.IDRegistry
	} else if provider, ok := parentCtx.(IDRegistryAware); ok {
		ctx.idRegistry = provider.GetIDRegistry()
	} else {
		panic("ID registry must be set")
	}

	return ctx
}

// GetVariableProvider returns with the variable provider
func (c *Context) GetVariableProvider() VariableProvider {
	return c.variableProvider
}

// GetFunctionRegistry returns with the function registry
func (c *Context) GetFunctionRegistry() FunctionRegistry {
	return c.functionRegistry
}

// GetIDRegistry returns with the identifier registry
func (c *Context) GetIDRegistry() IDRegistry {
	return c.idRegistry
}
