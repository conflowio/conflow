// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package fixtures

import "github.com/opsidian/conflow/basil"

// @block
type BlockRequiredField struct {
	// @id
	IDField basil.ID
	// @required
	required interface{}
}

func (b *BlockRequiredField) ID() basil.ID {
	return b.IDField
}
