// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"time"

	"github.com/opsidian/basil/examples/common"

	"github.com/opsidian/basil/basil/block"

	"github.com/opsidian/basil/basil"
)

//go:generate basil generate
type Timeout struct {
	id      basil.ID      `basil:"id"`
	timeout time.Duration `basil:"required,stage=init"`
}

func (t *Timeout) ID() basil.ID {
	return t.id
}

func (t *Timeout) Context(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, t.timeout)
}

func (t *Timeout) ParseContextOverride() basil.ParseContextOverride {
	return basil.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"sleep": common.SleepInterpreter{},
		},
	}
}
