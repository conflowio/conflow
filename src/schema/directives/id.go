// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"fmt"

	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/conflow/annotations"
	"github.com/conflowio/conflow/src/schema"
)

// @block "directive"
type ID struct {
	// @id
	id conflow.ID
}

func (i *ID) ID() conflow.ID {
	return i.id
}

func (i *ID) ApplyToSchema(s schema.Schema) error {
	ss, ok := s.(*schema.String)
	if !ok {
		return fmt.Errorf("id annotation can only be set on a conflow.ID field")
	}

	if ss.Format != schema.FormatConflowID {
		return fmt.Errorf("id annotation can only be set on a conflow.ID field")
	}

	return schema.UpdateMetadata(s, func(meta schema.MetadataAccessor) error {
		meta.SetAnnotation(annotations.ID, "true")
		return nil
	})
}
