// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package blocks

import (
	"fmt"

	"github.com/opsidian/basil/basil"
)

// Print will write a string to the standard output
//go:generate basil generate
type Print struct {
	id    basil.ID    `basil:"id"`
	value interface{} `basil:"value,required"`
}

func (p *Print) ID() basil.ID {
	return p.id
}

func (p *Print) Main(ctx basil.BlockContext) error {
	fmt.Print(p.value)
	return nil
}
