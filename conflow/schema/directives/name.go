// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"github.com/opsidian/conflow/conflow"
	"github.com/opsidian/conflow/conflow/schema"
)

// @block
type Name struct {
	// @id
	id conflow.ID
	// @value
	// @required
	Value string
}

func (n *Name) ID() conflow.ID {
	return n.id
}

func (n *Name) ApplyToSchema(s schema.Schema) error {
	return nil
}
