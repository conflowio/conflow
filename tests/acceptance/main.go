// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package acceptance

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/block"
	"github.com/opsidian/basil/blocks"
)

//go:generate basil generate
type Main struct {
	id basil.ID `basil:"id"`
}

func (m *Main) ID() basil.ID {
	return m.id
}

func (m *Main) ParseContextOverride() basil.ParseContextOverride {
	return basil.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"print":   blocks.PrintInterpreter{},
			"println": blocks.PrintlnInterpreter{},
		},
	}
}
