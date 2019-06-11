// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package fixtures

import (
	"context"
	"time"

	"github.com/opsidian/basil/basil"
)

//go:generate basil generate
type BlockWithContext struct {
	IDField basil.ID `basil:"id"`
}

func (b *BlockWithContext) ID() basil.ID {
	return b.IDField
}

func (b *BlockWithContext) Context(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, 1*time.Second)
}
