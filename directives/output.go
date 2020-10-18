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
type Output struct {
	id         basil.ID `basil:"id"`
	outputType string   `basil:"name=type,required"`
}

func (o *Output) ID() basil.ID {
	return o.id
}

func (o *Output) ApplyToParameterConfig(config *basil.ParameterConfig) {
	config.Output = util.BoolPtr(true)
	config.Type = util.StringPtr(o.outputType)
}

func (o *Output) EvalStage() basil.EvalStage {
	return basil.EvalStageParse
}
