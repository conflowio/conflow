// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

import "github.com/opsidian/parsley/parsley"

// VariableNode stores a variable reference. It always refers to a named block's parameter.
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . VariableNode
type VariableNode interface {
	parsley.Node
	Identifiable
	ParentID() ID
	ParamName() ID
}
