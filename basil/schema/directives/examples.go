// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/schema"
)

// @block
type Examples struct {
	// @id
	id basil.ID
	// @value
	// @required
	value []interface{}
}

func (e *Examples) ID() basil.ID {
	return e.id
}

func (e *Examples) ApplyToSchema(s schema.Schema) error {
	return schema.UpdateMetadata(s, func(meta schema.MetadataAccessor) error {
		meta.SetExamples(e.value)
		return nil
	})
}
