// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

import (
	"github.com/opsidian/parsley/parsley"

	"github.com/opsidian/basil/basil/schema"
)

// FunctionNameRegExpPattern defines a valid function name
const FunctionNameRegExpPattern = IDRegExpPattern + "(?:\\." + IDRegExpPattern + ")?"

// FunctionNode is the AST node for a function
//counterfeiter:generate . FunctionNode
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
//counterfeiter:generate . FunctionInterpreter
type FunctionInterpreter interface {
	Schema() schema.Schema
	Eval(ctx interface{}, node FunctionNode) (interface{}, parsley.Error)
}
