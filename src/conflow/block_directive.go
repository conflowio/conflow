// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conflow

import (
	"time"
)

// BlockDirective provides a way to add metadata or define alternate runtime execution for blocks
//
//counterfeiter:generate . BlockDirective
type BlockDirective interface {
	Block
	RuntimeConfigOption
}

type RuntimeConfigOption interface {
	ApplyToRuntimeConfig(*RuntimeConfig)
}

var _ RuntimeConfigOption = &RuntimeConfig{}

type RuntimeConfig struct {
	Skip        *bool
	Timeout     *time.Duration
	Triggers    []ID
	RetryConfig *RetryConfig
}

type RetryConfig struct {
	Limit int
}

func (r *RuntimeConfig) ApplyToRuntimeConfig(r2 *RuntimeConfig) {
	if r.Skip != nil {
		r2.Skip = r.Skip
	}
	if r.Timeout != nil {
		r2.Timeout = r.Timeout
	}
	if r.Triggers != nil {
		r2.Triggers = r.Triggers
	}
	if r.RetryConfig != nil {
		r2.RetryConfig = r.RetryConfig
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
