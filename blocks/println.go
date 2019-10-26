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

// Println will write a string followed by a new line to the standard output
//go:generate basil generate
type Println struct {
	id    basil.ID    `basil:"id"`
	value interface{} `basil:"value,required"`
}

func (p *Println) ID() basil.ID {
	return p.id
}

func (p *Println) Main(ctx basil.BlockContext) error {
	fmt.Println(p.value)
	return nil
}
