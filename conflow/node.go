// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conflow

import (
	"github.com/opsidian/parsley/parsley"
)

// Node is an identifiable node which has dependencies and has an evaluation stage
//counterfeiter:generate . Node
type Node interface {
	parsley.Node
	Identifiable
	EvalStage() EvalStage
	Dependencies() Dependencies
	Directives() []BlockNode
	Provides() []ID
	Generates() []ID
	CreateContainer(
		ctx *EvalContext,
		runtimeConfig RuntimeConfig,
		parent BlockContainer,
		value interface{},
		wgs []WaitGroup,
		pending bool,
	) JobContainer
	Value(userCtx interface{}) (interface{}, parsley.Error)
}

// Dependencies is a variable list
type Dependencies map[ID]VariableNode

func (d Dependencies) Add(d2 Dependencies) {
	for k, v := range d2 {
		d[k] = v
	}
}
