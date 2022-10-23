// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import "github.com/conflowio/conflow/pkg/conflow"

//	@block {
//	  type = "directive"
//	  eval_stage = "ignore"
//	}
type Deprecated struct {
	// @id
	id conflow.ID
	// @value
	// @required
	description string
}

func (d *Deprecated) ID() conflow.ID {
	return d.id
}

func (d *Deprecated) ApplyToRuntimeConfig(*conflow.RuntimeConfig) {
}

func (d *Deprecated) ApplyToParameterConfig(*conflow.ParameterConfig) {
}
