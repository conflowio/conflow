// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"fmt"

	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/schema"
)

// @block "directive"
type ExclusiveMaximum struct {
	// @id
	id conflow.ID
	// @value
	// @required
	value interface{}
}

func (e *ExclusiveMaximum) ID() conflow.ID {
	return e.id
}

func (e *ExclusiveMaximum) ApplyToSchema(s schema.Schema) error {
	value, err := s.ValidateValue(e.value)
	if err != nil {
		return fmt.Errorf("exclusive_maximum value is invalid: %w", err)
	}

	switch st := s.(type) {
	case *schema.Integer:
		st.ExclusiveMaximum = schema.IntegerPtr(value.(int64))
	case *schema.Number:
		st.ExclusiveMaximum = schema.NumberPtr(value.(float64))
	default:
		return fmt.Errorf("exclusive_maximum directive can not be applied to %T", s)
	}

	return nil
}
