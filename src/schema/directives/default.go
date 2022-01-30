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
type Default struct {
	// @id
	id conflow.ID
	// @value
	// @required
	value interface{}
}

func (d *Default) ID() conflow.ID {
	return d.id
}

func (d *Default) ApplyToSchema(s schema.Schema) error {
	value, err := s.ValidateValue(d.value)
	if err != nil {
		return fmt.Errorf("default value is invalid: %w", err)
	}

	switch st := s.(type) {
	case *schema.Array:
		st.Default = value.([]interface{})
	case *schema.Object:
		st.Default = schema.ObjectPtr(value.(map[string]interface{}))
	case *schema.Boolean:
		st.Default = schema.BooleanPtr(value.(bool))
	case *schema.Integer:
		st.Default = schema.IntegerPtr(value.(int64))
	case *schema.Number:
		st.Default = schema.NumberPtr(value.(float64))
	case *schema.String:
		st.Default = schema.StringPtr(value.(string))
	default:
		return fmt.Errorf("default directive can not be applied to %T", s)
	}

	return nil
}