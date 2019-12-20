// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

import (
	"time"
)

// BlockDirective provides a way to add metadata or define alternate runtime execution for blocks
//go:generate counterfeiter . BlockDirective
type BlockDirective interface {
	Block
	EvalStageAware
	RuntimeConfigOption
}

type RuntimeConfigOption interface {
	ApplyToRuntimeConfig(*RuntimeConfig)
}

var _ RuntimeConfigOption = &RuntimeConfig{}

type RuntimeConfig struct {
	Skip     *bool
	Timeout  *time.Duration
	Retry    Retryable
	Triggers []ID
}

func (r *RuntimeConfig) ApplyToRuntimeConfig(r2 *RuntimeConfig) {
	if r.Skip != nil {
		r2.Skip = r.Skip
	}
	if r.Timeout != nil {
		r2.Timeout = r.Timeout
	}
	if r.Retry != nil {
		r2.Retry = r.Retry
	}
	if r.Triggers != nil {
		r2.Triggers = r.Triggers
	}
}

func (r *RuntimeConfig) IsTrigger(trigger ID) bool {
	if len(r.Triggers) == 0 {
		return true
	}
	for _, t := range r.Triggers {
		if t == trigger {
			return true
		}
	}
	return false
}
