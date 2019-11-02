// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

import "github.com/opsidian/parsley/parsley"

// ParameterDescriptor describes a parameter
type ParameterDescriptor struct {
	Type       string
	EvalStage  EvalStage
	IsRequired bool
	IsOutput   bool
}

// ParameterNode is the AST node for a parameter
//go:generate counterfeiter . ParameterNode
type ParameterNode interface {
	Node
	Name() ID
	ValueNode() parsley.Node
	IsDeclaration() bool
	SetDescriptor(ParameterDescriptor)
}

// ParameterContainer is a parameter container
//go:generate counterfeiter . ParameterContainer
type ParameterContainer interface {
	Container
	BlockContainer() BlockContainer
}
