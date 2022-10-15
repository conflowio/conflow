// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"errors"

	"github.com/conflowio/conflow/src/conflow/annotations"

	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/schema"
)

// @block "directive"
type Key struct {
	// @id
	id conflow.ID
}

func (k *Key) ID() conflow.ID {
	return k.id
}

func (k *Key) ApplyToSchema(s schema.Schema) error {
	switch s.(type) {
	case *schema.String:
		return schema.UpdateMetadata(s, func(meta schema.MetadataAccessor) error {
			meta.SetAnnotation(annotations.Key, "true")
			return nil
		})
	default:
		return errors.New("key directive can only be applied to string fields")
	}

}
