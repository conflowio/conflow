// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ocl

import (
	"github.com/opsidian/ocl/function"
	"github.com/opsidian/ocl/variable"
	"github.com/opsidian/parsley/parsley"
)

// Context is the evaluation context
type Context struct {
	variableProvider variable.Provider
	functionRegistry function.Registry
}

// NewContext creates a new context
func NewContext(variableProvider variable.Provider, functionRegistry function.Registry) *Context {
	return &Context{
		variableProvider: variableProvider,
		functionRegistry: functionRegistry,
	}
}

// GetVar returns with the named variable
func (c *Context) GetVar(name string) (interface{}, bool) {
	return c.variableProvider.GetVar(name)
}

// LookupVar searches for the given complex variable
func (c *Context) LookupVar(lookup variable.LookUp) (interface{}, error) {
	return c.variableProvider.LookupVar(lookup)
}

// CallFunction calls the named function with the given parameters
func (c *Context) CallFunction(ctx interface{}, function parsley.Node, params []parsley.Node) (interface{}, parsley.Error) {
	return c.functionRegistry.CallFunction(ctx, function, params)
}

// FunctionExists checks whether the function exists
func (c *Context) FunctionExists(name string) bool {
	return c.functionRegistry.FunctionExists(name)
}

// RegisterFunction registers a new function
func (c *Context) RegisterFunction(name string, callable function.Callable) {
	c.functionRegistry.RegisterFunction(name, callable)
}
