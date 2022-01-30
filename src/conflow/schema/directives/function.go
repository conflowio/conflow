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

// Function is the directive for marking functions as conflow functions
//
// @block "directive"
type Function struct {
	// @id
	id   conflow.ID
	Path string
}

func (f *Function) ID() conflow.ID {
	return f.id
}

func (f *Function) ApplyToSchema(s schema.Schema) error {
	if _, ok := s.(*schema.Function); !ok {
		return fmt.Errorf("@function can only be used on a function")
	}

	return nil
}
