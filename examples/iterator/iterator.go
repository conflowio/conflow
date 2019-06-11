// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"errors"

	"github.com/opsidian/basil/basil/block"

	"github.com/opsidian/basil/basil"
)

//go:generate basil generate
type Iterator struct {
	id    basil.ID                `basil:"id"`
	count int64                   `basil:"required"`
	it    chan basil.BlockMessage `basil:"output,block"`
}

func (i *Iterator) ID() basil.ID {
	return i.id
}

func (it *Iterator) Main(ctx basil.BlockContext) error {
	for i := int64(0); i < it.count; i++ {
		message := basil.NewBlockMessage(&It{value: i})
		it.it <- message
		ctx.Logger().Debug().ID("id", it.id).Int64("i", i).Msg("ITERATOR WAITING")
		select {
		case <-message.WaitGroup().Wait():
		case <-ctx.Context().Done():
			return errors.New("aborted")
		}
		ctx.Logger().Debug().ID("id", it.id).Int64("i", i).Msg("ITERATOR FINISHED")
	}

	return nil
}

func (it *Iterator) ParseContextOverride() basil.ParseContextOverride {
	return basil.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"it": ItInterpreter{},
		},
	}
}

//go:generate basil generate
type It struct {
	id    basil.ID `basil:"id"`
	value int64    `basil:"output"`
}

func (i *It) ID() basil.ID {
	return i.id
}
