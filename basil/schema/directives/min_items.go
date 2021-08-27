// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"fmt"

	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/schema"
)

// @block
type MinItems struct {
	// @id
	id basil.ID
	// @value
	// @required
	value int64
}

func (m *MinItems) ID() basil.ID {
	return m.id
}

func (m *MinItems) ApplyToSchema(s schema.Schema) error {
	switch st := s.(type) {
	case *schema.Array:
		st.MinItems = schema.IntegerPtr(m.value)
		return nil
	default:
		return fmt.Errorf("min_items directive can not be applied to %T", s)
	}
}
