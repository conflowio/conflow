// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package directives

import (
	"github.com/opsidian/basil/basil/block"
)

func Registry() block.InterpreterRegistry {
	return block.InterpreterRegistry{
		"const":       ConstInterpreter{},
		"default":     DefaultInterpreter{},
		"deprecated":  DeprecatedInterpreter{},
		"description": DescriptionInterpreter{},
		"enum":        EnumInterpreter{},
		"examples":    ExamplesInterpreter{},
		"generated":   GeneratedInterpreter{},
		"id":          IDInterpreter{},
		"ignore":      IgnoreInterpreter{},
		"name":        NameInterpreter{},
		"read_only":   ReadOnlyInterpreter{},
		"required":    RequiredInterpreter{},
		"result_type": ResultTypeInterpreter{},
		"eval_stage":  EvalStageInterpreter{},
		"title":       TitleInterpreter{},
		"types":       TypesInterpreter{},
		"value":       ValueInterpreter{},
		"write_only":  WriteOnlyInterpreter{},
	}
}
