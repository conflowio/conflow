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
type Types struct {
	// @id
	id basil.ID
	// @value
	// @required
	value []string
}

func (t *Types) ID() basil.ID {
	return t.id
}

func (t *Types) ApplyToSchema(s schema.Schema) error {
	u, ok := s.(*schema.Untyped)
	if !ok {
		return fmt.Errorf("@types can only be used on an interface{} type")
	}

	u.Types = t.value

	return nil
}
