// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package fixtures

import (
	"context"

	"github.com/conflowio/conflow/conflow"
)

// @block "generator"
type BlockGenerator struct {
	// @id
	id conflow.ID
	// @required
	items []interface{}
	// @generated
	result *BlockGeneratorResult
	// @dependency
	blockPublisher conflow.BlockPublisher
}

func (b *BlockGenerator) ID() conflow.ID {
	return b.id
}

func (b *BlockGenerator) Run(ctx context.Context) (conflow.Result, error) {
	for _, item := range b.items {
		res := &BlockGeneratorResult{
			id:    b.result.id,
			value: item,
		}
		if _, err := b.blockPublisher.PublishBlock(res, nil); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

// @block "configuration"
type BlockGeneratorResult struct {
	// @id
	id    conflow.ID
	value interface{}
}

func (b *BlockGeneratorResult) ID() conflow.ID {
	return b.id
}
