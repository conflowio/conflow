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
type DependentRequired struct {
	// @id
	id conflow.ID
	// @value
	// @required
	Value map[string]schema.UniqueNonEmptyStringList
}

func (d *DependentRequired) ID() conflow.ID {
	return d.id
}

func (d *DependentRequired) ApplyToSchema(s schema.Schema) error {
	switch st := s.(type) {
	case *schema.Object:
		st.DependentRequired = d.Value
		return nil
	default:
		return fmt.Errorf("dependent_required directive can not be applied to %T", s)
	}
}
