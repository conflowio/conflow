// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blocks

import (
	"errors"

	"github.com/opsidian/basil/basil"
)

//go:generate basil generate
type Fail struct {
	id  basil.ID `basil:"id"`
	msg string   `basil:"value,required"`
}

func (f *Fail) ID() basil.ID {
	return f.id
}

func (f *Fail) Main(blockCtx basil.BlockContext) error {
	return errors.New(f.msg)
}
