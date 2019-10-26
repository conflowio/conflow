// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"time"

	"github.com/opsidian/basil/basil"
)

//go:generate basil generate
type Timeout struct {
	id       basil.ID      `basil:"id"`
	duration time.Duration `basil:"value"`
}

func (t Timeout) ID() basil.ID {
	return t.id
}

func (t Timeout) ApplyDirective(blockCtx basil.BlockContext, container basil.BlockContainer) error {
	if t.duration > 0 {
		container.SetTimeout(t.duration)
	}
	return nil
}
