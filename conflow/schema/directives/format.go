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
type Format struct {
	// @id
	id conflow.ID
	// @value
	// @required
	value string
}

func (f *Format) ID() conflow.ID {
	return f.id
}

func (f *Format) ApplyToSchema(s schema.Schema) error {
	switch st := s.(type) {
	case *schema.String:
		st.Format = f.value
		return nil
	default:
		return fmt.Errorf("format directive can not be applied to %T", s)
	}
}
