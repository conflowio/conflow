// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blocks

import (
	"context"
	"fmt"

	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/block"
)

// @block "generator"
type Iterator struct {
	// @id
	id    conflow.ID
	value interface{}
	// @generated
	it *It
	// @dependency
	blockPublisher conflow.BlockPublisher
}

func (i *Iterator) ParseContextOverride() conflow.ParseContextOverride {
	return conflow.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"it": ItInterpreter{},
		},
	}
}

func (i *Iterator) ID() conflow.ID {
	return i.id
}

func (i *Iterator) Run(ctx context.Context) (conflow.Result, error) {
	switch val := i.value.(type) {
	case []interface{}:
		for k, v := range val {
			_, err := i.blockPublisher.PublishBlock(&It{
				id:    i.it.id,
				key:   k,
				value: v,
			}, nil)
			if err != nil {
				return nil, err
			}
		}
	case map[string]interface{}:
		for k, v := range val {
			_, err := i.blockPublisher.PublishBlock(&It{
				id:    i.it.id,
				key:   k,
				value: v,
			}, nil)
			if err != nil {
				return nil, err
			}
		}
	default:
		return nil, fmt.Errorf("invalid value for range: %T", i.value)
	}

	return nil, nil
}

// @block "configuration"
type It struct {
	// @id
	id conflow.ID
	// @read_only
	key interface{}
	// @read_only
	value interface{}
}

func (i *It) ID() conflow.ID {
	return i.id
}
