// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package fixtures

import "github.com/opsidian/basil/basil"

//go:generate basil generate
type BlockWithDefault struct {
	IDField basil.ID `basil:"id"`
	Value   string   `basil:"default=\"foo\""`
}

func (b *BlockWithDefault) ID() basil.ID {
	return b.IDField
}