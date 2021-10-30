// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"github.com/opsidian/conflow/basil"
	"github.com/opsidian/conflow/basil/schema"
)

// @block
type Ignore struct {
	// @id
	id basil.ID
}

func (i *Ignore) ID() basil.ID {
	return i.id
}

func (i *Ignore) ApplyToSchema(schema.Schema) error {
	return nil
}
