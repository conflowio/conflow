// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"fmt"
	"regexp"

	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/schema"
)

// @block "directive"
type Pattern struct {
	// @id
	id conflow.ID
	// @value
	// @required
	// @format "regex"
	value regexp.Regexp
}

func (p *Pattern) ID() conflow.ID {
	return p.id
}

func (p *Pattern) ApplyToSchema(s schema.Schema) error {
	switch st := s.(type) {
	case *schema.String:
		st.Pattern = &p.value
		return nil
	default:
		return fmt.Errorf("format directive can not be applied to %T", s)
	}
}
