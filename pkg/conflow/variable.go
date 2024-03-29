// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conflow

import "github.com/conflowio/parsley/parsley"

// VariableNode stores a variable reference. It always refers to a named block's parameter.
//
//counterfeiter:generate . VariableNode
type VariableNode interface {
	parsley.Node
	Identifiable
	ParentID() ID
	ParamName() ID
}
