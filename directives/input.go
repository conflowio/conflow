// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/util"
)

// @block {
//   type = "directive"
//   eval_stage = "parse"
// }
type Input struct {
	// @id
	id       conflow.ID
	required bool
}

func (i *Input) ID() conflow.ID {
	return i.id
}

func (i *Input) ApplyToParameterConfig(config *conflow.ParameterConfig) {
	config.Input = util.BoolPtr(true)
	config.Required = util.BoolPtr(i.required)
}
