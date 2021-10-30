// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import "github.com/opsidian/conflow/basil"

// @block {
//   eval_stage = "ignore"
// }
type Bug struct {
	// @id
	id basil.ID
	// @value
	// @required
	description string
}

func (b *Bug) ID() basil.ID {
	return b.id
}

func (b *Bug) ApplyToRuntimeConfig(*basil.RuntimeConfig) {
}

func (b *Bug) ApplyToParameterConfig(*basil.ParameterConfig) {
}
