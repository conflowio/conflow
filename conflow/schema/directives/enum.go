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

// @block
type Enum struct {
	// @id
	id conflow.ID
	// @value
	// @required
	values []interface{}
}

func (e *Enum) ID() conflow.ID {
	return e.id
}

func (e *Enum) ApplyToSchema(s schema.Schema) error {
	for _, e := range e.values {
		if err := s.ValidateValue(e); err != nil {
			return fmt.Errorf("enum value %q is invalid: %w", schema.UntypedValue().StringValue(e), err)
		}
	}

	switch st := s.(type) {
	case *schema.Array:
		st.Enum = make([][]interface{}, len(e.values))
		for i, v := range e.values {
			st.Enum[i] = v.([]interface{})
		}
	case *schema.Object:
		st.Enum = make([]map[string]interface{}, len(e.values))
		for i, v := range e.values {
			st.Enum[i] = v.(map[string]interface{})
		}
	case *schema.Boolean:
		st.Enum = make([]bool, len(e.values))
		for i, v := range e.values {
			st.Enum[i] = v.(bool)
		}
	case *schema.Integer:
		st.Enum = make([]int64, len(e.values))
		for i, v := range e.values {
			st.Enum[i] = v.(int64)
		}
	case *schema.Number:
		st.Enum = make([]float64, len(e.values))
		for i, v := range e.values {
			st.Enum[i] = v.(float64)
		}
	case *schema.String:
		st.Enum = make([]string, len(e.values))
		for i, v := range e.values {
			st.Enum[i] = v.(string)
		}
	default:
		return fmt.Errorf("const directive can not be applied to %T", s)
	}

	return nil
}
