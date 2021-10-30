// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"fmt"

	"github.com/opsidian/conflow/basil/schema"
	"github.com/opsidian/conflow/conflow"
)

// @block
type Types struct {
	// @id
	id conflow.ID
	// @value
	// @required
	value []string
}

func (t *Types) ID() conflow.ID {
	return t.id
}

func (t *Types) ApplyToSchema(s schema.Schema) error {
	u, ok := s.(*schema.Untyped)
	if !ok {
		return fmt.Errorf("@types can only be used on an interface{} type")
	}

	u.Types = t.value

	return nil
}
