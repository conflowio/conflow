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

// Function is the directive for marking functions as basil functions
//
// @block
type Function struct {
	// @id
	id   basil.ID
	Path string
}

func (f *Function) ID() basil.ID {
	return f.id
}

func (f *Function) ApplyToSchema(s schema.Schema) error {
	if _, ok := s.(*schema.Function); !ok {
		return fmt.Errorf("@function can only be used on a function")
	}

	return nil
}
