// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package fixtures

import (
	"context"

	"github.com/opsidian/basil/basil"
)

// @block
type BlockGenerator struct {
	// @id
	id basil.ID
	// @required
	items []interface{}
	// @generated
	result *BlockGeneratorResult
	// @dependency
	blockPublisher basil.BlockPublisher
}

func (b *BlockGenerator) ID() basil.ID {
	return b.id
}

func (b *BlockGenerator) Run(ctx context.Context) error {
	for _, item := range b.items {
		res := &BlockGeneratorResult{
			id:    b.result.id,
			value: item,
		}
		if _, err := b.blockPublisher.PublishBlock(res, nil); err != nil {
			return err
		}
	}
	return nil
}

// @block
type BlockGeneratorResult struct {
	// @id
	id    basil.ID
	value interface{}
}

func (b *BlockGeneratorResult) ID() basil.ID {
	return b.id
}
