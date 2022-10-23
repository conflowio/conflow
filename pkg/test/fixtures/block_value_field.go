// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package fixtures

import "github.com/conflowio/conflow/pkg/conflow"

// @block "configuration"
type BlockValueField struct {
	// @id
	IDField conflow.ID
	// @value
	value interface{}
}

func (b *BlockValueField) ID() conflow.ID {
	return b.IDField
}
