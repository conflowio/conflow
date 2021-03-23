// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/schema"
)

//go:generate basil generate
type Name struct {
	// @id
	id basil.ID
	// @value
	// @required
	Value string
}

func (n *Name) ID() basil.ID {
	return n.id
}

func (n *Name) ApplyToSchema(s schema.Schema) error {
	return nil
}
