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
type Required struct {
	// @id
	id conflow.ID
}

func (r *Required) ID() conflow.ID {
	return r.id
}

func (r *Required) ApplyToSchema(s schema.Schema) error {
	return nil
}
