// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blocks

import (
	"context"
	"errors"
	"time"

	"github.com/conflowio/conflow/conflow"
)

// @block
type Sleep struct {
	// @id
	id conflow.ID
	// @value
	// @required
	duration time.Duration
}

func (s *Sleep) ID() conflow.ID {
	return s.id
}

func (s *Sleep) Run(ctx context.Context) (conflow.Result, error) {
	select {
	case <-time.After(s.duration):
		return nil, nil
	case <-ctx.Done():
		return nil, errors.New("aborted")
	}
}
