// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/schema"
)

// @block
type Dependency struct {
	// @id
	id conflow.ID
	// @value
	Name string
}

func (d *Dependency) ID() conflow.ID {
	return d.id
}

func (d *Dependency) ApplyToSchema(s schema.Schema) error {
	return schema.UpdateMetadata(s, func(meta schema.MetadataAccessor) error {
		meta.SetAnnotation("dependency", d.Name)
		return nil
	})
}
