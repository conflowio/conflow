// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package basil

// ParameterDirective provides a way to add metadata for parameters
//go:generate counterfeiter . ParameterDirective
type ParameterDirective interface {
	Block
	ApplyToParameterConfig(*ParameterConfig)
}

type ParameterConfigOption interface {
	ApplyToParameterConfig(*ParameterConfig)
}

var _ ParameterConfigOption = ParameterConfig{}

type ParameterConfig struct {
	Input    *bool
	Required *bool
	Output   *bool
	Type     *string
}

func (p ParameterConfig) ApplyToParameterConfig(p2 *ParameterConfig) {
	if p.Input != nil {
		p2.Input = p.Input
	}

	if p.Required != nil {
		p2.Required = p.Required
	}

	if p.Output != nil {
		p2.Output = p.Output
	}

	if p.Type != nil {
		p2.Type = p.Type
	}
}
