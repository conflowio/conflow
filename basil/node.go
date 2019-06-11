// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

import "github.com/opsidian/parsley/parsley"

// Node is an identifiable node which has dependencies and has an evaluation stage
//go:generate counterfeiter . Node
type Node interface {
	parsley.Node
	Identifiable
	EvalStage() EvalStage
	Dependencies() Dependencies
	Provides() []ID
	Generates() []ID
	Generated() bool
}

// Dependencies is a variable list
type Dependencies map[ID]VariableNode
