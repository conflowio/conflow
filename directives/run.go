// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/util"
)

//go:generate basil generate
type Run struct {
	id   basil.ID `basil:"id"`
	when bool     `basil:"value,default=true"`
}

func (r *Run) ID() basil.ID {
	return r.id
}

func (r *Run) ApplyToRuntimeConfig(config *basil.RuntimeConfig) {
	config.Skip = util.BoolPtr(!r.when)
}

func (r *Run) EvalStage() basil.EvalStage {
	return basil.EvalStageInit
}
