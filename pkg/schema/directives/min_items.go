// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"fmt"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/schema"
)

// @block "directive"
type MinItems struct {
	// @id
	id conflow.ID
	// @value
	// @required
	value int64
}

func (m *MinItems) ID() conflow.ID {
	return m.id
}

func (m *MinItems) ApplyToSchema(s schema.Schema) error {
	switch st := s.(type) {
	case *schema.Array:
		st.MinItems = m.value
		return nil
	default:
		return fmt.Errorf("min_items directive can not be applied to %T", s)
	}
}
