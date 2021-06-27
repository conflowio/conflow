// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import "github.com/opsidian/basil/basil"

// @block {
//   eval_stage = "resolve"
// }
type Triggers struct {
	// @id
	id basil.ID
	// @value
	// @required
	// @name "block_ids"
	blockIDs []interface{}
}

func (t *Triggers) ID() basil.ID {
	return t.id
}

func (t *Triggers) ApplyToRuntimeConfig(config *basil.RuntimeConfig) {
	// TODO: introduce the []basil.ID type for blockIDs
	triggers := make([]basil.ID, len(t.blockIDs))
	for i, id := range t.blockIDs {
		triggers[i] = basil.ID(id.(string))
	}
	config.Triggers = triggers
}
