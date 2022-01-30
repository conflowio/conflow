// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package fixtures

import "github.com/conflowio/conflow/src/conflow"

// @block "configuration"
type BlockWithDefault struct {
	// @id
	IDField conflow.ID
	// @value
	// @default "foo"
	Value string
}

func (b *BlockWithDefault) ID() conflow.ID {
	return b.IDField
}
