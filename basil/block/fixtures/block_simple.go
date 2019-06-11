// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package fixtures

import "github.com/opsidian/basil/basil"

//go:generate basil generate
type BlockSimple struct {
	IDField basil.ID    `basil:"id"`
	Value   interface{} `basil:"value"`
}

func (b *BlockSimple) ID() basil.ID {
	return b.IDField
}

func (b *BlockSimple) Foo() string {
	return "bar"
}
