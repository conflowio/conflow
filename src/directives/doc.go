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
type Doc struct {
	// @id
	id conflow.ID
	// @value
	// @required
	description string
}

func (d *Doc) ID() conflow.ID {
	return d.id
}

func (d *Doc) ApplyToRuntimeConfig(*conflow.RuntimeConfig) {
}

func (d *Doc) ApplyToParameterConfig(*conflow.ParameterConfig) {
}
