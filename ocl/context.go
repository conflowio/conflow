// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ocl

// Context is the evaluation context
type Context struct {
	variableProvider VariableProvider
	functionRegistry FunctionRegistry
	blockRegistry    BlockRegistry
}

// NewContext creates a new context
func NewContext(
	variableProvider VariableProvider,
	functionRegistry FunctionRegistry,
	blockRegistry BlockRegistry,
) *Context {
	return &Context{
		variableProvider: variableProvider,
		functionRegistry: functionRegistry,
		blockRegistry:    blockRegistry,
	}
}

func (c *Context) GetVariableProvider() VariableProvider {
	return c.variableProvider
}

func (c *Context) GetFunctionRegistry() FunctionRegistry {
	return c.functionRegistry
}

func (c *Context) GetBlockRegistry() BlockRegistry {
	return c.blockRegistry
}
