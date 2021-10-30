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
type Skip struct {
	// @id
	id basil.ID
	// @value
	// @default true
	when bool
}

func (s *Skip) ID() basil.ID {
	return s.id
}

func (s *Skip) ApplyToRuntimeConfig(config *basil.RuntimeConfig) {
	config.Skip = util.BoolPtr(s.when)
}
