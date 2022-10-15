// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"fmt"

	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/schema"
)

// @block "directive"
type MaxProperties struct {
	// @id
	id conflow.ID
	// @value
	// @required
	value int64
}

func (m *MaxProperties) ID() conflow.ID {
	return m.id
}

func (m *MaxProperties) ApplyToSchema(s schema.Schema) error {
	switch st := s.(type) {
	case *schema.Map:
		st.MaxProperties = schema.Pointer(m.value)
		return nil
	case *schema.Object:
		st.MaxProperties = schema.Pointer(m.value)
		return nil
	default:
		return fmt.Errorf("max_properties directive can not be applied to %T", s)
	}
}
