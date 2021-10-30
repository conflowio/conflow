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
//   eval_stage = "parse"
// }
type Output struct {
	// @id
	id conflow.ID
}

func (o *Output) ID() conflow.ID {
	return o.id
}

func (o *Output) ApplyToParameterConfig(config *conflow.ParameterConfig) {
	config.Output = util.BoolPtr(true)
}
