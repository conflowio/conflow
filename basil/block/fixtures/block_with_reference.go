// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package fixtures

import "github.com/opsidian/basil/basil"

//go:generate basil generate
type BlockWithReference struct {
	IDField basil.ID `basil:"id,reference"`
}

func (b *BlockWithReference) ID() basil.ID {
	return b.IDField
}
