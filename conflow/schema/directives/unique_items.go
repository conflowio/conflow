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
type UniqueItems struct {
	// @id
	id conflow.ID
}

func (u *UniqueItems) ID() conflow.ID {
	return u.id
}

func (u *UniqueItems) ApplyToSchema(s schema.Schema) error {
	switch st := s.(type) {
	case *schema.Array:
		st.UniqueItems = true
	default:
		return fmt.Errorf("unique_items directive can not be applied to %T", s)
	}

	return nil
}
