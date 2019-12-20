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
type Input struct {
	id        basil.ID `basil:"id"`
	inputType string   `basil:"name=type,required"`
	required  bool
}

func (i *Input) ID() basil.ID {
	return i.id
}

func (i *Input) ApplyToParameterConfig(config *basil.ParameterConfig) {
	config.Input = util.BoolPtr(true)
	config.Type = util.StringPtr(i.inputType)
	config.Required = util.BoolPtr(i.required)
}

func (i *Input) EvalStage() basil.EvalStage {
	return basil.EvalStageParse
}
