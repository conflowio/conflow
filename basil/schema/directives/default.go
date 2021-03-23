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

//go:generate basil generate
type Default struct {
	// @id
	id basil.ID
	// @value
	// @required
	value interface{}
}

func (d *Default) ID() basil.ID {
	return d.id
}

func (d *Default) ApplyToSchema(s schema.Schema) error {
	if err := s.ValidateValue(d.value); err != nil {
		return fmt.Errorf("default value is invalid: %w", err)
	}

	switch st := s.(type) {
	case *schema.Array:
		st.Default = d.value.([]interface{})
	case *schema.Object:
		st.Default = schema.ObjectPtr(d.value.(map[string]interface{}))
	case *schema.Boolean:
		st.Default = schema.BooleanPtr(d.value.(bool))
	case *schema.Integer:
		st.Default = schema.IntegerPtr(d.value.(int64))
	case *schema.Number:
		st.Default = schema.NumberPtr(d.value.(float64))
	case *schema.String:
		st.Default = schema.StringPtr(d.value.(string))
	default:
		return fmt.Errorf("default directive can not be applied to %T", s)
	}

	return nil
}
