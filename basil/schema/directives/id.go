// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"fmt"

	"github.com/opsidian/conflow/basil"
	"github.com/opsidian/conflow/basil/schema"
)

// @block
type ID struct {
	// @id
	id basil.ID
}

func (i *ID) ID() basil.ID {
	return i.id
}

func (i *ID) ApplyToSchema(s schema.Schema) error {
	ss, ok := s.(*schema.String)
	if !ok {
		return fmt.Errorf("id annotation can only be set on a basil.ID field")
	}

	if ss.Format != schema.FormatBasilID {
		return fmt.Errorf("id annotation can only be set on a basil.ID field")
	}

	return schema.UpdateMetadata(s, func(meta schema.MetadataAccessor) error {
		meta.SetAnnotation("id", "true")
		return nil
	})
}
