// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import "github.com/conflowio/conflow/conflow"

// @block {
//   eval_stage = "ignore"
// }
type Todo struct {
	// @id
	id conflow.ID
	// @value
	// @required
	description string
}

func (t *Todo) ID() conflow.ID {
	return t.id
}

func (t *Todo) ApplyToRuntimeConfig(*conflow.RuntimeConfig) {
}

func (t *Todo) ApplyToParameterConfig(*conflow.ParameterConfig) {
}
