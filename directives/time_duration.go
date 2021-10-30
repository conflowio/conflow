// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"github.com/opsidian/conflow/basil"
	"github.com/opsidian/conflow/basil/schema"
)

// @block {
//   eval_stage = "parse"
// }
type TimeDuration struct {
	schema.TimeDuration
}

func (t *TimeDuration) ApplyToParameterConfig(config *basil.ParameterConfig) {
	config.Schema = &t.TimeDuration
}
