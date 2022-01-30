// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"

	"github.com/conflowio/conflow/src/conflow"
)

// Fail will error for the given tries
// @block "task"
type Fail struct {
	// @id
	id conflow.ID
	// @required
	triesRequired int64
	// @read_only
	tries int64
}

func (f *Fail) ID() conflow.ID {
	return f.id
}

func (f *Fail) Run(ctx context.Context) (conflow.Result, error) {
	f.tries++
	if f.tries < f.triesRequired {
		return conflow.Retry("I am destined to fail, but eventually I will succeed"), nil
	}

	return nil, nil
}
