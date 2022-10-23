// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"time"

	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/conflow/types"
	"github.com/conflowio/conflow/pkg/util/ptr"
)

//	@block {
//	  type = "directive"
//	  eval_stage = "init"
//	}
type Timeout struct {
	// @id
	id conflow.ID
	// @value
	// @required
	duration types.Duration
}

func (t *Timeout) ID() conflow.ID {
	return t.id
}

func (t *Timeout) ApplyToRuntimeConfig(config *conflow.RuntimeConfig) {
	config.Timeout = ptr.To(time.Duration(t.duration))
}
