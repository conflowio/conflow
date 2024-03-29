// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"github.com/conflowio/conflow/pkg/conflow"
	"github.com/conflowio/conflow/pkg/schema"
	"github.com/conflowio/conflow/pkg/util/ptr"
)

//	@block {
//	  type = "directive"
//	  eval_stage = "parse"
//	}
type Output struct {
	// @id
	id conflow.ID
	// @name "type"
	// @required
	schema schema.Schema
}

func (o *Output) ID() conflow.ID {
	return o.id
}

func (o *Output) ApplyToParameterConfig(config *conflow.ParameterConfig) {
	config.Output = ptr.To(true)
	config.Schema = o.schema
}

func (o *Output) ParseContextOverride() conflow.ParseContextOverride {
	return conflow.ParseContextOverride{
		BlockTransformerRegistry: schemaRegistry,
	}
}
