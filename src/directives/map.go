// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"github.com/conflowio/conflow/src/conflow"
	"github.com/conflowio/conflow/src/conflow/schema"
)

// @block {
//   type = "directive"
//   eval_stage = "parse"
// }
type Map struct {
	schema.Map
}

func (m *Map) ApplyToParameterConfig(config *conflow.ParameterConfig) {
	config.Schema = &m.Map
}

func (m *Map) ParseContextOverride() conflow.ParseContextOverride {
	return conflow.ParseContextOverride{
		BlockTransformerRegistry: schemaRegistry,
	}
}
