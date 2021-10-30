// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package fixtures

import (
	"github.com/opsidian/conflow/basil"
	"github.com/opsidian/conflow/basil/block"
)

// @block
type BlockWithOneBlock struct {
	// @id
	IDField basil.ID
	Block   *Block
}

func (b *BlockWithOneBlock) ID() basil.ID {
	return b.IDField
}

func (b *BlockWithOneBlock) ParseContextOverride() basil.ParseContextOverride {
	return basil.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"block": BlockInterpreter{},
		},
	}
}
