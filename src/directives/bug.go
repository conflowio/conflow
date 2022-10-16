// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import "github.com/conflowio/conflow/src/conflow"

//	@block {
//	  type = "directive"
//	  eval_stage = "ignore"
//	}
type Bug struct {
	// @id
	id conflow.ID
	// @value
	// @required
	description string
}

func (b *Bug) ID() conflow.ID {
	return b.id
}

func (b *Bug) ApplyToRuntimeConfig(*conflow.RuntimeConfig) {
}

func (b *Bug) ApplyToParameterConfig(*conflow.ParameterConfig) {
}
