// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/schema"
	"github.com/conflowio/conflow/src/util/ptr"
)

//	@block {
//	  type = "directive"
//	  eval_stage = "parse"
//	}
type Input struct {
	// @id
	id       conflow.ID
	required bool
	// @name "type"
	// @required
	schema schema.Schema
}

func (i *Input) ID() conflow.ID {
	return i.id
}

func (i *Input) ApplyToParameterConfig(config *conflow.ParameterConfig) {
	config.Input = ptr.To(true)
	config.Required = ptr.To(i.required)
	config.Schema = i.schema
}

func (i *Input) ParseContextOverride() conflow.ParseContextOverride {
	return conflow.ParseContextOverride{
		BlockTransformerRegistry: schemaRegistry,
	}
}
