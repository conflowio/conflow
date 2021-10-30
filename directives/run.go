// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"github.com/opsidian/conflow/basil"
	"github.com/opsidian/conflow/util"
)

// @block {
//   eval_stage = "init"
// }
type Run struct {
	// @id
	id basil.ID
	// @value
	// @default true
	when bool
}

func (r *Run) ID() basil.ID {
	return r.id
}

func (r *Run) ApplyToRuntimeConfig(config *basil.RuntimeConfig) {
	config.Skip = util.BoolPtr(!r.when)
}
