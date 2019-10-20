// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directive

import "github.com/opsidian/basil/basil"

// TODO: introduce the []basil.ID type
//go:generate basil generate
type Triggers struct {
	id       basil.ID      `basil:"id"`
	blockIDs []interface{} `basil:"value,name=block_ids"`
}

func (t Triggers) ID() basil.ID {
	return t.id
}

func (t Triggers) ApplyDirective(blockCtx basil.BlockContext, container basil.BlockContainer) error {
	tid := container.Trigger()
	if tid == "" {
		return nil
	}

	for _, id := range t.blockIDs {
		if tid == id {
			return nil
		}
	}

	container.Skip()
	return nil
}
