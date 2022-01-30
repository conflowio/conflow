// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package fixtures

import (
	"context"
	"time"

	"github.com/conflowio/conflow/src/conflow"
)

// @block "configuration"
type BlockWithContext struct {
	// @id
	IDField conflow.ID
	// @required
	// @eval_stage "init"
	timeout time.Duration
}

func (b *BlockWithContext) Context(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, b.timeout)
}

func (b *BlockWithContext) ID() conflow.ID {
	return b.IDField
}
