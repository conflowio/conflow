// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"fmt"

	"github.com/opsidian/conflow/basil"
	"github.com/opsidian/conflow/basil/schema"
)

// @block
type MultipleOf struct {
	// @id
	id basil.ID
	// @value
	// @required
	value interface{}
}

func (m *MultipleOf) ID() basil.ID {
	return m.id
}

func (m *MultipleOf) ApplyToSchema(s schema.Schema) error {
	if err := s.ValidateValue(m.value); err != nil {
		return fmt.Errorf("multiple_of value is invalid: %w", err)
	}

	switch st := s.(type) {
	case *schema.Integer:
		st.MultipleOf = schema.IntegerPtr(m.value.(int64))
	case *schema.Number:
		st.MultipleOf = schema.NumberPtr(m.value.(float64))
	default:
		return fmt.Errorf("multiple_of directive can not be applied to %T", s)
	}

	return nil
}
