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
type Value struct {
	// @id
	id conflow.ID
}

func (v *Value) ID() conflow.ID {
	return v.id
}

func (v *Value) ApplyToSchema(s schema.Schema) error {
	return schema.UpdateMetadata(s, func(meta schema.MetadataAccessor) error {
		meta.SetAnnotation(conflow.AnnotationValue, "true")
		return nil
	})
}
