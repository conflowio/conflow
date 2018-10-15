// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ocl

import (
	"github.com/opsidian/ocl/variable"
)

// Context is the evaluation context
type Context struct {
	variableProvider variable.Provider
}

// NewContext creates a new context
func NewContext(variableProvider variable.Provider) *Context {
	return &Context{
		variableProvider: variableProvider,
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
