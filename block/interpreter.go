// Copyright (c) 2018 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package block

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/variable"
	"github.com/opsidian/parsley/parsley"
)

// Interpreter defines an interpreter for blocks
type Interpreter interface {
	StaticCheck(ctx interface{}, node basil.BlockNode) (string, parsley.Error)
	Eval(ctx interface{}, node basil.BlockNode) (block basil.Block, err parsley.Error)
	EvalBlock(ctx interface{}, node basil.BlockNode, stage string, block basil.Block) parsley.Error
	ValueParamName() variable.ID
	HasForeignID() bool
	basil.BlockRegistryAware
}
