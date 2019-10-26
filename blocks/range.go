// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blocks

import (
	"fmt"

	"github.com/opsidian/basil/basil/block"

	"github.com/opsidian/basil/basil"
)

//go:generate basil generate
type Range struct {
	id    basil.ID `basil:"id"`
	value interface{}
	entry *RangeEntry `basil:"generated"`
}

func (r *Range) ParseContextOverride() basil.ParseContextOverride {
	return basil.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"entry": RangeEntryInterpreter{},
		},
	}
}

func (r *Range) ID() basil.ID {
	return r.id
}

func (r *Range) Main(ctx basil.BlockContext) error {
	switch val := r.value.(type) {
	case []interface{}:
		for k, v := range val {
			err := ctx.PublishBlock(&RangeEntry{
				id:    r.entry.id,
				key:   k,
				value: v,
			})
			if err != nil {
				return err
			}
		}
	case map[string]interface{}:
		for k, v := range val {
			err := ctx.PublishBlock(&RangeEntry{
				id:    r.entry.id,
				key:   k,
				value: v,
			})
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("invalid value for range: %T", r.value)
	}

	return nil
}

//go:generate basil generate
type RangeEntry struct {
	id    basil.ID    `basil:"id"`
	key   interface{} `basil:"output"`
	value interface{} `basil:"output"`
}

func (r *RangeEntry) ID() basil.ID {
	return r.id
}
