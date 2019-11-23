// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package fixtures

import "github.com/opsidian/basil/basil"

//go:generate basil generate
type BlockGenerator struct {
	id     basil.ID              `basil:"id"`
	items  []interface{}         `basil:"required"`
	result *BlockGeneratorResult `basil:"generated"`
}

func (b *BlockGenerator) ID() basil.ID {
	return b.id
}

func (b *BlockGenerator) Main(blockCtx basil.BlockContext) error {
	for _, item := range b.items {
		res := &BlockGeneratorResult{
			id:    b.result.id,
			value: item,
		}
		if _, err := blockCtx.PublishBlock(res, nil); err != nil {
			return err
		}
	}
	return nil
}

//go:generate basil generate
type BlockGeneratorResult struct {
	id    basil.ID `basil:"id"`
	value interface{}
}

func (b *BlockGeneratorResult) ID() basil.ID {
	return b.id
}
