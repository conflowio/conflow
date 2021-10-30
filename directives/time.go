// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"github.com/conflowio/conflow/conflow"
	"github.com/conflowio/conflow/conflow/schema"
)

// @block {
//   eval_stage = "parse"
// }
type Time struct {
	schema.Time
}

func (t *Time) ApplyToParameterConfig(config *conflow.ParameterConfig) {
	config.Schema = &t.Time
}
