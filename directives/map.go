// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"github.com/opsidian/basil/basil"
	"github.com/opsidian/basil/basil/schema"
)

// @block {
//   eval_stage = "parse"
// }
type Map struct {
	schema.Map
}

func (m *Map) ApplyToParameterConfig(config *basil.ParameterConfig) {
	config.Schema = &m.Map
}

func (m *Map) ParseContextOverride() basil.ParseContextOverride {
	return basil.ParseContextOverride{
		BlockTransformerRegistry: schemaRegistry,
	}
}
