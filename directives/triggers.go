// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import "github.com/opsidian/basil/basil"

//go:generate basil generate
type Triggers struct {
	id       basil.ID      `basil:"id"`
	blockIDs []interface{} `basil:"value,name=block_ids,required"`
}

func (t *Triggers) ID() basil.ID {
	return t.id
}

func (t *Triggers) RuntimeConfig() basil.RuntimeConfig {
	// TODO: introduce the []basil.ID type for blockIDs
	triggers := make([]basil.ID, len(t.blockIDs))
	for i, id := range t.blockIDs {
		triggers[i] = basil.ID(id.(string))
	}
	return basil.RuntimeConfig{
		Triggers: triggers,
	}
}

func (t *Triggers) EvalStage() basil.EvalStage {
	return basil.EvalStageResolve
}
