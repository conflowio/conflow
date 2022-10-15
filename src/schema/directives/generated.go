// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/conflow/annotations"
	"github.com/conflowio/conflow/src/schema"
)

// @block "directive"
type Generated struct {
	// @id
	id conflow.ID
}

func (g *Generated) ID() conflow.ID {
	return g.id
}

func (g *Generated) ApplyToSchema(s schema.Schema) error {
	return schema.UpdateMetadata(s, func(meta schema.MetadataAccessor) error {
		meta.SetAnnotation(annotations.Generated, "true")
		return nil
	})
}
