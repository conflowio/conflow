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

// ParameterNode is the AST node for a parameter
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . ParameterNode
type ParameterNode interface {
	Node
	parsley.StaticCheckable
	Name() ID
	ValueNode() parsley.Node
	IsDeclaration() bool
	SetSchema(schema.Schema)
}

// ParameterContainer is a parameter container
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . ParameterContainer
type ParameterContainer interface {
	Container
	BlockContainer() BlockContainer
}
