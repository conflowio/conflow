// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

import (
	"time"

	"github.com/opsidian/parsley/parsley"
)

// Directive provides a way to describe alternate runtime execution
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Directive
type Directive interface {
	Block
	EvalStageAware
	RuntimeConfig() RuntimeConfig
}

type RuntimeConfig struct {
	Skip     bool
	Timeout  time.Duration
	Retry    Retryable
	Triggers []ID
}

func (r RuntimeConfig) Merge(r2 RuntimeConfig) RuntimeConfig {
	if r2.Skip {
		r.Skip = r2.Skip
	}
	if r2.Timeout > 0 {
		r.Timeout = r2.Timeout
	}
	if r2.Retry != nil {
		r.Retry = r2.Retry
	}
	if len(r2.Triggers) > 0 {
		r.Triggers = r2.Triggers
	}
	return r
}

func (r RuntimeConfig) IsTrigger(trigger ID) bool {
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

// DirectiveTransformerRegistryAware is an interface to get a block node transformer registry
type DirectiveTransformerRegistryAware interface {
	DirectiveTransformerRegistry() parsley.NodeTransformerRegistry
}
