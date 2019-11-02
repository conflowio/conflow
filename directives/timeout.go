// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"time"

	"github.com/opsidian/basil/basil"
)

//go:generate basil generate
type Timeout struct {
	id       basil.ID      `basil:"id"`
	duration time.Duration `basil:"value,required"`
}

func (t *Timeout) ID() basil.ID {
	return t.id
}

func (t *Timeout) RuntimeConfig() basil.RuntimeConfig {
	return basil.RuntimeConfig{
		Timeout: t.duration,
	}
}

func (t *Timeout) EvalStage() basil.EvalStage {
	return basil.EvalStageInit
}
