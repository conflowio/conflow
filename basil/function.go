// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

import (
	"github.com/opsidian/parsley/parsley"
)

// FunctionNameRegExpPattern defines a valid function name
const FunctionNameRegExpPattern = IDRegExpPattern + "(?:\\." + IDRegExpPattern + ")?"

// FunctionNode is the AST node for a function
//go:generate counterfeiter . FunctionNode
type FunctionNode interface {
	parsley.Node
	parsley.StaticCheckable
	Name() ID
	ArgumentNodes() []parsley.Node
}

// FunctionTransformerRegistryAware is an interface to get a function node transformer registry
type FunctionTransformerRegistryAware interface {
	FunctionTransformerRegistry() parsley.NodeTransformerRegistry
}

// FunctionInterpreter defines an interpreter for functions
//go:generate counterfeiter . FunctionInterpreter
type FunctionInterpreter interface {
	StaticCheck(ctx interface{}, node FunctionNode) (string, parsley.Error)
	Eval(ctx interface{}, node FunctionNode) (interface{}, parsley.Error)
}
