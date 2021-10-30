// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import "github.com/opsidian/conflow/conflow"

// @block {
//   eval_stage = "init"
// }
type Retry struct {
	// @id
	id conflow.ID
	// @value
	// @required
	// @default -1
	// @minimum -1
	// @maximum 2147483647
	limit int64
}

func (r *Retry) ID() conflow.ID {
	return r.id
}

func (r *Retry) ApplyToRuntimeConfig(config *conflow.RuntimeConfig) {
	config.RetryConfig = &conflow.RetryConfig{Limit: int(r.limit)}
}
