// Copyright (c) 2018 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package function

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/parsley/parsley"
)

// Interpreter defines an interpreter for functions
type Interpreter interface {
	StaticCheck(ctx interface{}, node basil.FunctionNode) (string, parsley.Error)
	Eval(ctx interface{}, node basil.FunctionNode) (interface{}, parsley.Error)
}
