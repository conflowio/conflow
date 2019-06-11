// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package fixtures

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/block"
)

type BlockInterface interface {
	basil.Block
	Foo() string
}

//go:generate basil generate
type BlockWithBlockInterface struct {
	IDField     basil.ID         `basil:"id"`
	BlockSimple []BlockInterface `basil:"block"`
}

func (b *BlockWithBlockInterface) ID() basil.ID {
	return b.IDField
}

func (b *BlockWithBlockInterface) ParseContextOverride(ctx *basil.ParseContext) basil.ParseContextOverride {
	return basil.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"block_simple": BlockSimpleInterpreter{},
		},
	}
}
