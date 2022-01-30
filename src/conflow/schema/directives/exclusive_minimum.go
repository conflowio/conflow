// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"fmt"

	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/conflow/schema"
)

// @block "directive"
type ExclusiveMinimum struct {
	// @id
	id conflow.ID
	// @value
	// @required
	value interface{}
}

func (e *ExclusiveMinimum) ID() conflow.ID {
	return e.id
}

func (e *ExclusiveMinimum) ApplyToSchema(s schema.Schema) error {
	value, err := s.ValidateValue(e.value)
	if err != nil {
		return fmt.Errorf("exclusive_minimum value is invalid: %w", err)
	}

	switch st := s.(type) {
	case *schema.Integer:
		st.ExclusiveMinimum = schema.IntegerPtr(value.(int64))
	case *schema.Number:
		st.ExclusiveMinimum = schema.NumberPtr(value.(float64))
	default:
		return fmt.Errorf("exclusive_minimum directive can not be applied to %T", s)
	}

	return nil
}
