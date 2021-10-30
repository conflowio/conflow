// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blocks

import (
	"context"
	"fmt"

	"github.com/opsidian/conflow/basil/block"
	"github.com/opsidian/conflow/conflow"
)

// @block
type Range struct {
	// @id
	id    conflow.ID
	value interface{}
	// @generated
	entry *RangeEntry
	// @dependency
	blockPublisher conflow.BlockPublisher
}

func (r *Range) ParseContextOverride() conflow.ParseContextOverride {
	return conflow.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"entry": RangeEntryInterpreter{},
		},
	}
}

func (r *Range) ID() conflow.ID {
	return r.id
}

func (r *Range) Run(ctx context.Context) (conflow.Result, error) {
	switch val := r.value.(type) {
	case []interface{}:
		for k, v := range val {
			_, err := r.blockPublisher.PublishBlock(&RangeEntry{
				id:    r.entry.id,
				key:   k,
				value: v,
			}, nil)
			if err != nil {
				return nil, err
			}
		}
	case map[string]interface{}:
		for k, v := range val {
			_, err := r.blockPublisher.PublishBlock(&RangeEntry{
				id:    r.entry.id,
				key:   k,
				value: v,
			}, nil)
			if err != nil {
				return nil, err
			}
		}
	default:
		return nil, fmt.Errorf("invalid value for range: %T", r.value)
	}

	return nil, nil
}

// @block
type RangeEntry struct {
	// @id
	id conflow.ID
	// @read_only
	key interface{}
	// @read_only
	value interface{}
}

func (r *RangeEntry) ID() conflow.ID {
	return r.id
}
