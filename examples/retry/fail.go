// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"

	"github.com/opsidian/conflow/basil"
)

// Fail will error for the given tries
// @block
type Fail struct {
	// @id
	id basil.ID
	// @required
	triesRequired int64
	// @read_only
	tries int64
}

func (f *Fail) ID() basil.ID {
	return f.id
}

func (f *Fail) Run(ctx context.Context) (basil.Result, error) {
	f.tries++
	if f.tries < f.triesRequired {
		return basil.Retry("I am destined to fail, but eventually I will succeed"), nil
	}

	return nil, nil
}
