// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/util"
)

// @block {
//   type = "directive"
//   eval_stage = "init"
// }
type Skip struct {
	// @id
	id conflow.ID
	// @value
	// @default true
	when bool
}

func (s *Skip) ID() conflow.ID {
	return s.id
}

func (s *Skip) ApplyToRuntimeConfig(config *conflow.RuntimeConfig) {
	config.Skip = util.BoolPtr(s.when)
}
