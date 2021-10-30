// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blocks

import (
	"context"
	"errors"

	"github.com/conflowio/conflow/conflow"
)

// @block
type Fail struct {
	// @id
	id conflow.ID
	// @value
	// @required
	msg string
}

func (f *Fail) ID() conflow.ID {
	return f.id
}

func (f *Fail) Run(ctx context.Context) (conflow.Result, error) {
	return nil, errors.New(f.msg)
}
