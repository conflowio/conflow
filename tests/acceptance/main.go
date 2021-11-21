// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package acceptance

import (
	"github.com/conflowio/conflow/blocks"
	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/block"
)

// @block "main"
type Main struct {
	// @id
	id conflow.ID
}

func (m *Main) ID() conflow.ID {
	return m.id
}

func (m *Main) ParseContextOverride() conflow.ParseContextOverride {
	return conflow.ParseContextOverride{
		BlockTransformerRegistry: block.InterpreterRegistry{
			"print":   blocks.PrintInterpreter{},
			"println": blocks.PrintlnInterpreter{},
		},
	}
}
