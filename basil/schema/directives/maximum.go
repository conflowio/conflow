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
type Maximum struct {
	// @id
	id basil.ID
	// @value
	// @required
	value interface{}
}

func (m *Maximum) ID() basil.ID {
	return m.id
}

func (m *Maximum) ApplyToSchema(s schema.Schema) error {
	if err := s.ValidateValue(m.value); err != nil {
		return fmt.Errorf("maximum value is invalid: %w", err)
	}

	switch st := s.(type) {
	case *schema.Integer:
		st.Maximum = schema.IntegerPtr(m.value.(int64))
	case *schema.Number:
		st.Maximum = schema.NumberPtr(m.value.(float64))
	default:
		return fmt.Errorf("maximum directive can not be applied to %T", s)
	}

	return nil
}
